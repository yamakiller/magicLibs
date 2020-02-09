package listener

import "net"

//Listener 监听接口
type Listener interface {
	Addr() net.Addr
	Wait()
	Close() error
}
