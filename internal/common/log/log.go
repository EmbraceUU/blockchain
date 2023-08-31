package log

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

const (
	defaultPath         = "logs"
	defaultLevel        = "info"
	defaultMaxAge       = 604800
	defaultRotationTime = 86400
)

func WithLoggerName(name string) *logrus.Logger {
	if name == "" {
		return defaultLogger
	}

	logger, err := CreateNewLogger(defaultPath, name, defaultLevel, defaultMaxAge, defaultRotationTime)
	if err != nil {
		return defaultLogger
	}

	return logger
}

func WithLoggerNameAndLevel(name string, level string) *logrus.Logger {
	if name == "" || level == "" {
		return defaultLogger
	}

	logger, err := CreateNewLogger(defaultPath, name, level, defaultMaxAge, defaultRotationTime)
	if err != nil {
		return defaultLogger
	}

	return logger
}

func JsonPrint(obj interface{}) string {
	brs, _ := json.Marshal(obj)
	return string(brs)
}
