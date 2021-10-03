package redblacktree

import (
	"fmt"

	"github.com/yamakiller/magicLibs/st/comparator"
)

const (
	colorRed   = true
	colorBlock = false

	alloterMalloc = true
	alloterFree   = false
)

// Node red-black Tree is Node
type Node struct {
	//node
	_parent, _left, _right *Node
	//property
	_key, _val interface{}
	_color     bool
}

func (n *Node) grandparent() *Node {
	if n != nil && n._parent != nil {
		return n._parent._parent
	}
	return nil
}

func (n *Node) uncle() *Node {
	if n == nil || n._parent == nil || n._parent._parent == nil {
		return nil
	}

	return n._parent.sibling()
}

func (n *Node) sibling() *Node {
	if n == nil || n._parent == nil {
		return nil
	}
	if n == n._parent._left {
		return n._parent._right
	}
	return n._parent._left
}

func (n *Node) maximumNode() *Node {
	if n == nil {
		return nil
	}
	for n._right != nil {
		n = n._right
	}
	return n
}

func (n *Node) String() string {
	return fmt.Sprintf("%v", n._key)
}

// Alloter Tree alloter
type Alloter struct {
	A func(k, v interface{}, c bool) *Node
	F func(p *Node)
}

// Tree red-black tree
type Tree struct {
	_root    *Node
	_size    int
	_compare comparator.Comparator
	_alloter *Alloter
}

// NewTree instantiates a red-black tree with the custom comparator
func NewTree(comparator comparator.Comparator, alloc *Alloter) *Tree {
	return &Tree{_compare: comparator, _alloter: alloc}
}

// NewTreeIntComparator instantiates a red-black tree with the IntComparator, i.e. keys are of type int.
func NewTreeIntComparator(alloc *Alloter) *Tree {
	return &Tree{_compare: comparator.IntComparator, _alloter: alloc}
}

// NewTreeStringComparator instantiates a red-black tree with the IntComparator, i.e. keys are of type string.
func NewTreeStringComparator(alloc *Alloter) *Tree {
	return &Tree{_compare: comparator.StringComparator, _alloter: alloc}
}

var (
	defaultAlloter = Alloter{A: func(k, v interface{}, c bool) *Node {
		return &Node{_key: k, _val: v, _color: c}
	},
		F: func(p *Node) {

		},
	}
)

func (t *Tree) newNode(k, v interface{}, c bool) *Node {
	if t._alloter == nil {
		return defaultAlloter.A(k, v, c)
	}

	return t._alloter.A(k, v, c)
}

func (t *Tree) freeNode(n *Node) {
	if t._alloter == nil {
		defaultAlloter.F(n)
		return
	}

	t._alloter.F(n)
	return
}

func (t *Tree) isRed(n *Node) bool {
	return n != nil && n._color
}

// Empty returns true if tree does not contain any nodes
func (t *Tree) Empty() bool {
	return (t._size == 0)
}

// Size returns number of nodes in the tree
func (t *Tree) Size() int {
	return t._size
}

// Clear removes all nodes from the tree
func (t *Tree) Clear() {
	//! With memory management, nodes need to be released one by one
	if t._alloter != nil {
		it := t.Iterator()
		for i := 0; it.Next(); i++ {
			if it.It() != nil {
				t.freeNode(it.It())
			}
		}
	}

	t._root = nil
	t._size = 0
}

func (t *Tree) String() string {
	str := "RedBlackTree\n"
	if !t.Empty() {
		output(t._root, "", true, &str)
	}
	return str
}

func output(n *Node, prefix string, isTail bool, str *string) {
	if n._right != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		output(n._right, newPrefix, false, str)
	}
	*str += prefix
	if isTail {
		*str += "└── "
	} else {
		*str += "┌── "
	}
	*str += n.String() + "\n"
	if n._left != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		output(n._left, newPrefix, true, str)
	}
}

func (t *Tree) getDepth() int {
	var depthCall func(n *Node) int
	depthCall = func(n *Node) int {
		if n == nil {
			return 0
		}

		if n._left == nil && n._right == nil {
			return 1
		}

		ld := depthCall(n._left)
		rd := depthCall(n._right)

		if ld > rd {
			return ld + 1
		}
		return rd + 1
	}

	return depthCall(t._root)
}

func (t *Tree) leftRotate(n *Node) {
	if n._right == nil {
		return
	}

	r := n._right
	n._right = r._left
	if n._right != nil {
		n._right._parent = n
	}

	r._parent = n._parent
	if n._parent == nil {
		t._root = r
	} else {
		if n._parent._left == n {
			n._parent._left = r
		} else {
			n._parent._right = r
		}
	}
	r._left = n
	n._parent = r
}

func (t *Tree) rightRotate(n *Node) {
	if n._left == nil {
		return
	}

	l := n._left
	n._left = l._right
	if n._left != nil {
		n._left._parent = n
	}

	l._parent = n._parent
	if n._parent == nil {
		t._root = l
	} else {
		if n._parent._left == n {
			n._parent._left = l
		} else {
			n._parent._right = l
		}
	}
	l._right = n
	n._parent = l
}

func (t *Tree) lookup(key interface{}) *Node {
	n := t._root
	for n != nil {
		compare := t._compare(key, n._key)
		if compare > 0 {
			n = n._right
		} else if compare < 0 {
			n = n._left
		} else {
			return n
		}
	}
	return nil
}

// Left returns the left-most (min) node or nil if tree is empty
func (t *Tree) left() *Node {
	var parent *Node
	cur := t._root
	for cur != nil {
		parent = cur
		cur = cur._left
	}
	return parent
}

// Right returns the right-most (max) node or nil if tree is empty
func (t *Tree) right() *Node {
	var parent *Node
	cur := t._root
	for cur != nil {
		parent = cur
		cur = cur._right
	}
	return parent
}

// Keys returns all keys
func (t *Tree) Keys() []interface{} {
	keys := make([]interface{}, t._size)
	it := t.Iterator()
	for i := 0; it.Next(); i++ {
		keys[i] = it.Key()
	}
	return keys
}

// Values returns all values
func (t *Tree) Values() []interface{} {
	vals := make([]interface{}, t._size)
	it := t.Iterator()
	for i := 0; it.Next(); i++ {
		vals[i] = it.Value()
	}
	return vals
}

// Get sreach the node in thre tree by key and returns its value or nil if key is not found in tree
func (t *Tree) Get(k interface{}) (interface{}, bool) {
	n := t.lookup(k)
	if n != nil {
		return n._val, true
	}
	return nil, false
}

// Insert inserts node into the tree
func (t *Tree) Insert(k, v interface{}) {
	var n *Node
	if t.Empty() {
		t._root = t.newNode(k, v, colorRed)
		return
	}

	cur := t._root
	for {
		if t._compare(k, cur._key) > 0 {
			if cur._right == nil {
				cur._right = t.newNode(k, v, colorRed)
				n = cur._right
				break
			}

			cur = cur._right
			continue
		} else if t._compare(k, cur._key) < 0 {
			if cur._left == nil {
				cur._left = t.newNode(k, v, colorRed)
				n = cur._left
				break
			}

			cur = cur._left
			continue
		}

		cur._key = k
		cur._val = v
		return
	}
	n._parent = cur

	t.insertFixup1(n)
	t._size++
}

// Erase remove the node from the tree by key.
func (t *Tree) Erase(k interface{}) {
	var c *Node
	n := t.lookup(k)
	if n == nil {
		return
	}
	if n._left != nil && n._right != nil {
		pred := n._left.maximumNode()
		n._key = pred._key
		n._val = pred._val
		n = pred
	}
	if n._left == nil || n._right == nil {
		if n._right == nil {
			c = n._left
		} else {
			c = n._right
		}
		if n._color == colorBlock {
			n._color = t.isRed(c)
			t.eraseFixup1(n)
		}
		if n._parent == nil {
			t._root = c
		} else {
			if n == n._left._parent {
				n._parent._left = c
			} else {
				n._parent._right = c
			}
		}
		if c != nil {
			c._parent = n._parent
		}
		if n._parent == nil && c != nil {
			c._color = colorBlock
		}
	}
	t._size--
	if n._parent == nil &&
		n._left == nil &&
		n._right == nil {
		// ?  Logging and testing. Can it be performed here?
		t.freeNode(n)
	}
}

func (t *Tree) insertFixup1(n *Node) {
	if n._parent == nil {
		n._color = colorBlock
	} else {
		t.insertFixup2(n)
	}
}

func (t *Tree) insertFixup2(n *Node) {
	if !t.isRed(n) {
		return
	}
	t.insertFixup3(n)
}

func (t *Tree) insertFixup3(n *Node) {
	uncle := n.uncle()
	if t.isRed(uncle) {
		n._parent._color, uncle._color = colorBlock, colorBlock
		n.grandparent()._color = colorRed
		t.insertFixup1(n.grandparent())
	} else {
		t.insertFixup4(n)
	}
}

func (t *Tree) insertFixup4(n *Node) {
	grandparent := n.grandparent()
	if n == n._parent._right && n._parent == grandparent._left {
		t.leftRotate(n._parent)
		n = n._left
	} else if n == n._parent._left && n._parent == grandparent._right {
		t.rightRotate(n._parent)
		n = n._right
	}
	t.insertFixup5(n)
}

func (t *Tree) insertFixup5(n *Node) {
	n._parent._color = colorBlock
	grandparent := n.grandparent()
	grandparent._color = colorRed
	if n == n._parent._left && n._parent == grandparent._left {
		t.rightRotate(grandparent)
	} else if n == n._parent._right && n._parent == grandparent._right {
		t.leftRotate(grandparent)
	}
}

func (t *Tree) eraseFixup1(n *Node) {
	if n._parent == nil {
		return
	}
	t.eraseFixup2(n)
}

func (t *Tree) eraseFixup2(n *Node) {
	sibling := n.sibling()
	if t.isRed(sibling) {
		n._parent._color = colorRed
		sibling._color = colorBlock
		if n == n._parent._left {
			t.leftRotate(n._parent)
		} else {
			t.rightRotate(n._parent)
		}
	}
	t.eraseFixup3(n)
}

func (t *Tree) eraseFixup3(n *Node) {
	sibling := n.sibling()
	if !t.isRed(n._parent) &&
		!t.isRed(sibling) &&
		!t.isRed(sibling._left) &&
		!t.isRed(sibling._right) {
		sibling._color = colorRed
		t.eraseFixup1(n._parent)
	} else {
		t.eraseFixup4(n)
	}
}

func (t *Tree) eraseFixup4(n *Node) {
	sibling := n.sibling()
	if t.isRed(n._parent) &&
		!t.isRed(sibling) &&
		!t.isRed(sibling._left) &&
		!t.isRed(sibling._right) {
		sibling._color = colorRed
		n._parent._color = colorBlock
	} else {
		t.eraseFixup5(n)
	}
}

func (t *Tree) eraseFixup5(n *Node) {
	sibling := n.sibling()
	if n == n._parent._left &&
		!t.isRed(sibling) &&
		t.isRed(sibling._left) &&
		!t.isRed(sibling._right) {
		sibling._color = colorRed
		sibling._left._color = colorBlock
		t.rightRotate(sibling)
	} else if n == n._parent._right &&
		!t.isRed(sibling) &&
		t.isRed(sibling._right) &&
		!t.isRed(sibling._left) {
		sibling._color = colorRed
		sibling._right._color = colorBlock
		t.leftRotate(sibling)
	}
	t.eraseFixup6(n)
}

func (t *Tree) eraseFixup6(n *Node) {
	sibling := n.sibling()
	sibling._color = t.isRed(n._parent)
	n._parent._color = colorBlock
	if n == n._parent._left && t.isRed(sibling._right) {
		sibling._right._color = colorBlock
		t.leftRotate(n._parent)
	} else if t.isRed(sibling._left) {
		sibling._left._color = colorBlock
		t.rightRotate(n._parent)
	}
}
