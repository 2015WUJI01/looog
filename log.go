package looog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

var logger = DefaultLogger()

type Logger struct {
	l     *zap.Logger
	mx    *sync.Mutex
	style string // "json" "console"
	level zapcore.Level
	ws    zapcore.WriteSyncer

	timeenc  zapcore.TimeEncoder
	levelenc zapcore.LevelEncoder

	caller    bool
	callerenc zapcore.CallerEncoder

	cores []zapcore.Core
}

type Option func(*Logger)

func DefaultLogger() *Logger {
	return New(
		OptionUseLevelFormat(LF_CAPITAL|LF_COLOR),
		OptionSetLogLevel(DebugLevel),
		OptionEnableCaller(true),
	)
}

func New(opts ...Option) *Logger {
	l := &Logger{
		mx: &sync.Mutex{},
	}

	for _, f := range opts {
		f(l)
	}

	var enccfg = zap.NewProductionEncoderConfig()
	enccfg.CallerKey = "file"
	enccfg.EncodeDuration = zapcore.MillisDurationEncoder // 执行时间，以秒为单位

	if l.timeenc != nil {
		enccfg.EncodeTime = l.timeenc
	} else {
		enccfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	}

	if l.levelenc != nil {
		enccfg.EncodeLevel = l.levelenc
	} else {
		enccfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	if l.ws == nil {
		l.ws = os.Stdout
	}

	if l.level == 0 {
		l.level = DebugLevel
	}

	var zapopts []zap.Option

	if l.caller {
		// 调用文件和行号，内部使用 runtime.Caller
		zapopts = append(zapopts, zap.AddCaller(), zap.AddCallerSkip(0))
		enccfg.EncodeCaller = l.callerenc
	}

	var enc zapcore.Encoder
	if l.style == "json" {
		enc = zapcore.NewJSONEncoder(enccfg)
	} else {
		enc = zapcore.NewConsoleEncoder(enccfg)
	}

	l.cores = append(l.cores, zapcore.NewCore(enc, l.ws, l.level))
	l.l = zap.New(zapcore.NewTee(l.cores...), zapopts...)
	return l
}

func OptionSetStyle(s string) Option {
	return func(l *Logger) { l.style = s }
}

func OptionSetLogLevel(lv Level) Option {
	return func(l *Logger) { l.level = lv }
}

// OptionSetOutput "stdout" "file"
func OptionSetOutput(ws zapcore.WriteSyncer) Option {
	return func(l *Logger) { l.ws = ws }
}

func OptionSetOutputFile(path string) Option {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		fmt.Printf("set log output file error: %s\n", err.Error())
	}
	return OptionSetOutput(f)
}

func OptionSetTimeFormat(layout string) Option {
	return func(l *Logger) {
		l.timeenc = zapcore.TimeEncoderOfLayout(layout)
	}
}

// level format
const (
	LF_NONE    = 0
	LF_CAPITAL = 1 << iota
	LF_COLOR
)

func OptionUseLevelFormat(flag int) Option {
	var enc zapcore.LevelEncoder
	switch {
	case flag == LF_NONE:
		enc = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {}
	case flag&LF_CAPITAL == LF_CAPITAL && flag&LF_COLOR == LF_COLOR:
		enc = zapcore.CapitalColorLevelEncoder
	case flag&LF_CAPITAL == LF_CAPITAL:
		enc = zapcore.CapitalLevelEncoder
	case flag&LF_COLOR == LF_COLOR:
		enc = zapcore.LowercaseColorLevelEncoder
	default:
		enc = zapcore.LowercaseLevelEncoder
	}
	return func(l *Logger) {
		l.levelenc = enc
	}
}

// caller style
const (
	CS_NONE = iota
	CS_SHORT
	CS_FULL
)

// OptionEnableCaller 开启 caller，可以指定 style 为 CS_FULL 为全路径，默认不指定则为为简短路径
func OptionEnableCaller(caller bool, style ...int) Option {
	s := CS_SHORT
	if len(style) > 0 {
		s = style[0]
	}
	return func(l *Logger) {
		l.caller = caller
		switch s {
		case CS_NONE:
			l.callerenc = func(c zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {}
		case CS_SHORT:
			l.callerenc = zapcore.ShortCallerEncoder
		case CS_FULL:
			l.callerenc = zapcore.FullCallerEncoder
		}
	}
}
