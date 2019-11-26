package mutex

import (
	"sync"

	"github.com/yamakiller/magicLibs/util"
)

//ReMutex desc
//@struct ReMutex desc: Reentrant mutex
type ReMutex struct {
	_mutex *sync.Mutex
	_owner int
	_count int
}

//Width desc
//@method Width desc: Sync lock association reentrant lock
//@param (*sync.Mutex) mutex object
func (slf *ReMutex) Width(m *sync.Mutex) {
	slf._mutex = m
}

//Lock desc
//@method Lock desc: locking
func (slf *ReMutex) Lock() {
	me := util.GetCurrentGoroutineID()
	if slf._owner == me {
		slf._count++
		return
	}

	slf._mutex.Lock()
}

//Unlock desc
//@method Unlock desc : unlocking
func (slf *ReMutex) Unlock() {
	util.Assert(slf._owner == util.GetCurrentGoroutineID(), "illegalMonitorStateError")
	if slf._count > 0 {
		slf._count--
	} else {
		slf._mutex.Unlock()
	}
}
