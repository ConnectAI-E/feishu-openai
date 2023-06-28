package logger

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strings"
)

var logger = logrus.New()

func init() {

	logger.SetFormatter(&formatter{})

	logger.SetReportCaller(true)

	gin.DefaultWriter = logger.Out

	// 设置日志级别 支持
	//PanicLevel
	//FatalLevel
	//ErrorLevel
	//WarnLevel
	//InfoLevel
	//DebugLevel
	logger.Level = logrus.InfoLevel

}

type Fields logrus.Fields

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Debugf(format, args...)
	}
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Warnf(format, args...)
	}
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Errorf(format, args...)
	}
}

// Fatalf logs a message at level Fatal on the standard logger.
func Fatalf(format string, args ...interface{}) {
	if logger.Level >= logrus.FatalLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Fatalf(format, args...)
	}
}

func Debug(format string, args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Debug(format, args)
	}
}

// Info logs a message at level Info on the standard logger.
func Info(format string, args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Info(format, args)
	}
}

// Warn logs a message at level Warn on the standard logger.
func Warn(format string, args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Warn(format, args)
	}
}

// Error logs a message at level Error on the standard logger.
func Error(format string, args ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Error(format, args)
	}
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(format string, args ...interface{}) {
	if logger.Level >= logrus.FatalLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Fatal(format, args)
	}
}

// Formatter implements logrus.Formatter interface.
type formatter struct {
	prefix string
}

// Format building log message.
func (f *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var sb bytes.Buffer

	sb.WriteString("[" + strings.ToUpper(entry.Level.String()) + "]")
	sb.WriteString(entry.Time.Format("2006-01-02 15:04:05"))
	sb.WriteString(" ")
	//sb.WriteString(" ")
	//sb.WriteString(f.prefix)
	sb.WriteString(entry.Message)
	sb.WriteString("\n")

	return sb.Bytes(), nil
}
