package router

type IRouteCtrl interface {
	GetKey() string
	IncRef()
	DecRef() int
	Call(string, interface{}, interface{}) error
}
