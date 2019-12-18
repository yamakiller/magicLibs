package aoi

/*
import (
	"errors"

	"github.com/yamakiller/magicLibs/mmath"
)

const (
	nodeLeftUp    = 0
	nodeRightUp   = 1
	nodeLeftDown  = 2
	nodeRightDown = 3
	nodeRoot      = 4
)

type quadNode struct {
	_rect     mmath.Rect2
	_quadrant int
	_shape    interface{}
	_parent   *quadNode
	_child    []*quadNode
}

func newQuadNode() *quadNode {
	return &quadNode{_child: make([]*quadNode, 4)}
}

//QuadTree Quadtree
type QuadTree struct {
	_root     *quadNode
	_depthCap int
	_alloc    func() interface{}
}

//Initial initial build quad tree
func (slf *QuadTree) Initial(rect mmath.Rect2,
	depth int,
	alloc func() interface{}) error {
	if depth > 6 {
		return errors.New("Too deep")
	}
	slf._depthCap = depth
	slf._root = newQuadNode()
	slf._root._rect = rect
	slf._root._quadrant = nodeRoot
	slf._alloc = alloc

	currDepth := 0
	slf.spawn(slf._root, currDepth)

	return nil
}

func (slf *QuadTree) spawn(parent *quadNode, depth int) {
	if parent._child[nodeLeftUp] == nil && depth+1 < slf._depthCap {
		parent._child[nodeLeftUp] = newQuadNode()
		parent._child[nodeLeftUp]._parent = parent
		parent._child[nodeLeftUp]._quadrant = nodeLeftUp
		leftUp := mmath.NewVector2(parent._rect.GetLeftUp().GetX(), parent._rect.GetLeftUp().GetY())
		rigthUp := mmath.NewVector2(parent._rect.GetRightUp().GetX()*0.5, parent._rect.GetRightUp().GetY())
		leftDown := mmath.NewVector2(parent._rect.GetLeftDown().GetX(), parent._rect.GetLeftDown().GetY()*0.5)
		rigthDown := mmath.NewVector2(parent._rect.GetRightDown().GetX()*0.5, parent._rect.GetRightDown().GetY()*0.5)
		parent._child[nodeLeftUp]._rect = *mmath.NewRect2(*leftUp, *rigthUp, *leftDown, *rigthDown)
		parent._child[nodeLeftUp]._shape = slf._alloc()
		slf.spawn(parent._child[nodeLeftUp], depth+1)
	}

	if parent._child[nodeRightUp] == nil && depth+1 < slf._depthCap {
		parent._child[nodeRightUp] = newQuadNode()
		parent._child[nodeRightUp]._parent = parent
		parent._child[nodeRightUp]._quadrant = nodeRightUp
		leftUp := mmath.NewVector2(parent._rect.GetRightUp().GetX()*0.5, parent._rect.GetLeftUp().GetY())
		rigthUp := mmath.NewVector2(parent._rect.GetRightUp().GetX(), parent._rect.GetLeftUp().GetY())
		leftDown := mmath.NewVector2(parent._rect.GetRightUp().GetX()*0.5, parent._rect.GetLeftDown().GetY()*0.5)
		rigthDown := mmath.NewVector2(parent._rect.GetRightDown().GetX(), parent._rect.GetRightDown().GetY()*0.5)
		parent._child[nodeRightUp]._rect = *mmath.NewRect2(*leftUp, *rigthUp, *leftDown, *rigthDown)
		parent._child[nodeRightUp]._shape = slf._alloc()
		slf.spawn(parent._child[nodeRightUp], depth+1)
	}

	if parent._child[nodeLeftDown] == nil && depth+1 < slf._depthCap {
		parent._child[nodeLeftDown] = newQuadNode()
		parent._child[nodeLeftDown]._parent = parent
		parent._child[nodeLeftDown]._quadrant = nodeLeftDown
		leftUp := mmath.NewVector2(parent._rect.GetRightUp().GetX(), parent._rect.GetLeftDown().GetY()*0.5)
		rigthUp := mmath.NewVector2(parent._rect.GetRightDown().GetX()*0.5, parent._rect.GetLeftDown().GetY()*0.5)
		leftDown := mmath.NewVector2(parent._rect.GetLeftDown().GetX(), parent._rect.GetLeftDown().GetY())
		rigthDown := mmath.NewVector2(parent._rect.GetRightDown().GetX()*0.5, parent._rect.GetRightDown().GetY())
		parent._child[nodeLeftDown]._rect = *mmath.NewRect2(*leftUp, *rigthUp, *leftDown, *rigthDown)
		parent._child[nodeLeftDown]._shape = slf._alloc()
		slf.spawn(parent._child[nodeLeftDown], depth+1)
	}

	if parent._child[nodeRightDown] == nil && depth+1 < slf._depthCap {
		parent._child[nodeRightDown] = newQuadNode()
		parent._child[nodeRightDown]._parent = parent
		parent._child[nodeRightDown]._quadrant = nodeRightDown
		leftUp := mmath.NewVector2(parent._rect.GetRightDown().GetX()*0.5, parent._rect.GetRightDown().GetY()*0.5)
		rigthUp := mmath.NewVector2(parent._rect.GetRightDown().GetX(), parent._rect.GetLeftDown().GetY()*0.5)
		leftDown := mmath.NewVector2(parent._rect.GetRightDown().GetX()*0.5, parent._rect.GetRightDown().GetY())
		rigthDown := mmath.NewVector2(parent._rect.GetRightDown().GetX(), parent._rect.GetRightDown().GetY())
		parent._child[nodeRightDown]._rect = *mmath.NewRect2(*leftUp, *rigthUp, *leftDown, *rigthDown)
		parent._child[nodeRightDown]._shape = slf._alloc()
		slf.spawn(parent._child[nodeRightDown], depth+1)
	}
}

//Match match all node
func (slf *QuadTree) Match(f func(interface{})) {

}

//QueryShape Returs space informat
func (slf *QuadTree) QueryShape(point mmath.Vector2) interface{} {
	return slf.query(point, slf._root)
}

func (slf *QuadTree) query(point mmath.Vector2, cur *quadNode) interface{} {
	if !cur._rect.IsInvolve(point) {
		return nil
	}

	if cur._child[nodeLeftUp] == nil &&
		cur._child[nodeLeftDown] == nil &&
		cur._child[nodeRightUp] == nil &&
		cur._child[nodeRightDown] == nil {
		return cur._shape
	}

	if r := slf.query(point, cur._child[nodeLeftUp]); r != nil {
		return r
	}

	if r := slf.query(point, cur._child[nodeLeftDown]); r != nil {
		return r
	}

	if r := slf.query(point, cur._child[nodeRightUp]); r != nil {
		return r
	}

	if r := slf.query(point, cur._child[nodeRightDown]); r != nil {
		return r
	}

	return nil
}
*/
