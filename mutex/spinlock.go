package mutex

import "sync/atomic"

//SpinLock doc
//@Struct SpinLick @Summary spin lock
//@Member (uint32)
type SpinLock struct {
	_kernel uint32
}

//Trylock doc
//@Method Trylock @Summary try lock if unlock return false
//@Return (bool)
func (slf *SpinLock) Trylock() bool {
	return atomic.CompareAndSwapUint32(&slf._kernel, 0, 1)
}

//Lock doc
//@Method Lock @Summary locking
func (slf *SpinLock) Lock() {
	for !atomic.CompareAndSwapUint32(&slf._kernel, 0, 1) {
	}
}

//Unlock doc
//@Method Unlock @Summary unlocking
func (slf *SpinLock) Unlock() {
	atomic.StoreUint32(&slf._kernel, 0)
}
