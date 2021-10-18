package interfaces

import (
	"io"
)

//Exception 异常处理接口
type Exception interface {
	Error(error)
}

//Serialization 序列化反序列化接口
type Serialization interface {
	UnSeria(io.Reader) (interface{}, int, error)
	Seria(interface{}, io.Writer) (int, error)
}
