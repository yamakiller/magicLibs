package connection

import (
	"crypto/tls"
	"io"
	"time"
)

type CLIENT_STATUS uint8

const (
	CS_UNCONNECT = CLIENT_STATUS(0)
	CS_CONNECTING = CLIENT_STATUS(1)
	CS_CONNECTED = CLIENT_STATUS(2)
	CS_CLOSING   = CLIENT_STATUS(3)
)

//Client 客户端接口
type Client interface {
	Connect(string, time.Duration) error
	ConnectTls(string, time.Duration, *tls.Config) error
	IsConnected() bool
	Parse() (interface{}, error)
	SendTo(interface{}) error
	Close() error
}

//Exception 异常处理接口
type Exception interface {
	Error(error)
}

//Serialization 序列化反序列化接口
type Serialization interface {
	UnSeria(io.Reader) (interface{}, int, error)
	Seria(interface{}, io.Writer) (int, error)
}
