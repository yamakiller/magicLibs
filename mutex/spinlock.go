package mutex

import "sync/atomic"

//SpinLock desc
//@struct SpinLick desc: spin lock
//@member (uint32)
type SpinLock struct {
	_kernel uint32
}

//Trylock desc
//@method Trylock desc: try lock if unlock return false
//@return (bool)
func (slf *SpinLock) Trylock() bool {
	return atomic.CompareAndSwapUint32(&slf._kernel, 0, 1)
}

//Lock desc
//@method Lock desc: locking
func (slf *SpinLock) Lock() {
	for !atomic.CompareAndSwapUint32(&slf._kernel, 0, 1) {
	}
}

//Unlock desc
//@method Unlock desc: unlocking
func (slf *SpinLock) Unlock() {
	atomic.StoreUint32(&slf._kernel, 0)
}
