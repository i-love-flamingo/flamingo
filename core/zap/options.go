package zap

func WithLogSession(value bool) Option {
	return func(logger *Logger) {
		logger.logSession = value
	}
}

func WithFieldMap(fieldMap map[string]string) Option {
	return func(logger *Logger) {
		logger.fieldMap = fieldMap
	}
}

func WithArea(area string) Option {
	return func(logger *Logger) {
		logger.configArea = area
	}
}
