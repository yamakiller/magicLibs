package lists

import (
	"sync/atomic"
	"unsafe"
)

type node struct {
	_next *node
	_val  interface{}
}

//Queue desc
//@Struct Queue desc: Simple queue
//@Member (*node) header
//@Member (*node) tail
type Queue struct {
	_head, _tail *node
}

//NewQueue desc
//@Method NewQueue desc: Create a queue object
//@Return (*Queue) Queue object
func NewQueue() *Queue {
	q := &Queue{}
	stub := &node{}
	q._head = stub
	q._tail = stub
	return q
}

//Push desc
//@Method Push desc: Insert an Object into the queue
//@Param  (interface{}) insert value
func (slf *Queue) Push(t interface{}) {
	n := new(node)
	n._val = t
	prev := (*node)(atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&slf._head)), unsafe.Pointer(n)))
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&prev._next)), unsafe.Pointer(n))
}

//Pop desc
//@Method Pop desc: An object pops up in the re-queue
//@Return (interface{}) pop header elements
func (slf *Queue) Pop() interface{} {
	tail := slf._tail
	next := (*node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&tail._next))))
	if next != nil {
		slf._tail = next
		v := next._val
		next._val = nil
		return v
	}
	return nil
}

//IsEmpty desc
//@Method IsEmpty desc: Whether the queue is empty
//@Return (bool) null:true, not null:false
func (slf *Queue) IsEmpty() bool {
	tail := slf._tail
	next := (*node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&tail._next))))
	return next == nil
}
