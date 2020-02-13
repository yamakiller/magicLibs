package borker

import "github.com/yamakiller/magicLibs/net/listener"

//Borker 网络代理服务
type Borker interface {
	ListenAndServe(string) error
	Listener() listener.Listener
	Serve() error
	Shutdown()
}
