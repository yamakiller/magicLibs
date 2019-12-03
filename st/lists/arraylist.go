package lists

import (
	"fmt"
	"strings"
)

//List doc
//@Struct List @Summary holds the elements in a slice
type List struct {
	_es   []interface{}
	_size int
}

const (
	growthFactor = float32(2.0)
	shrinkFactor = float32(0.25)
)

//NewList doc
//@Method NewList @Summary instantiates a new list and adds the passed values, if any, to the list
//@Param (...interface{}) elements
//@Return (*List)
func NewList(values ...interface{}) *List {
	list := &List{}
	if len(values) > 0 {
		list.Add(values...)
	}
	return list
}

//Add doc
//@Method Add @Summary appends a value at the end of the list
//@Param (...interface{}) elements
func (list *List) Add(values ...interface{}) {
	list.growBy(len(values))
	for _, value := range values {
		list._es[list._size] = value
		list._size++
	}
}

//Get doc
//@Method Get @Summary returns the element at index.
//@Param  (int) index
//@Return (interface{}) Returns value
//@Return (bool)  Second return parameter is true if index is within bounds of the array and array is not empty, otherwise false.
func (list *List) Get(index int) (interface{}, bool) {

	if !list.withinRange(index) {
		return nil, false
	}

	return list._es[index], true
}

//Remove doc
//@Method Remove @Summary removes the element at the given index from the list.
//@Param (int) index
func (list *List) Remove(index int) {

	if !list.withinRange(index) {
		return
	}

	list._es[index] = nil                                // cleanup reference
	copy(list._es[index:], list._es[index+1:list._size]) // shift to the left by one (slow operation, need ways to optimize this)
	list._size--

	list.shrink()
}

//Contains doc
//@Method Contains @Summary checks if elements (one or more) are present in the set.
// All elements have to be present in the set for the method to return true.
// Performance time complexity of n^2.
// Returns true if no arguments are passed at all, i.e. set is always super-set of empty set.
//@Param  (...interface{}) elements
//@Return (bool)
func (list *List) Contains(values ...interface{}) bool {

	for _, searchValue := range values {
		found := false
		for _, element := range list._es {
			if element == searchValue {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

//Values doc
//@Method Values @Summary returns all elements in the list.
//@Return ([]interface{})
func (list *List) Values() []interface{} {
	newElements := make([]interface{}, list._size, list._size)
	copy(newElements, list._es[:list._size])
	return newElements
}

//IndexOf doc
//@Method IndexOf @Summary returns index of provided element
//@Param  (interface{}) element
//@Return (int) index
func (list *List) IndexOf(value interface{}) int {
	if list._size == 0 {
		return -1
	}
	for index, element := range list._es {
		if element == value {
			return index
		}
	}
	return -1
}

//IsEmpty doc
//@Method IsEmpty @Summary returns true if list does not contain any elements.
//@Return (bool)
func (list *List) IsEmpty() bool {
	return list._size == 0
}

//Size doc
//@Method Size @Summary returns number of elements within the list.
//@Return (int) size
func (list *List) Size() int {
	return list._size
}

//Clear doc
//@Method Clear @Summary removes all elements from the list.
func (list *List) Clear() {
	list._size = 0
	list._es = []interface{}{}
}

//Swap doc
//@Method Swap @Summary swaps the two values at the specified positions.
//@Param (int)
//@Param (int)
func (list *List) Swap(i, j int) {
	if list.withinRange(i) && list.withinRange(j) {
		list._es[i], list._es[j] = list._es[j], list._es[i]
	}
}

//Insert doc
//@Method Insert @Summary inserts values at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
//@Param (int) index
//@Param (...interface{}) values
// Does not do anything if position is negative or bigger than list's size
// Note: position equal to list's size is valid, i.e. append.
func (list *List) Insert(index int, values ...interface{}) {

	if !list.withinRange(index) {
		// Append
		if index == list._size {
			list.Add(values...)
		}
		return
	}

	l := len(values)
	list.growBy(l)
	list._size += l
	copy(list._es[index+l:], list._es[index:list._size-l])
	copy(list._es[index:], values)
}

//Set doc
//@Method Set @Summarythe value at specified index
//@Param (int) index
//@Param (interface{}) value
// Does not do anything if position is negative or bigger than list's size
// Note: position equal to list's size is valid, i.e. append.
func (list *List) Set(index int, value interface{}) {

	if !list.withinRange(index) {
		// Append
		if index == list._size {
			list.Add(value)
		}
		return
	}

	list._es[index] = value
}

//String doc
//@Method String @Summary returns a string representation of container
//@Return (string)
func (list *List) String() string {
	str := "ArrayList\n"
	values := []string{}
	for _, value := range list._es[:list._size] {
		values = append(values, fmt.Sprintf("%v", value))
	}
	str += strings.Join(values, ", ")
	return str
}

// Check that the index is within bounds of the list
func (list *List) withinRange(index int) bool {
	return index >= 0 && index < list._size
}

func (list *List) resize(cap int) {
	newElements := make([]interface{}, cap, cap)
	copy(newElements, list._es)
	list._es = newElements
}

// Expand the array if necessary, i.e. capacity will be reached if we add n elements
func (list *List) growBy(n int) {
	// When capacity is reached, grow by a factor of growthFactor and add number of elements
	currentCapacity := cap(list._es)
	if list._size+n >= currentCapacity {
		newCapacity := int(growthFactor * float32(currentCapacity+n))
		list.resize(newCapacity)
	}
}

// Shrink the array if necessary, i.e. when size is shrinkFactor percent of current capacity
func (list *List) shrink() {
	if shrinkFactor == 0.0 {
		return
	}
	// Shrink when size is at shrinkFactor * capacity
	currentCapacity := cap(list._es)
	if list._size <= int(float32(currentCapacity)*shrinkFactor) {
		list.resize(list._size)
	}
}
