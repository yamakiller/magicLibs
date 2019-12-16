package mmath

//NewVector2 doc
//@Method NewVector2 @Summary create vector2 object
//@Param  (float64) x
//@Param  (float64) y
//@Return (*Vector2)
func NewVector2(x, y float64) *Vector2 {
	return &Vector2{_x: x, _y: y}
}

//Vector2 doc
//@Struct Vector2 doc
//@Member (float64) x
//@Member (float64) y
type Vector2 struct {
	_x float64
	_y float64
}

//Initial doc
//@Method Initial @Summary initialization vector2
//@Param (float64) x
//@Param (float64) y
func (slf *Vector2) Initial(x, y float64) {
	slf._x = x
	slf._y = y
}

//GetX doc
//@Method GetX @Summary return x
//@Return (float64) x
func (slf *Vector2) GetX() float64 {
	return slf._x
}

//GetY doc
//@Method GetY @Summary return y
//@Return (float64) y
func (slf *Vector2) GetY() float64 {
	return slf._y
}

//SetX doc
//@Method SetX @Summary Setting x
//@Param (float64) x
func (slf *Vector2) SetX(x float64) {
	slf._x = x
}

//SetY doc
//@Method SetY @Summary Setting y
//@Param (float64) y
func (slf *Vector2) SetY(y float64) {
	slf._y = y
}
