package hashset

import (
	"fmt"
	"strings"
)

//Set doc
//@Struct Set @Summary holds elements in go`s native map
type Set struct {
	items map[interface{}]struct{}
}

//Push doc
//@Method Push @Summary adds the es (one or more) to the set
//@Param  (...interface{}) insert elements
func (s *Set) Push(es ...interface{}) {
	for _, it := range es {
		s.items[it] = struct{}{}
	}
}

//PushAll doc
//@Method PushAll @Summary adds the st(set in element) to the set.
//@Param (*Set) sets
func (s *Set) PushAll(st *Set) {
	for _, it := range st.items {
		s.items[it] = struct{}{}
	}
}

//Retain doc
//@Method Retain @Summary retain the es (one or more) to the set.
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

//RetainAll doc
//@Method RetainAll @Summary retain the st(set in element) to the set.
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

//Erase doc
//@Method Erase @Summary removes the es (one or more) from the set
//@Param  (...interface{}) elements
func (s *Set) Erase(es ...interface{}) {
	for _, it := range es {
		delete(s.items, it)
	}
}

//EraseAll doc
//@Method EraseAll @Summary removes this st(set in element) from the set
//@Param (*Set) sets
func (s *Set) EraseAll(st *Set) {
	for _, it := range st.items {
		delete(s.items, it)
	}
}

//Contains doc
//@Method Contains @Summary check if es (one or more) are present in the set.
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

//Size doc
//@Method Size @Summary returns number of elements within the set.
//@Return (int) size
func (s *Set) Size() int {
	return len(s.items)
}

//IsEmpty doc
//@Method IsEmpty @Summary returns true if set does not contain any elements.
//@Param (bool)
func (s *Set) IsEmpty() bool {
	return s.Size() == 0
}

//Clear doc
//@Method Clear @Summary clears all values in the set.
func (s *Set) Clear() {
	s.items = make(map[interface{}]struct{})
}

//Values doc
//@Method Values @Summary returns all items in the set.
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

//String doc
//@Method String @Summary Returns a string
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
