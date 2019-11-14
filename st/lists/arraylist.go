package lists

import (
	"fmt"
	"strings"
)

//List desc
//@struct List desc: holds the elements in a slice
type List struct {
	es   []interface{}
	size int
}

const (
	growthFactor = float32(2.0)
	shrinkFactor = float32(0.25)
)

//NewList desc
//@method NewList desc: instantiates a new list and adds the passed values, if any, to the list
//@param (...interface{}) elements
//@return (*List)
func NewList(values ...interface{}) *List {
	list := &List{}
	if len(values) > 0 {
		list.Add(values...)
	}
	return list
}

//Add desc
//@method Add desc: appends a value at the end of the list
//@param (...interface{}) elements
func (list *List) Add(values ...interface{}) {
	list.growBy(len(values))
	for _, value := range values {
		list.es[list.size] = value
		list.size++
	}
}

//Get desc
//@method Get desc: returns the element at index.
//@param  (int) index
//@return (interface{}) Returns value
//@return (bool)  Second return parameter is true if index is within bounds of the array and array is not empty, otherwise false.
func (list *List) Get(index int) (interface{}, bool) {

	if !list.withinRange(index) {
		return nil, false
	}

	return list.es[index], true
}

//Remove desc
//@method Remove desc: removes the element at the given index from the list.
//@param (int) index
func (list *List) Remove(index int) {

	if !list.withinRange(index) {
		return
	}

	list.es[index] = nil                              // cleanup reference
	copy(list.es[index:], list.es[index+1:list.size]) // shift to the left by one (slow operation, need ways to optimize this)
	list.size--

	list.shrink()
}

//Contains desc
//@method Contains desc: checks if elements (one or more) are present in the set.
// All elements have to be present in the set for the method to return true.
// Performance time complexity of n^2.
// Returns true if no arguments are passed at all, i.e. set is always super-set of empty set.
//@param  (...interface{}) elements
//@return (bool)
func (list *List) Contains(values ...interface{}) bool {

	for _, searchValue := range values {
		found := false
		for _, element := range list.es {
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

//Values desc
//@method Values desc: returns all elements in the list.
//@return ([]interface{})
func (list *List) Values() []interface{} {
	newElements := make([]interface{}, list.size, list.size)
	copy(newElements, list.es[:list.size])
	return newElements
}

//IndexOf desc
//@method IndexOf desc: returns index of provided element
//@param  (interface{}) element
//@return (int) index
func (list *List) IndexOf(value interface{}) int {
	if list.size == 0 {
		return -1
	}
	for index, element := range list.es {
		if element == value {
			return index
		}
	}
	return -1
}

//IsEmpty desc
//@method IsEmpty desc: returns true if list does not contain any elements.
//@return (bool)
func (list *List) IsEmpty() bool {
	return list.size == 0
}

//Size desc
//@method Size desc: returns number of elements within the list.
//@return (int) size
func (list *List) Size() int {
	return list.size
}

//Clear desc
//@method Clear desc: removes all elements from the list.
func (list *List) Clear() {
	list.size = 0
	list.es = []interface{}{}
}

//Swap desc
//@method Swap desc: swaps the two values at the specified positions.
//@param (int)
//@param (int)
func (list *List) Swap(i, j int) {
	if list.withinRange(i) && list.withinRange(j) {
		list.es[i], list.es[j] = list.es[j], list.es[i]
	}
}

//Insert desc
//@method Insert desc: inserts values at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
//@param (int) index
//@param (...interface{}) values
// Does not do anything if position is negative or bigger than list's size
// Note: position equal to list's size is valid, i.e. append.
func (list *List) Insert(index int, values ...interface{}) {

	if !list.withinRange(index) {
		// Append
		if index == list.size {
			list.Add(values...)
		}
		return
	}

	l := len(values)
	list.growBy(l)
	list.size += l
	copy(list.es[index+l:], list.es[index:list.size-l])
	copy(list.es[index:], values)
}

//Set desc
//@method Set desc:the value at specified index
//@param (int) index
//@param (interface{}) value
// Does not do anything if position is negative or bigger than list's size
// Note: position equal to list's size is valid, i.e. append.
func (list *List) Set(index int, value interface{}) {

	if !list.withinRange(index) {
		// Append
		if index == list.size {
			list.Add(value)
		}
		return
	}

	list.es[index] = value
}

//String desc
//@method String desc: returns a string representation of container
//@return (string)
func (list *List) String() string {
	str := "ArrayList\n"
	values := []string{}
	for _, value := range list.es[:list.size] {
		values = append(values, fmt.Sprintf("%v", value))
	}
	str += strings.Join(values, ", ")
	return str
}

// Check that the index is within bounds of the list
func (list *List) withinRange(index int) bool {
	return index >= 0 && index < list.size
}

func (list *List) resize(cap int) {
	newElements := make([]interface{}, cap, cap)
	copy(newElements, list.es)
	list.es = newElements
}

// Expand the array if necessary, i.e. capacity will be reached if we add n elements
func (list *List) growBy(n int) {
	// When capacity is reached, grow by a factor of growthFactor and add number of elements
	currentCapacity := cap(list.es)
	if list.size+n >= currentCapacity {
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
	currentCapacity := cap(list.es)
	if list.size <= int(float32(currentCapacity)*shrinkFactor) {
		list.resize(list.size)
	}
}
