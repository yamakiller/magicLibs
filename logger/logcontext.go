package logger

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

//LogContext doc
//@Summary Log context
//@Struct LogContext
//@Member (string) log file path
//@Member (*os.File) log file handle
//@Member (logrus.Levle) log level limit
//@Member (*logrus.Logger) log object
//@Member (int32) log mail queue number
//@Member (int32) log mail queue max
//@Member (chan LogMessage) log mail queue chan
//@Member (chan struct{}) log system shutdown chan
//@Member (sync.WaitGroup) log system is shutdown success
type LogContext struct {
	_filPath    string
	_filHandle  *os.File
	_logLevel   logrus.Level
	_logHandle  *logrus.Logger
	_logMailNum int32
	_logMailMax int32
	_logMailbox chan LogMessage
	_logStop    chan struct{}
	_logWait    sync.WaitGroup
}

//SetFilPath doc
//@Summary Setting log file name
//@Method SetFilPath
//@Param (string) file name
func (slf *LogContext) SetFilPath(v string) {
	slf._filPath = v
}

//SetLevel doc
//@Summary Setting log level limit
//@Method SetLevel
//@Param (logrus.Level) log level
func (slf *LogContext) SetLevel(v logrus.Level) {
	slf._logLevel = v
}

//SetFilHandle doc
//@Summary Setting log file handle
//@Method SetFilHandle
//@Param (*os.File) log file
func (slf *LogContext) SetFilHandle(v *os.File) {
	slf._filHandle = v
}

//SetHandle doc
//@SummarySetting log object
//@Method SetHandle
//@Param (*logrus.Logger)
func (slf *LogContext) SetHandle(v *logrus.Logger) {
	slf._logHandle = v
}

//SetMailMax doc
//@Summary Setting log mail max
//@Method SetMailMax
//@Param (int)
func (slf *LogContext) SetMailMax(v int) {
	slf._logMailMax = int32(v)
}

//SetFormatter doc
//@Summary Setting log format
//@Method SetFormatter
//@Param (logrus.Formatter)
func (slf *LogContext) SetFormatter(f logrus.Formatter) {
	slf._logHandle.SetFormatter(f)
}

//Initial doc
//@Summary initail logger
//@Method Initial
func (slf *LogContext) Initial() {
	slf._logMailbox = make(chan LogMessage, slf._logMailMax)
	slf._logStop = make(chan struct{})
	if slf._logHandle != nil {
		slf._logHandle.SetLevel(slf._logLevel)
	}
}

func (slf *LogContext) run() int {
	select {
	case <-slf._logStop:
		return -1
	case msg := <-slf._logMailbox:
		slf.write(&msg)
		atomic.AddInt32(&slf._logMailNum, -1)
		return 0
	}
}

func (slf *LogContext) exit() {
	slf._logWait.Done()
}

func (slf *LogContext) write(msg *LogMessage) {
	switch msg._level {
	case uint32(logrus.ErrorLevel):
		slf._logHandle.WithFields(logrus.Fields{"prefix": msg._prefix}).Errorln(msg._message)
	case uint32(logrus.InfoLevel):
		slf._logHandle.WithFields(logrus.Fields{"prefix": msg._prefix}).Infoln(msg._message)
	case uint32(logrus.TraceLevel):
		slf._logHandle.WithFields(logrus.Fields{"prefix": msg._prefix}).Traceln(msg._message)
	case uint32(logrus.DebugLevel):
		slf._logHandle.WithFields(logrus.Fields{"prefix": msg._prefix}).Debugln(msg._message)
	case uint32(logrus.WarnLevel):
		slf._logHandle.WithFields(logrus.Fields{"prefix": msg._prefix}).Warningln(msg._message)
	case uint32(logrus.FatalLevel):
		slf._logHandle.WithFields(logrus.Fields{"prefix": msg._prefix}).Fatalln(msg._message)
	case uint32(logrus.PanicLevel):
		slf._logHandle.WithFields(logrus.Fields{"prefix": msg._prefix}).Panicln(msg._message)
	}
}

func (slf *LogContext) getPrefix(owner uint32) string {
	return fmt.Sprintf("[.%08x]", owner)
}

func (slf *LogContext) push(data LogMessage) {
	select {
	case slf._logMailbox <- data:
	}

	atomic.AddInt32(&slf._logMailNum, 1)
}

//Redirect doc
//@Summary Redirect log file
//@Method Redirect
func (slf *LogContext) Redirect() {

	if slf._filPath == "" {
		slf._logHandle.SetOutput(os.Stdout)
		return
	}

	f, err := os.OpenFile(slf._filPath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return
	}

	slf._filHandle = f
	slf._logHandle.SetOutput(f)
}

//Mount doc
//@Summary Mount log module
//@Method Mount
func (slf *LogContext) Mount() {
	slf._logWait.Add(1)
	go func() {
		for {
			if slf.run() != 0 {
				break
			}
		}
		slf.exit()
	}()
}

//Close doc
//@Summary Turn off the logging system
//@Method Close
func (slf *LogContext) Close() {
	for {
		if atomic.LoadInt32(&slf._logMailNum) > 0 {
			time.Sleep(time.Millisecond * 10)
			continue
		}
		break
	}

	close(slf._logStop)
	slf._logWait.Wait()
	close(slf._logMailbox)
	if slf._filHandle != nil {
		slf._filHandle.Close()
	}
}

//Error doc
//@Summary Output error log
//@Method Error
//@Param (int32) owner
//@Param (string) format
//@Param (...interface{}) args
func (slf *LogContext) Error(owner uint32, fmrt string, args ...interface{}) {
	slf.push(spawnMessage(uint32(logrus.ErrorLevel), slf.getPrefix(owner), fmt.Sprintf(fmrt, args...)))
}

//Info doc
//@Summary Output information log
//@Method Info
//@Param (int32) owner
//@Param (string) format
//@Param (...interface{}) args
func (slf *LogContext) Info(owner uint32, fmrt string, args ...interface{}) {
	slf.push(spawnMessage(uint32(logrus.InfoLevel), slf.getPrefix(owner), fmt.Sprintf(fmrt, args...)))
}

//Warning doc
//@Summary Output warning log
//@Method Warning
//@Param (int32) owner
//@Param (string) format
//@Param (...interface{}) args
func (slf *LogContext) Warning(owner uint32, fmrt string, args ...interface{}) {
	slf.push(spawnMessage(uint32(logrus.WarnLevel), slf.getPrefix(owner), fmt.Sprintf(fmrt, args...)))
}

//Panic doc
//@Summary Output program crash log
//@Method Panic
//@Param (int32) owner
//@Param (string) format
//@Param (...interface{}) args
func (slf *LogContext) Panic(owner uint32, fmrt string, args ...interface{}) {
	slf.push(spawnMessage(uint32(logrus.PanicLevel), slf.getPrefix(owner), fmt.Sprintf(fmrt, args...)))
}

//Fatal doc
//@Summary Output critical error log
//@Method Fatal
//@Param (int32) owner
//@Param (string) format
//@Param (...interface{}) args
func (slf *LogContext) Fatal(owner uint32, fmrt string, args ...interface{}) {
	slf.push(spawnMessage(uint32(logrus.FatalLevel), slf.getPrefix(owner), fmt.Sprintf(fmrt, args...)))
}

//Debug doc
//@Summary Output Debug log
//@Method Debug
//@Param (int32) owner
//@Param (string) format
//@Param (...interface{}) args
func (slf *LogContext) Debug(owner uint32, fmrt string, args ...interface{}) {
	slf.push(spawnMessage(uint32(logrus.DebugLevel), slf.getPrefix(owner), fmt.Sprintf(fmrt, args...)))
}

//Trace doc
//@Summary Output trace log
//@Method Trace
//@Param (int32) owner
//@Param (string) format
//@Param (...interface{}) args
func (slf *LogContext) Trace(owner uint32, fmrt string, args ...interface{}) {
	slf.push(spawnMessage(uint32(logrus.TraceLevel), slf.getPrefix(owner), fmt.Sprintf(fmrt, args...)))
}

func spawnMessage(level uint32, prefix string, message string) LogMessage {
	return LogMessage{_level: level, _prefix: prefix, _message: message}
}
