package router

//IRouteCtrl route control interface
type IRouteCtrl interface {
	GetName() string
	IncRef()
	DecRef() int
	Call(string, interface{}, interface{}) error
	Shutdown()
}
