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
func New(del func(p IRouteCtrl),
	replicas int) *RouteGroup {
	return &RouteGroup{_delete: del,
		_replicas: replicas,
		_rmaps:    make(map[string]*RouteHandle)}
}

//RouteGroup route group
type RouteGroup struct {
	_delete   func(p IRouteCtrl)
	_replicas int
	_serkey   int
	_rmaps    map[string]*RouteHandle
	_sync     sync.RWMutex
}

func (slf *RouteGroup) getAddrs() []string {
	i := 0
	slf._sync.Lock()
	addrs := make([]string, len(slf._rmaps))
	for k := range slf._rmaps {
		addrs[i] = k
		i++
	}
	slf._sync.Unlock()
	return addrs
}

//Size Return Group size
func (slf *RouteGroup) Size() int {
	slf._sync.Lock()
	defer slf._sync.Unlock()
	return len(slf._rmaps)
}

//WithReplicas Set replicas
func (slf *RouteGroup) WithReplicas(replicas int) {
	slf._replicas = replicas
}

//IsExist addrss is exist
func (slf *RouteGroup) IsExist(addr, key string) bool {
	slf._sync.RLock()
	defer slf._sync.RUnlock()
	if h, ok := slf._rmaps[addr]; ok {
		return h.isexist(key)
	}
	return false
}

//Register register router address and control
func (slf *RouteGroup) Register(addr, key string, ctrl IRouteCtrl) {
	slf._sync.Lock()
	defer slf._sync.Unlock()
	if h, ok := slf._rmaps[addr]; ok {
		h.register(key, ctrl)
		h.release(ctrl)
	} else {
		h := NewHandle(addr, slf._replicas)
		h.withDelete(slf._delete)
		h.register(key, ctrl)
		slf._rmaps[addr] = h
		h.release(ctrl)
	}
}

//UnRegister unregister a router
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

//Call Address Call remote function
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

//Shutdown shutdown route system
func (slf *RouteGroup) Shutdown() {
	for {
		if slf.Size() <= 0 {
			break
		}

		addrs := slf.getAddrs()
		for _, k := range addrs {
			slf._sync.Lock()
			if h, ok := slf._rmaps[k]; ok {
				delete(slf._rmaps, k)
				slf._sync.Unlock()
				h.shutdown()
				continue
			}
			slf._sync.Unlock()
		}
	}
}
