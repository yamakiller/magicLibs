package hashset

import (
	"fmt"
	"strings"
)

//Set desc
//@Struct Set desc: holds elements in go`s native map
type Set struct {
	items map[interface{}]struct{}
}

//Push desc
//@Method Push desc: adds the es (one or more) to the set
//@Param  (...interface{}) insert elements
func (s *Set) Push(es ...interface{}) {
	for _, it := range es {
		s.items[it] = struct{}{}
	}
}

//PushAll desc
//@Method PushAll desc: adds the st(set in element) to the set.
//@Param (*Set) sets
func (s *Set) PushAll(st *Set) {
	for _, it := range st.items {
		s.items[it] = struct{}{}
	}
}

//Retain desc
//@Method Retain desc: retain the es (one or more) to the set.
//@Param (...interface{}) elements
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
//@Method RetainAll desc: retain the st(set in element) to the set.
//@Param (*Set) sets
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
//@Method Erase desc: removes the es (one or more) from the set
//@Param  (...interface{}) elements
func (s *Set) Erase(es ...interface{}) {
	for _, it := range es {
		delete(s.items, it)
	}
}

//EraseAll desc
//@Method EraseAll desc: removes this st(set in element) from the set
//@Param (*Set) sets
func (s *Set) EraseAll(st *Set) {
	for _, it := range st.items {
		delete(s.items, it)
	}
}

//Contains desc
//@Method Contains desc: check if es (one or more) are present in the set.
//@Param  (...interface{}) elements
//@Return (bool)
func (s *Set) Contains(es ...interface{}) bool {
	for _, it := range es {
		if _, cs := s.items[it]; !cs {
			return false
		}
	}
	return true
}

//Size desc
//@Method Size desc: returns number of elements within the set.
//@Return (int) size
func (s *Set) Size() int {
	return len(s.items)
}

//IsEmpty desc
//@Method IsEmpty desc: returns true if set does not contain any elements.
//@Param (bool)
func (s *Set) IsEmpty() bool {
	return s.Size() == 0
}

//Clear desc
//@Method Clear desc: clears all values in the set.
func (s *Set) Clear() {
	s.items = make(map[interface{}]struct{})
}

//Values desc
//@Method Values desc: returns all items in the set.
//@Return ([]interface{})
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
//@Method String desc: Returns a string
//@Return (string)
func (s *Set) String() string {
	str := "HashSet\n"
	items := []string{}
	for k := range s.items {
		items = append(items, fmt.Sprintf("%v", k))
	}
	str += strings.Join(items, ", ")
	return str
}
