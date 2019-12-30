package listen

type ClientSocket interface {
	GetID() int32
	GetAddr() string
	Write([]byte) (int, error)
	Read([]byte) (int, error)
	Close()
}

type ListenSocket interface {
	Listen(addr string, cap int) error
	Accept() (ClientSocket, error)
	Close(int32)
	Shutdown()
}
