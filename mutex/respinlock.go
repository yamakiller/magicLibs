package mutex

import "github.com/yamakiller/magicLibs/util"

//ReSpinLock desc
//@struct ReSpinLock desc: Re-entrant spin lock
type ReSpinLock struct {
	mutex *SpinLock
	owner int
	count int
}

//Width desc
//@method Width desc : Spinlock association reentrant spin lock
//@param (*SpinLock) width spinlock object
func (slf *ReSpinLock) Width(m *SpinLock) {
	slf.mutex = m
}

//TryLock desc
//@method Trylock desc : Try to lock if you fail to get the lock return failure will not try again
//@return (bool)
func (slf *ReSpinLock) Trylock() bool {
	me := util.GetCurrentGoroutineID()
	if slf.owner == me {
		slf.count++
		return true
	}

	return slf.mutex.Trylock()
}

//Lock desc
//@method Lock desc: locking
func (slf *ReSpinLock) Lock() {
	me := util.GetCurrentGoroutineID()
	if slf.owner == me {
		slf.count++
		return
	}

	slf.mutex.Lock()
}

//Unlock desc
//@method Unlock desc: unlocking
func (slf *ReSpinLock) Unlock() {
	util.Assert(slf.owner == util.GetCurrentGoroutineID(), "illegalMonitorStateError")
	if slf.count > 0 {
		slf.count--
	} else {
		slf.mutex.Unlock()
	}
}
