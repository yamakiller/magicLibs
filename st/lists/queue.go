package lists

import (
	"sync/atomic"
	"unsafe"
)

type node struct {
	next *node
	val  interface{}
}

//Queue desc
//@struct Queue desc: Simple queue
//@member (*node) header
//@member (*node) tail
type Queue struct {
	head, tail *node
}

//NewQueue desc
//@method NewQueue desc: Create a queue object
//@return (*Queue) Queue object
func NewQueue() *Queue {
	q := &Queue{}
	stub := &node{}
	q.head = stub
	q.tail = stub
	return q
}

//Push desc
//@method Push desc: Insert an Object into the queue
//@param  (interface{}) insert value
func (slf *Queue) Push(t interface{}) {
	n := new(node)
	n.val = t
	prev := (*node)(atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&slf.head)), unsafe.Pointer(n)))
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&prev.next)), unsafe.Pointer(n))
}

//Pop desc
//@method Pop desc: An object pops up in the re-queue
//@return (interface{}) pop header elements
func (slf *Queue) Pop() interface{} {
	tail := slf.tail
	next := (*node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&tail.next))))
	if next != nil {
		slf.tail = next
		v := next.val
		next.val = nil
		return v
	}
	return nil
}

//IsEmpty desc
//@method IsEmpty desc: Whether the queue is empty
//@return (bool) null:true, not null:false
func (slf *Queue) IsEmpty() bool {
	tail := slf.tail
	next := (*node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&tail.next))))
	return next == nil
}
