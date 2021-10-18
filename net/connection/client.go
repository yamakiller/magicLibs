package connection

import (
	"crypto/tls"
	"time"
)

//Client 客户端接口
type Client interface {
	Connect(string, time.Duration) error
	ConnectTls(string, time.Duration, *tls.Config) error
	Parse() (interface{}, error)
	SendTo(interface{}) error
	GetWriteBytes() uint64
	GetReadBytes() uint64
	Close() error
	Shutdown()
}
