package looog

import (
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

// LogConfig 用于定制一个 zap.core
type LogConfig struct {
	style string // "json" "console"
	level zapcore.Level
	out   zapcore.WriteSyncer

	// 高级设置
	timekey     string
	timeEncoder zapcore.TimeEncoder

	levelkey     string
	levelEncoder zapcore.LevelEncoder
}

func (cfg *LogConfig) zapcore() zapcore.Core {
	return zapcore.NewCore(cfg.encoder(), cfg.out, cfg.level)
}

func (cfg *LogConfig) encoderConfig() zapcore.EncoderConfig {
	// 用自带的生产环境配置来改
	c := zap.NewProductionEncoderConfig()
	// 改掉不太喜欢的配置项
	c.CallerKey = "file"
	c.EncodeDuration = zapcore.MillisDurationEncoder // 执行时间，以秒为单位
	// 用用户自定义的覆盖掉默认项
	c.TimeKey = cfg.timekey
	c.EncodeTime = cfg.timeEncoder
	c.LevelKey = cfg.levelkey
	c.EncodeLevel = cfg.levelEncoder
	return c
}

func (cfg *LogConfig) encoder() zapcore.Encoder {
	c := cfg.encoderConfig()
	if cfg.style == string(LS_JSON) {
		return zapcore.NewJSONEncoder(c)
	}
	return zapcore.NewConsoleEncoder(c)
}

type AdvanceLogConfig func(*LogConfig)

type LogStyle string

const (
	LS_JSON    LogStyle = "json"
	LS_CONSOLE          = "console"
)

func DefaultLogConfig() LogConfig {
	return NewLogConfig(LS_CONSOLE, DebugLevel,
		EnableTime(true), EnableLevel(true),
		SetLevelFormat(LFMTcapital|LFMTcolor),
	)
}

func NewLogConfig(style LogStyle, level Level, advance ...AdvanceLogConfig) LogConfig {
	cfg := LogConfig{
		style:       string(style),
		level:       level,
		out:         os.Stdout,
		timeEncoder: zapcore.TimeEncoderOfLayout("2006-01-02 15:04:15"),
	}
	for _, f := range advance {
		f(&cfg)
	}
	return cfg
}

func EnableTime(ok bool) AdvanceLogConfig {
	return func(cfg *LogConfig) {
		if !ok {
			cfg.timekey = ""
		} else if cfg.timekey == "" {
			cfg.timekey = "time"
		}
	}
}

func SetTimeFormat(layout string) AdvanceLogConfig {
	return func(cfg *LogConfig) {
		cfg.timeEncoder = zapcore.TimeEncoderOfLayout(layout)
	}
}

type LevelFormatFlag int

const (
	LFMTcapital LevelFormatFlag = 1 << iota
	LFMTcolor
)

func EnableLevel(ok bool) AdvanceLogConfig {
	return func(cfg *LogConfig) {
		if !ok {
			cfg.levelkey = ""
		} else if cfg.levelkey == "" {
			cfg.levelkey = "time"
		}
	}
}

func SetLevelFormat(flag LevelFormatFlag) AdvanceLogConfig {
	var enc zapcore.LevelEncoder
	switch {
	case flag == LFMTcapital^LFMTcolor:
		enc = zapcore.CapitalColorLevelEncoder
	case flag == LFMTcapital:
		enc = zapcore.CapitalLevelEncoder
	case flag == LFMTcolor:
		enc = zapcore.LowercaseColorLevelEncoder
	default:
		enc = zapcore.LowercaseLevelEncoder
	}
	return func(cfg *LogConfig) {
		cfg.levelEncoder = enc
	}
}

func SetOutputFile(path string) AdvanceLogConfig {
	return func(cfg *LogConfig) {
		f, _ := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
		cfg.out = f
	}
}

type Option func(*Logger)

func DefaultLogger() *Logger {
	return New(DefaultLogConfig(), OptionEnableCaller(false))
}

func New(config LogConfig, opts ...Option) *Logger {
	l := &Logger{mx: &sync.Mutex{}}
	for _, f := range opts {
		f(l)
	}
	l.l = zap.New(zapcore.NewTee(append(l.cores, config.zapcore())...), zap.WithCaller(l.caller), zap.AddCallerSkip(1))
	return l
}

func (l *Logger) Add(config LogConfig) {
	logger.mx.Lock()
	l.cores = append(l.cores, config.zapcore())
	logger.mx.Unlock()
}

func (l *Logger) Rebuild() {
	_ = logger.l.Sync()
	l.mx.Lock()
	l.l = zap.New(zapcore.NewTee(l.cores...), zap.WithCaller(l.caller), zap.AddCallerSkip(1))
	l.mx.Unlock()
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
