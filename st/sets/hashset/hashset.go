package hashset

import (
	"fmt"
	"strings"
)

//Set desc
//@struct Set desc: holds elements in go`s native map
type Set struct {
	items map[interface{}]struct{}
}

//Push desc
//@method Push desc: adds the es (one or more) to the set
//@param  (...interface{}) insert elements
func (s *Set) Push(es ...interface{}) {
	for _, it := range es {
		s.items[it] = struct{}{}
	}
}

//PushAll desc
//@method PushAll desc: adds the st(set in element) to the set.
//@param (*Set) sets
func (s *Set) PushAll(st *Set) {
	for _, it := range st.items {
		s.items[it] = struct{}{}
	}
}

//Retain desc
//@method Retain desc: retain the es (one or more) to the set.
//@param (...interface{}) elements
func (s *Set) Retain(eds ...interface{}) {
	vs := make(map[interface{}]struct{})
	for _, it := range eds {
		if v, ok := s.items[it]; ok {
			vs[it] = v
		}
	}

	s.items = vs
}

//RetainAll desc
//@method RetainAll desc: retain the st(set in element) to the set.
//@param (*Set) sets
func (s *Set) RetainAll(st *Set) {
	vs := make(map[interface{}]struct{})
	for _, it := range st.items {
		if v, ok := s.items[it]; ok {
			vs[it] = v
		}
	}
	s.items = vs
}

//Erase desc
//@method Erase desc: removes the es (one or more) from the set
//@param  (...interface{}) elements
func (s *Set) Erase(es ...interface{}) {
	for _, it := range es {
		delete(s.items, it)
	}
}

//EraseAll desc
//@method EraseAll desc: removes this st(set in element) from the set
//@param (*Set) sets
func (s *Set) EraseAll(st *Set) {
	for _, it := range st.items {
		delete(s.items, it)
	}
}

//Contains desc
//@method Contains desc: check if es (one or more) are present in the set.
//@param  (...interface{}) elements
//@return (bool)
func (s *Set) Contains(es ...interface{}) bool {
	for _, it := range es {
		if _, cs := s.items[it]; !cs {
			return false
		}
	}
	return true
}

//Size desc
//@method Size desc: returns number of elements within the set.
//@return (int) size
func (s *Set) Size() int {
	return len(s.items)
}

//IsEmpty desc
//@method IsEmpty desc: returns true if set does not contain any elements.
//@param (bool)
func (s *Set) IsEmpty() bool {
	return s.Size() == 0
}

//Clear desc
//@method Clear desc: clears all values in the set.
func (s *Set) Clear() {
	s.items = make(map[interface{}]struct{})
}

//Values desc
//@method Values desc: returns all items in the set.
//@return ([]interface{})
func (s *Set) Values() []interface{} {
	vs := make([]interface{}, s.Size())
	icnt := 0
	for it := range s.items {
		vs[icnt] = it
		icnt++
	}
	return vs
}

//String desc
//@method String desc: Returns a string
//@return (string)
func (s *Set) String() string {
	str := "HashSet\n"
	items := []string{}
	for k := range s.items {
		items = append(items, fmt.Sprintf("%v", k))
	}
	str += strings.Join(items, ", ")
	return str
}
