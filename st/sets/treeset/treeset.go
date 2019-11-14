package treeset

import (
	"fmt"
	"strings"

	rbt "github.com/yamakiller/magicNet/st/trees/redblacktree"
)

//Set desc
//@struct Set desc: holds elements in a red-black tree
type Set struct {
	tree *rbt.Tree
}

//Push desc
//@method Push desc: adds the es (one or more) to the set
//@param (...interface{}) elements
func (s *Set) Push(es ...interface{}) {
	for _, it := range es {
		s.tree.Insert(it, struct{}{})
	}
}

//PushAll desc
//@method PushAll desc: adds the st(set in element) to the set.
//@param (*Set) sets
func (s *Set) PushAll(st *Set) {
	it := st.tree.Iterator()
	for i := 0; it.Next(); i++ {
		if it.It() != nil {
			s.tree.Insert(it.Key(), it.Value())
		}
	}
}

//Retain desc
//@method Retain desc: retain the es (one or more) to the set.
//@param (...interface{}) elements
func (s *Set) Retain(eds ...interface{}) {
	var vs []interface{}
	var ic int
	for _, it := range eds {
		if _, ok := s.tree.Get(it); ok {
			vs = append(vs, it)
			ic++
		}
	}

	s.tree.Clear()
	for i := 0; i < ic; i++ {
		s.tree.Insert(vs[i], struct{}{})
	}
	vs = nil
}

//RetainAll desc
//@method RetainAll desc: retain the st(set in element) to the set.
//@param (*Set) sets
func (s *Set) RetainAll(st *Set) {
	var vs []interface{}
	var ic int

	it := st.tree.Iterator()
	for i := 0; it.Next(); i++ {
		if it.It() != nil {
			if _, ok := s.tree.Get(it.Key()); ok {
				vs = append(vs, it)
				ic++
			}
		}
	}

	s.tree.Clear()
	for i := 0; i < ic; i++ {
		s.tree.Insert(vs[i], struct{}{})
	}
	vs = nil
}

//Erase desc
//@method Erase desc: removes the es (one or more) from the set
//@param (...interface{}) elements
func (s *Set) Erase(es ...interface{}) {
	for _, it := range es {
		s.tree.Erase(it)
	}
}

//EraseAll desc
//@method EraseAll desc: removes this st(set in element) from the set
//@param (*Set) sets
func (s *Set) EraseAll(st *Set) {
	it := st.tree.Iterator()
	for i := 0; it.Next(); i++ {
		if it.It() != nil {
			s.tree.Erase(it.Key())
		}
	}
}

//Contains desc
//@method Contains desc: check if es (one or more) are present in the set.
//@param  (...interface{}) elements
//@return (bool)
func (s *Set) Contains(es ...interface{}) bool {
	for _, it := range es {
		if _, cs := s.tree.Get(it); !cs {
			return false
		}
	}
	return true
}

//Size desc
//@method Size desc: Returns number of elements within the set.
//@return (int) size
func (s *Set) Size() int {
	return s.tree.Size()
}

//IsEmpty desc
//@method IsEmpty desc: Returns true if set does not contain any elements.
//@return (bool)
func (s *Set) IsEmpty() bool {
	return s.Size() == 0
}

//Clear desc
//@method Clear desc: clears all values in the set.
func (s *Set) Clear() {
	s.tree.Clear()
}

//Values desc
//@method Values desc: Returns all items in the set.
//@return ([]interface{})
func (s *Set) Values() []interface{} {
	return s.tree.Keys()
}

//String desc
//@method String desc: Returns a string
func (s *Set) String() string {
	str := "TreeSet\n"
	items := []string{}
	for _, v := range s.tree.Keys() {
		items = append(items, fmt.Sprintf("%v", v))
	}
	str += strings.Join(items, ", ")
	return str
}
