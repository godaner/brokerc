package log

type Logger interface {
	SetDebug(debug bool)
	Debugf(fms string, arg ...interface{})
	Debug(arg ...interface{})
	Infof(fms string, arg ...interface{})
	Info(arg ...interface{})
	Warningf(fms string, arg ...interface{})
	Warning(arg ...interface{})
	Errorf(fms string, arg ...interface{})
	Error(arg ...interface{})
}
