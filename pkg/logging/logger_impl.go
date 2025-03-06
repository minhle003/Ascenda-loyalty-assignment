package logging

import (
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.StandardLogger()
)

type logrusLogger struct {
	logger *logrus.Logger
}

func LogrusLogger() Logger {
	return &logrusLogger{logger: log}
}

func (l *logrusLogger) Trace(message string, data ...interface{}) {
	l.logger.Traceln(append(data, message))
}

func (l *logrusLogger) Debug(message string, data ...interface{}) {
	l.logger.Debugln(append(data, message))
}

func (l *logrusLogger) Info(message string, data ...interface{}) {
	l.logger.Infoln(append(data, message))
}

func (l *logrusLogger) Warn(message string, data ...interface{}) {
	l.logger.Warnln(append(data, message))
}

func (l *logrusLogger) Error(message string, data ...interface{}) {
	l.logger.Errorln(append(data, message))
}

func (l *logrusLogger) Critical(message string, data ...interface{}) {
	l.logger.Fatalln(append(data, message))
}
