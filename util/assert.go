package util

import (
	"fmt"
	"runtime"
	"strings"
)

//Assert doc
//@Method Assert @Summary Assert boolean and output error message
//@Param (bool) false assert
//@Param (string) error message
func Assert(isAs bool, errMsg string) {
	if !isAs {
		_, file, inline, ok := runtime.Caller(2)
		panic(fmt.Sprintf("%s %d %v\n%s", file, inline, ok, errMsg))
	}
}

//AssertEmpty doc
//@Method AssertEmtpy @Summary Assert Nil and output an error message
//@Param (interface{}) is null assert
//@Param (string) error message
func AssertEmpty(isNull interface{}, errMsg string) {
	if isNull == nil {
		_, file, inline, ok := runtime.Caller(2)
		panic(fmt.Sprintf("%s %d %v\n%s", file, inline, ok, errMsg))
	}
}

//GetStack doc
//@Method GetStack @Summary Return current stack information
//@Return (string)
func GetStack() string {
	var name, file string
	var line int
	var pc [16]uintptr

	n := runtime.Callers(4, pc[:])
	callers := pc[:n]
	frames := runtime.CallersFrames(callers)
	for {
		frame, more := frames.Next()
		file = frame.File
		line = frame.Line
		name = frame.Function
		if !strings.HasPrefix(name, "runtime.") || !more {
			break
		}
	}

	var str string
	switch {
	case name != "":
		str = fmt.Sprintf("%v:%v", name, line)
	case file != "":
		str = fmt.Sprintf("%v:%v", file, line)
	default:
		str = fmt.Sprintf("pc:%x", pc)
	}
	return "stacktrace:\n" + str
}
