package test

import (
	"fmt"
	"sync/atomic"
	"testing"
	"unsafe"
)

type testValue struct {
	a int32
	b int32
}

func TestLoad(t *testing.T) {
	onePtr := &testValue{a: 1, b: 2}

	p := (*testValue)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&onePtr))))
	fmt.Println("p:", p, ",", onePtr)
	fmt.Println((p))
	//b := unsafe.Pointer(onePtr)
	//a := atomic.LoadPointer(&b)

}
