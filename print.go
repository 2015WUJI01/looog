package looog

func Sync() { _ = logger.l.Sync() }

func Debug(args ...interface{}) { logger.l.Sugar().Debug(args) }
func Debugw(msg string, keysAndValues ...interface{}) {
	logger.l.Sugar().Debugw(msg, keysAndValues...)
}
func Debugf(template string, args ...interface{}) {
	logger.l.Sugar().Debugf(template, args...)
}

func Info(args ...interface{}) { logger.l.Sugar().Info(args) }
func Infow(msg string, keysAndValues ...interface{}) {
	logger.l.Sugar().Infow(msg, keysAndValues...)
}
func Infof(template string, args ...interface{}) {
	logger.l.Sugar().Infof(template, args...)
}

func Warn(args ...interface{}) { logger.l.Sugar().Warn(args) }
func Warnw(msg string, keysAndValues ...interface{}) {
	logger.l.Sugar().Warnw(msg, keysAndValues...)
}
func Warnf(template string, args ...interface{}) {
	logger.l.Sugar().Warnf(template, args...)
}

func Error(args ...interface{}) { logger.l.Sugar().Error(args) }
func Errorw(msg string, keysAndValues ...interface{}) {
	logger.l.Sugar().Errorw(msg, keysAndValues...)
}
func Errorf(template string, args ...interface{}) {
	logger.l.Sugar().Errorf(template, args...)
}

func Panic(args ...interface{}) { logger.l.Sugar().Panic(args) }
func Panicw(msg string, keysAndValues ...interface{}) {
	logger.l.Sugar().Panicw(msg, keysAndValues...)
}
func Panicf(template string, args ...interface{}) {
	logger.l.Sugar().Panicf(template, args...)
}

func Fatal(args ...interface{}) { logger.l.Sugar().Fatal(args) }
func Fatalw(msg string, keysAndValues ...interface{}) {
	logger.l.Sugar().Fatalw(msg, keysAndValues...)
}
func Fatalf(template string, args ...interface{}) {
	logger.l.Sugar().Fatalf(template, args...)
}

func Print(v ...interface{}) {
	Info(v...)
}

func Println(v ...interface{}) {
	Info(v...)
}

func Printf(format string, v ...interface{}) {
	Infof(format, v...)
}
