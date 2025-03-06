package logging

type Logger interface {
	Trace(message string, data ...interface{})
	Debug(message string, data ...interface{})
	Info(message string, data ...interface{})
	Warn(message string, data ...interface{})
	Error(message string, data ...interface{})
	Critical(message string, data ...interface{})
}

type LogEntity struct {
	Name  string
	Value interface{}
}
