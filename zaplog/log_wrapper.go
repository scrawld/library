package zaplog

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	Logger.Sugar().Debug(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	Logger.Sugar().Info(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	Logger.Sugar().Warn(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	Logger.Sugar().Error(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(template string, args ...interface{}) {
	Logger.Sugar().Debugf(template, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(template string, args ...interface{}) {
	Logger.Sugar().Infof(template, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(template string, args ...interface{}) {
	Logger.Sugar().Warnf(template, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(template string, args ...interface{}) {
	Logger.Sugar().Errorf(template, args...)
}
