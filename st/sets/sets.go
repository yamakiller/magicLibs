package sets

import "github.com/yamakiller/magicLibs/st/containers"

//Set doc
//@Interface Set @Summary that all sets
type Set interface {
	Push(es ...interface{})
	PushAll(st *Set)
	Retain(eds ...interface{})
	RetainAll(st *Set)
	Erase(es ...interface{})
	EraseAll(st *Set)
	Contains(es ...interface{})

	containers.Container
}
