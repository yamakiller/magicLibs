package mmath

//NewRect2 new rect 2d
func NewRect2(leftUp, rightUp, leftDown, rightDown Vector2) *Rect2 {
	return &Rect2{leftUp, rightUp, leftDown, rightDown}
}

//Rect2 2d Rect
type Rect2 struct {
	_pointLeftUp    Vector2
	_pointRightUp   Vector2
	_pointLeftDown  Vector2
	_pointRightDown Vector2
}

//SetLeftUp Set Left Up point
func (slf *Rect2) SetLeftUp(point Vector2) {
	slf._pointLeftUp = point
}

//GetLeftUp Returns Left up point
func (slf *Rect2) GetLeftUp() *Vector2 {
	return &slf._pointLeftUp
}

//SetRightUp Set Right up point
func (slf *Rect2) SetRightUp(point Vector2) {
	slf._pointRightUp = point
}

//GetRightUp Returns Right up point
func (slf *Rect2) GetRightUp() *Vector2 {
	return &slf._pointRightUp
}

//SetLeftDown Set Left down point
func (slf *Rect2) SetLeftDown(point Vector2) {
	slf._pointLeftDown = point
}

//GetLeftDown Returns Left down point
func (slf *Rect2) GetLeftDown() *Vector2 {
	return &slf._pointLeftDown
}

//SetRightDown Set Right down point
func (slf *Rect2) SetRightDown(point Vector2) {
	slf._pointRightDown = point
}

//GetRightDown Returns Right down point
func (slf *Rect2) GetRightDown() *Vector2 {
	return &slf._pointRightDown
}

func (slf *Rect2) IsInvolve(point Vector2) bool {
	if point._x >= slf._pointLeftUp._x &&
		point._x < slf._pointRightUp._x &&
		point._y >= slf._pointLeftUp._y &&
		point._y < slf._pointRightDown._y {
		return true
	}
	return false
}

//Width Returns Rect width
func (slf *Rect2) Width() FValue {
	return Dist2(slf._pointLeftUp, slf._pointRightUp)
}

//Height Returns Rect height
func (slf *Rect2) Height() FValue {
	return Dist2(slf._pointLeftUp, slf._pointLeftDown)
}
