package aoi

import (
	"github.com/yamakiller/magicLibs/mmath"
)

//CrossNode cross linked list node
type CrossNode struct {
	_xPrev, _xNext *CrossNode
	_yPrev, _yNext *CrossNode
	_zPrev, _zNext *CrossNode
	_key           interface{}
	_pos           mmath.Vector3
}

//GetKey Returns key
func (slf *CrossNode) GetKey() interface{} {
	return slf._key
}

//CrossLinked cross linked list
type CrossLinked struct {
	Isz          bool
	_head, _tail *CrossNode
	_sz          int
}

func (slf *CrossLinked) add(node *CrossNode) {
	//x
	cur := slf._head._xNext
	for cur != nil {
		if (cur._pos.GetX() > node._pos.GetX()) || cur == slf._tail {
			node._xNext = cur
			node._xPrev = cur._xPrev
			cur._xPrev._xNext = node
			cur._xPrev = node
			break
		}

		cur = cur._xNext
	}

	//y
	cur = slf._head._yNext
	for cur != nil {
		if cur._pos.GetY() > node._pos.GetY() || cur == slf._tail {
			node._yNext = cur
			node._yPrev = cur._yPrev
			cur._yPrev._yNext = node
			cur._yPrev = node
		}
		cur = cur._yNext
	}

	//z
	if !slf.Isz {
		cur = slf._head._zNext
		for cur != nil {
			if cur._pos.GetZ() > node._pos.GetZ() || cur == slf._tail {
				node._zNext = cur
				node._zPrev = cur._yPrev
				cur._zPrev._zNext = node
				cur._zPrev = node
			}
			cur = cur._zNext
		}
	}
}

//Enter enter scenes
func (slf *CrossLinked) Enter(key interface{}, pos mmath.Vector3) *CrossNode {
	newNode := &CrossNode{_key: key, _pos: pos}
	slf.add(newNode)
	return newNode
}

//Leave leave scenes
func (slf *CrossLinked) Leave(node *CrossNode) {
	node._xPrev._xNext = node._xNext
	node._xNext._xPrev = node._xPrev
	node._yPrev._yPrev = node._yNext
	node._yNext._yNext = node._yPrev
	if !slf.Isz {
		node._zPrev._zNext = node._zNext
		node._zNext._zPrev = node._zPrev
	}

	node._xPrev = nil
	node._xNext = nil
	node._yPrev = nil
	node._yNext = nil
	if !slf.Isz {
		node._zPrev = nil
		node._zNext = nil
	}
}

//Move move scenes
func (slf *CrossLinked) Move(node *CrossNode, pos mmath.Vector3) {
	//x
	if node._pos.GetX() != pos.GetX() {
		if pos.GetX() > node._pos.GetX() {

			cur := node._xNext
			for cur != slf._tail {
				if pos.GetX() <= cur._pos.GetX() {
					break
				}

				cur._xPrev = node._xPrev
				node._xPrev._xNext = cur
				node._xNext = cur._xNext
				cur._xNext = node
				//next a
				cur = node._xNext
			}
		} else {
			cur := node._xPrev
			for cur != slf._head {
				if pos.GetX() >= cur._pos.GetX() {
					break
				}

				cur._xPrev._xNext = node
				node._xPrev = cur._xPrev
				cur._xPrev = node
				cur._xNext = node._xNext
				node._xNext = cur
				//prev a
				cur = node._xPrev
			}
		}

		node._pos.SetX(pos.GetX())
	}

	//y
	if node._pos.GetY() != pos.GetY() {
		if pos.GetY() > node._pos.GetY() {

			cur := node._yNext
			for cur != slf._tail {
				if pos.GetY() <= cur._pos.GetY() {
					break
				}

				cur._yPrev = node._yPrev
				node._yPrev._yNext = cur
				node._yNext = cur._yNext
				cur._yNext = node
				//next a
				cur = node._yNext
			}
		} else {
			cur := node._yPrev
			for cur != slf._head {
				if pos.GetY() >= cur._pos.GetY() {
					break
				}

				cur._yPrev._yNext = node
				node._yPrev = cur._yPrev
				cur._yPrev = node
				cur._yNext = node._yNext
				node._yNext = cur
				//prev a
				cur = node._yPrev
			}
		}

		node._pos.SetY(pos.GetY())
	}

	//z
	if !slf.Isz {
		if node._pos.GetZ() != pos.GetZ() {
			if pos.GetZ() > node._pos.GetZ() {

				cur := node._zNext
				for cur != slf._tail {
					if pos.GetZ() <= cur._pos.GetZ() {
						break
					}

					cur._zPrev = node._zPrev
					node._zPrev._zNext = cur
					node._zNext = cur._zNext
					cur._zNext = node
					//next a
					cur = node._zNext
				}
			} else {
				cur := node._zPrev
				for cur != slf._head {
					if pos.GetZ() >= cur._pos.GetZ() {
						break
					}

					cur._zPrev._zNext = node
					node._zPrev = cur._zPrev
					cur._zPrev = node
					cur._zNext = node._zNext
					node._zNext = cur
					//prev a
					cur = node._zPrev
				}
			}

			node._pos.SetZ(pos.GetZ())
		}
	}
}
