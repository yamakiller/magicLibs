package treeset

import (
	"fmt"
	"strings"

	rbt "github.com/yamakiller/magicLibs/st/trees/redblacktree"
)

//Set doc
//@Struct Set @Summary holds elements in a red-black tree
type Set struct {
	_tree *rbt.Tree
}

//Push doc
//@Method Push @Summary adds the es (one or more) to the set
//@Param (...interface{}) elements
func (s *Set) Push(es ...interface{}) {
	for _, it := range es {
		s._tree.Insert(it, struct{}{})
	}
}

//PushAll doc
//@Method PushAll @Summary adds the st(set in element) to the set.
//@Param (*Set) sets
func (s *Set) PushAll(st *Set) {
	it := st._tree.Iterator()
	for i := 0; it.Next(); i++ {
		if it.It() != nil {
			s._tree.Insert(it.Key(), it.Value())
		}
	}
}

//Retain doc
//@Method Retain @Summary retain the es (one or more) to the set.
//@Param (...interface{}) elements
func (s *Set) Retain(eds ...interface{}) {
	var vs []interface{}
	var ic int
	for _, it := range eds {
		if _, ok := s._tree.Get(it); ok {
			vs = append(vs, it)
			ic++
		}
	}

	s._tree.Clear()
	for i := 0; i < ic; i++ {
		s._tree.Insert(vs[i], struct{}{})
	}
	vs = nil
}

//RetainAll doc
//@Method RetainAll @Summary retain the st(set in element) to the set.
//@Param (*Set) sets
func (s *Set) RetainAll(st *Set) {
	var vs []interface{}
	var ic int

	it := st._tree.Iterator()
	for i := 0; it.Next(); i++ {
		if it.It() != nil {
			if _, ok := s._tree.Get(it.Key()); ok {
				vs = append(vs, it)
				ic++
			}
		}
	}

	s._tree.Clear()
	for i := 0; i < ic; i++ {
		s._tree.Insert(vs[i], struct{}{})
	}
	vs = nil
}

//Erase doc
//@Method Erase @Summary removes the es (one or more) from the set
//@Param (...interface{}) elements
func (s *Set) Erase(es ...interface{}) {
	for _, it := range es {
		s._tree.Erase(it)
	}
}

//EraseAll doc
//@Method EraseAll @Summary removes this st(set in element) from the set
//@Param (*Set) sets
func (s *Set) EraseAll(st *Set) {
	it := st._tree.Iterator()
	for i := 0; it.Next(); i++ {
		if it.It() != nil {
			s._tree.Erase(it.Key())
		}
	}
}

//Contains doc
//@Method Contains @Summary check if es (one or more) are present in the set.
//@Param  (...interface{}) elements
//@Return (bool)
func (s *Set) Contains(es ...interface{}) bool {
	for _, it := range es {
		if _, cs := s._tree.Get(it); !cs {
			return false
		}
	}
	return true
}

//Size doc
//@Method Size @Summary Returns number of elements within the set.
//@Return (int) size
func (s *Set) Size() int {
	return s._tree.Size()
}

//IsEmpty doc
//@Method IsEmpty @Summary Returns true if set does not contain any elements.
//@Return (bool)
func (s *Set) IsEmpty() bool {
	return s.Size() == 0
}

//Clear doc
//@Method Clear @Summary clears all values in the set.
func (s *Set) Clear() {
	s._tree.Clear()
}

//Values doc
//@Method Values @Summary Returns all items in the set.
//@Return ([]interface{})
func (s *Set) Values() []interface{} {
	return s._tree.Keys()
}

//String doc
//@Method String @Summary Returns a string
func (s *Set) String() string {
	str := "TreeSet\n"
	items := []string{}
	for _, v := range s._tree.Keys() {
		items = append(items, fmt.Sprintf("%v", v))
	}
	str += strings.Join(items, ", ")
	return str
}
