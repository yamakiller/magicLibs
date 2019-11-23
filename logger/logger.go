package logger

import (
	"runtime"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

const (
	EnvKey = "log"
)

//Logger desc
//@interface Logger desc: Log module interface
type Logger interface {
	run() int
	exit()

	write(msg *LogMessage)
	getPrefix(owner uint32) string

	Initial()
	Mount()
	Redirect()
	Close()
	Error(owner uint32, fmrt string, args ...interface{})
	Info(owner uint32, fmrt string, args ...interface{})
	Warning(owner uint32, fmrt string, args ...interface{})
	Panic(owner uint32, fmrt string, args ...interface{})
	Fatal(owner uint32, fmrt string, args ...interface{})
	Debug(owner uint32, fmrt string, args ...interface{})
	Trace(owner uint32, fmrt string, args ...interface{})
}

//Spawn desc
//@type Spawn desc: Log object maker
type Spawn func() Logger

var (
	defaultLevel       = logrus.PanicLevel
	defaultSize        = 512
	defaultFile        = ""
	defaultSpawnLogger = func() Logger {
		l := LogContext{}
		l.SetFilPath(defaultFile)
		l.SetHandle(logrus.New())
		l.SetMailMax(defaultSize)
		l.SetLevel(defaultLevel)

		formatter := new(prefixed.TextFormatter)
		formatter.FullTimestamp = true
		if runtime.GOOS == "windows" {
			formatter.DisableColors = true
		} else {
			formatter.SetColorScheme(&prefixed.ColorScheme{
				PrefixStyle: "blue+b"})
		}
		l.SetFormatter(formatter)
		l.Initial()
		l.Redirect()
		return &l
	}

	defaultHandle Logger
)

//New desc
//@method New desc: create an Logger object
//@param (Spawn) Logger make method
//@return (Logger) log object
func New(spawn Spawn) Logger {

	if spawn == nil {
		r := defaultSpawnLogger()
		return r
	}

	r := spawn()
	return r
}

//WithDefault desc
//@method WithDefault desc: Set the default log handle
//@param (Logger) logger object
func WithDefault(log Logger) {
	defaultHandle = log
}

//Error desc
//@method Error desc: Output error log
//@param (int32) owner
//@param (string) format
//@param (...interface{}) args
func Error(owner uint32, fmrt string, args ...interface{}) {
	if defaultHandle == nil {
		return
	}
	defaultHandle.Error(owner, fmrt, args...)
}

//Info desc
//@method Info desc: Output information log
//@param (int32) owner
//@param (string) format
//@param (...interface{}) args
func Info(owner uint32, fmrt string, args ...interface{}) {
	if defaultHandle == nil {
		return
	}
	defaultHandle.Info(owner, fmrt, args...)
}

//Warning desc
//@method Warning desc: Output warning log
//@param (int32) owner
//@param (string) format
//@param (...interface{}) args
func Warning(owner uint32, fmrt string, args ...interface{}) {
	if defaultHandle == nil {
		return
	}
	defaultHandle.Warning(owner, fmrt, args...)
}

//Panic desc
//@method Panic desc: Output program crash log
//@param (int32) owner
//@param (string) format
//@param (...interface{}) args
func Panic(owner uint32, fmrt string, args ...interface{}) {
	if defaultHandle == nil {
		return
	}
	defaultHandle.Panic(owner, fmrt, args...)
}

//Fatal desc
//@method Fatal desc: Output critical error log
//@param (int32) owner
//@param (string) format
//@param (...interface{}) args
func Fatal(owner uint32, fmrt string, args ...interface{}) {
	if defaultHandle == nil {
		return
	}
	defaultHandle.Fatal(owner, fmrt, args...)
}

//Debug desc
//@method Debug desc: Output Debug log
//@param (int32) owner
//@param (string) format
//@param (...interface{}) args
func Debug(owner uint32, fmrt string, args ...interface{}) {
	if defaultHandle == nil {
		return
	}
	defaultHandle.Debug(owner, fmrt, args...)
}

//Trace desc
//@method Trace desc: Output trace log
//@param (int32) owner
//@param (string) format
//@param (...interface{}) args
func Trace(owner uint32, fmrt string, args ...interface{}) {
	if defaultHandle == nil {
		return
	}
	defaultHandle.Trace(owner, fmrt, args...)
}
