package log

import (
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

func init() {
	setDefaultOutputs()
}

// SetupDefaultLogger configure logger instance with user provided settings
func SetupDefaultLogger(logPath string, fileName string, level string, maxAge int, rotationTime int) (err error) {
	//Logger.LogPath = logPath
	//Logger.File = fileName
	//AddLogRotateHook(defaultLogger)

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	defaultLogger.SetLevel(lvl)

	ConfigLocalFilesystemLogger(defaultLogger, logPath, fileName, maxAge, rotationTime)

	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		err = os.MkdirAll(logPath, 0755)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

// setDefaultOutputs() this setups defaults used by the logger
// This allows it to be used without any user configuration
func setDefaultOutputs() {
	defaultLogger = logrus.New()
	defaultLogger.SetLevel(logrus.DebugLevel)
}

func ConfigLocalFilesystemLogger(l *logrus.Logger, logPath string, logFileName string, maxAge int, rotationTime int) {
	maxAgeDuration := time.Second * time.Duration(int64(maxAge))
	rotationTimeDuration := time.Second * time.Duration(int64(rotationTime))

	baseLogPaht := path.Join(logPath, logFileName)
	writer, err := rotatelogs.New(
		baseLogPaht+".%Y%m%d",
		rotatelogs.WithMaxAge(maxAgeDuration),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTimeDuration), // 日志切割时间间隔
	)
	if err != nil {
		l.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.TextFormatter{DisableColors: true, TimestampFormat: "2006-01-02 15:04:05.000"})

	//注意上面这个 and &amp;符号被转义了
	l.AddHook(lfHook)
}

func CreateNewLogger(logPath string, fileName string, level string, maxAge int, rotationTime int) (*logrus.Logger, error) {
	newLogger := logrus.New()

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	newLogger.SetLevel(lvl)

	if _, err = os.Stat(logPath); os.IsNotExist(err) {
		err = os.MkdirAll(logPath, 0755)
		if err != nil {
			return nil, err
		}
	}

	ConfigLocalFilesystemLogger(newLogger, logPath, fileName, maxAge, rotationTime)
	return newLogger, nil
}
