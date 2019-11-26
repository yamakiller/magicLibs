package redblacktree

// Iterator holding the iterator`s state
type Iterator struct {
	_t *Tree
	_n *Node
	_p position
}

type position byte

const (
	begin, between, end position = 0, 1, 2
)

// Iterator returns a stateful iterator key/value pairs.
func (tr *Tree) Iterator() Iterator {
	return Iterator{_t: tr, _n: nil, _p: begin}
}

// Next moves the iterator to the next element
func (it *Iterator) Next() bool {
	if it._p == end {
		goto end
	}

	if it._p == begin {
		l := it._t.left()
		if l == nil {
			goto end
		}
		it._n = l
		goto between
	}
	if it._n._right != nil {
		it._n = it._n._right
		for it._n._left != nil {
			it._n = it._n._left
		}
		goto between
	}
	if it._n._parent != nil {
		n := it._n
		for it._n._parent != nil {
			it._n = it._n._parent
			if it._t._compare(n._key, it._n._key) >= 0 {
				goto between
			}
		}
	}
end:
	it._n = nil
	it._p = end
	return false
between:
	it._p = between
	return true
}

// Prev moves the iterator to the prev element
func (it *Iterator) Prev() bool {
	if it._p == begin {
		goto begin
	}
	if it._p == end {
		r := it._t.right()
		if r == nil {
			goto begin
		}
		it._n = r
		goto between
	}
	if it._n._left != nil {
		it._n = it._n._left
		for it._n._right != nil {
			it._n = it._n._right
		}
		goto between
	}
	if it._n._parent != nil {
		n := it._n
		for it._n._parent != nil {
			it._n = it._n._parent
			if it._t._compare(n._key, it._n._key) >= 0 {
				goto between
			}
		}
	}
begin:
	it._n = nil
	it._p = begin
	return false

between:
	it._p = between
	return true
}

// It returns current Node
func (it *Iterator) It() *Node {
	return it._n
}

// Value returns the current element`s value.
func (it *Iterator) Value() interface{} {
	return it._n._val
}

// Key returns the current element`s key.
func (it *Iterator) Key() interface{} {
	return it._n._key
}

// Begin resets the iterator to its start(one-before-first)
func (it *Iterator) Begin() {
	it._n = nil
	it._p = begin
}

// End resets the iterator to its end (one-past-the-end)
func (it *Iterator) End() {
	it._n = nil
	it._p = end
}

// First moves the iterator to the first element
func (it *Iterator) First() bool {
	it.Begin()
	return it.Next()
}

// Last moves the iterator to the last element
func (it *Iterator) Last() bool {
	it.End()
	return it.Prev()
}
