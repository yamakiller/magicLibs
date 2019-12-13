package router

import (
	"errors"
	"fmt"
	"math"
	"sync"
)

var (
	//ErrControlNot error
	ErrControlNot = errors.New("control does not exist")
	//ErrAddressNot error
	ErrAddressNot = errors.New("route address does not exist")
)

//New new an Router Group
func New(del func(p IRouteCtrl), replicas int) *RouteGroup {
	return &RouteGroup{_delete: del, _replicas: replicas, _rmaps: make(map[string]*RouteHandle)}
}

type RouteGroup struct {
	_delete   func(p IRouteCtrl)
	_replicas int
	_serkey   int
	_rmaps    map[string]*RouteHandle

	_sync sync.RWMutex
}

func (slf *RouteGroup) WithReplicas(replicas int) {
	slf._replicas = replicas
}

func (slf *RouteGroup) Register(addr, key string, ctrl IRouteCtrl) {
	slf._sync.Lock()
	defer slf._sync.Unlock()
	if h, ok := slf._rmaps[addr]; ok {
		h.register(key, ctrl)
	} else {
		h := NewHandle(addr, slf._replicas)
		h.withDelete(slf._delete)
		h.register(key, ctrl)
		slf._rmaps[addr] = h
	}
}

func (slf *RouteGroup) UnRegister(addr, key string) {
	slf._sync.Lock()
	if h, ok := slf._rmaps[addr]; ok {
		slf._sync.Unlock()
		h.unregister(key)
		slf._sync.Lock()
		if h.isempty() {
			if _, ok := slf._rmaps[addr]; ok {
				delete(slf._rmaps, addr)
			}
		}
	}
	slf._sync.Unlock()
}

func (slf *RouteGroup) Call(addr, method string, param, ret interface{}) error {
	slf._sync.RLock()
	if _, ok := slf._rmaps[addr]; ok {
		slf._sync.RUnlock()
		return ErrAddressNot
	}

	slf._serkey = ((slf._serkey + 1) & math.MaxInt32)
	v := slf._rmaps[addr]
	slf._sync.RUnlock()
	h, err := v.grap(fmt.Sprintf("key%d", slf._serkey))
	if err == nil {
		return err
	}
	defer v.release(h)

	return h.Call(method, param, ret)
}
