package log

import (
	"os"
)

// Info handler takes any input returns unformatted output to infoLogger writer
func Info(v ...interface{}) {
	defaultLogger.Info(v...)
}

// Infof handler takes any input infoLogger returns formatted output to infoLogger writer
func Infof(data string, v ...interface{}) {
	defaultLogger.Infof(data, v...)
}

// Infoln handler takes any input infoLogger returns formatted output to infoLogger writer
func Infoln(v ...interface{}) {
	defaultLogger.Infoln(v...)
}

// Debug handler takes any input returns unformatted output to infoLogger writer
func Debug(v ...interface{}) {
	defaultLogger.Debug(v...)
}

// Debugf handler takes any input infoLogger returns formatted output to infoLogger writer
func Debugf(data string, v ...interface{}) {
	defaultLogger.Debugf(data, v...)
}

// Debugln handler takes any input infoLogger returns formatted output to infoLogger writer
func Debugln(v ...interface{}) {
	defaultLogger.Debugln(v...)
}

// Warn handler takes any input returns unformatted output to warnLogger writer
func Warn(v ...interface{}) {
	defaultLogger.Warn(v...)
}

// Warnf handler takes any input returns unformatted output to warnLogger writer
func Warnf(data string, v ...interface{}) {
	defaultLogger.Warnf(data, v...)
}

// Error handler takes any input returns unformatted output to errorLogger writer
func Error(v ...interface{}) {
	defaultLogger.Error(v...)
}

// Errorf handler takes any input returns unformatted output to errorLogger writer
func Errorf(data string, v ...interface{}) {
	defaultLogger.Errorf(data, v...)
}

// Fatal handler takes any input returns unformatted output to fatalLogger writer
func Fatal(v ...interface{}) {
	// Send to Output instead of Fatal to allow us to increase the output depth by 1 to make sure the correct file is displayed
	defaultLogger.Fatal(v...)
	os.Exit(1)
}

// Fatalf handler takes any input returns unformatted output to fatalLogger writer
func Fatalf(data string, v ...interface{}) {
	// Send to Output instead of Fatal to allow us to increase the output depth by 1 to make sure the correct file is displayed
	defaultLogger.Fatalf(data, v...)
	os.Exit(1)
}
