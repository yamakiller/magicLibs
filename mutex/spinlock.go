package mutex

import "sync/atomic"

//SpinLock desc
//@Struct SpinLick desc: spin lock
//@Member (uint32)
type SpinLock struct {
	_kernel uint32
}

//Trylock desc
//@Method Trylock desc: try lock if unlock return false
//@Return (bool)
func (slf *SpinLock) Trylock() bool {
	return atomic.CompareAndSwapUint32(&slf._kernel, 0, 1)
}

//Lock desc
//@Method Lock desc: locking
func (slf *SpinLock) Lock() {
	for !atomic.CompareAndSwapUint32(&slf._kernel, 0, 1) {
	}
}

//Unlock desc
//@Method Unlock desc: unlocking
func (slf *SpinLock) Unlock() {
	atomic.StoreUint32(&slf._kernel, 0)
}
