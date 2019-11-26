package mutex

import "github.com/yamakiller/magicLibs/util"

//ReSpinLock desc
//@struct ReSpinLock desc: Re-entrant spin lock
type ReSpinLock struct {
	_mutex *SpinLock
	_owner int
	_count int
}

//Width desc
//@method Width desc : Spinlock association reentrant spin lock
//@param (*SpinLock) width spinlock object
func (slf *ReSpinLock) Width(m *SpinLock) {
	slf._mutex = m
}

//Trylock desc
//@method Trylock desc : Try to lock if you fail to get the lock return failure will not try again
//@return (bool)
func (slf *ReSpinLock) Trylock() bool {
	me := util.GetCurrentGoroutineID()
	if slf._owner == me {
		slf._count++
		return true
	}

	return slf._mutex.Trylock()
}

//Lock desc
//@method Lock desc: locking
func (slf *ReSpinLock) Lock() {
	me := util.GetCurrentGoroutineID()
	if slf._owner == me {
		slf._count++
		return
	}

	slf._mutex.Lock()
}

//Unlock desc
//@method Unlock desc: unlocking
func (slf *ReSpinLock) Unlock() {
	util.Assert(slf._owner == util.GetCurrentGoroutineID(), "illegalMonitorStateError")
	if slf._count > 0 {
		slf._count--
	} else {
		slf._mutex.Unlock()
	}
}
