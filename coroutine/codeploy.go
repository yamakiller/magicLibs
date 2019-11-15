package coroutine

import "math"

//Deploy desc
//@struct Deploy desc coroutine pool deploy informat (json)
type Deploy struct {
	Max  int `json:"max"`
	Min  int `json:"min"`
	Task int `json:"task-limit"`
}

//NewDefault desc
//@method NewDefault desc: create default coroutine pool deploy informat
func NewDefault() *Deploy {
	return &Deploy{32, math.MaxInt16, 64}
}
