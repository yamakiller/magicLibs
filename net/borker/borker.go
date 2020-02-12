package borker

//Borker 网络代理服务
type Borker interface {
	ListenAndServe(string) error
	Serve() error
	Shutdown()
}
