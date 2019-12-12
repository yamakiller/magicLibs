package logger

import (
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

const (
	//EnvKey Logger env key index
	EnvKey = "log"
)

//Logger doc
//@Summary Log module interface
//@Interface Logger
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

//Spawn doc
//@Summary Log object maker
//@type Spawn
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
		formatter.TimestampFormat = "2006-01-02 15:04:05"
		formatter.SetColorScheme(&prefixed.ColorScheme{
			PrefixStyle:    "white+h",
			TimestampStyle: "black+h"})
		l.SetFormatter(formatter)
		l.Initial()
		l.Redirect()
		return &l
	}

	defaultHandle Logger
)

//New doc
//@Summary create an Logger object
//@Method New
//@Param (Spawn) Logger make method
//@Return (Logger) log object
func New(spawn Spawn) Logger {

	if spawn == nil {
		r := defaultSpawnLogger()
		return r
	}

	r := spawn()
	return r
}

//WithDefault doc
//@Summary Set the default log handle
//@Method WithDefault
//@Param (Logger) logger object
func WithDefault(log Logger) {
	defaultHandle = log
}

//Error doc
//@Summary Output error log
//@Method Error
//@Param (int32) owner
//@Param (string) format
//@Param (...interface{}) args
func Error(owner uint32, fmrt string, args ...interface{}) {
	if defaultHandle == nil {
		return
	}
	defaultHandle.Error(owner, fmrt, args...)
}

//Info doc
//@Summary Output information log
//@Method Info
//@Param (int32) owner
//@Param (string) format
//@Param (...interface{}) args
func Info(owner uint32, fmrt string, args ...interface{}) {

	if defaultHandle == nil {
		return
	}

	defaultHandle.Info(owner, fmrt, args...)
}

//Warning doc
//@Summary Output warning log
//@Method Warning
//@Param (int32) owner
//@Param (string) format
//@Param (...interface{}) args
func Warning(owner uint32, fmrt string, args ...interface{}) {
	if defaultHandle == nil {
		return
	}
	defaultHandle.Warning(owner, fmrt, args...)
}

//Panic doc
//@Summary Output program crash log
//@Method Panic
//@Param (int32) owner
//@Param (string) format
//@Param (...interface{}) args
func Panic(owner uint32, fmrt string, args ...interface{}) {
	if defaultHandle == nil {
		return
	}
	defaultHandle.Panic(owner, fmrt, args...)
}

//Fatal doc
//@Summary Output critical error log
//@Method Fatal
//@Param (int32) owner
//@Param (string) format
//@Param (...interface{}) args
func Fatal(owner uint32, fmrt string, args ...interface{}) {
	if defaultHandle == nil {
		return
	}
	defaultHandle.Fatal(owner, fmrt, args...)
}

//Debug doc
//@Summary Output Debug log
//@Method Debug
//@Param (int32) owner
//@Param (string) format
//@Param (...interface{}) args
func Debug(owner uint32, fmrt string, args ...interface{}) {
	if defaultHandle == nil {
		return
	}
	defaultHandle.Debug(owner, fmrt, args...)
}

//Trace doc
//@Summary Output trace log
//@Method Trace
//@Param (int32) owner
//@Param (string) format
//@Param (...interface{}) args
func Trace(owner uint32, fmrt string, args ...interface{}) {
	if defaultHandle == nil {
		return
	}
	defaultHandle.Trace(owner, fmrt, args...)
}
