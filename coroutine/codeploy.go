package coroutine

import "math"

const (
	//EnvKey desc: coroutine pool env name
	EnvKey = "coroutine pool"
)

//Deploy desc
//@Struct Deploy desc coroutine pool deploy informat (json)
type Deploy struct {
	Max  int `json:"max"`
	Min  int `json:"min"`
	Task int `json:"task-limit"`
}

//NewDefault desc
//@Method NewDefault desc: create default coroutine pool deploy informat
func NewDefault() *Deploy {
	return &Deploy{math.MaxInt16, 32, 64}
}
