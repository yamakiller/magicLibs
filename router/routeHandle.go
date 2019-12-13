package router

import "github.com/yamakiller/magicLibs/st/hash"

func NewHandle(addr string, replicas int) *RouteHandle {
	return &RouteHandle{_addr: addr, _ctrls: hash.New(replicas)}
}

type RouteHandle struct {
	_delete func(ctrl IRouteCtrl)
	_addr   string
	_ctrls  *hash.Map
}

func (slf *RouteHandle) withDelete(f func(IRouteCtrl)) {
	slf._delete = f
}

func (slf *RouteHandle) getAddr() string {
	return slf._addr
}

func (slf *RouteHandle) isempty() bool {
	slf._ctrls.Lock()
	defer slf._ctrls.Unlock()
	return slf._ctrls.IsEmpty()
}

func (slf *RouteHandle) register(key string, ctrl IRouteCtrl) {
	slf._ctrls.Lock()
	defer slf._ctrls.Unlock()

	slf._ctrls.Add(key, ctrl)

	ctrl.IncRef()
	ctrl.IncRef()
}

func (slf *RouteHandle) unregister(key string) {
	slf._ctrls.Lock()
	v := slf._ctrls.Remove(key)
	if v != nil {
		n := v.(IRouteCtrl).DecRef()
		slf._ctrls.Unlock()
		if n <= 0 {
			slf._delete(v.(IRouteCtrl))
		}
		return
	}
	slf._ctrls.Unlock()
}

func (slf *RouteHandle) getKeys() []string {
	slf._ctrls.Lock()
	defer slf._ctrls.Unlock()
	return slf._ctrls.GetKeys()
}

func (slf *RouteHandle) grap(key string) (IRouteCtrl, error) {
	slf._ctrls.Lock()
	defer slf._ctrls.Unlock()
	ctrl, err := slf._ctrls.Get(key)
	if err != nil {
		return nil, err
	}
	ctrl.(IRouteCtrl).IncRef()
	return ctrl.(IRouteCtrl), nil
}

func (slf *RouteHandle) release(ctrl IRouteCtrl) {
	slf._ctrls.Lock()
	n := ctrl.DecRef()
	slf._ctrls.Unlock()
	if n <= 0 {
		slf._delete(ctrl)
	}
}
