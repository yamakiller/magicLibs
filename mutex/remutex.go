package mutex

import (
	"sync"

	"github.com/yamakiller/magicLibs/util"
)

//ReMutex desc
//@struct ReMutex desc: Reentrant mutex
type ReMutex struct {
	mutex *sync.Mutex
	owner int
	count int
}

//Width desc
//@method Width desc: Sync lock association reentrant lock
//@param (*sync.Mutex) mutex object
func (slf *ReMutex) Width(m *sync.Mutex) {
	slf.mutex = m
}

//Lock desc
//@method Lock desc: locking
func (slf *ReMutex) Lock() {
	me := util.GetCurrentGoroutineID()
	if slf.owner == me {
		slf.count++
		return
	}

	slf.mutex.Lock()
}

//Unlock desc
//@method Unlock desc : unlocking
func (slf *ReMutex) Unlock() {
	util.Assert(slf.owner == util.GetCurrentGoroutineID(), "illegalMonitorStateError")
	if slf.count > 0 {
		slf.count--
	} else {
		slf.mutex.Unlock()
	}
}
