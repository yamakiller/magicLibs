package math

//NewVector2 desc
//@Method NewVector2 desc: create vector2 object
//@Param  (float64) x
//@Param  (float64) y
//@Return (*Vector2)
func NewVector2(x, y float64) *Vector2 {
	return &Vector2{_x: x, _y: y}
}

//Vector2 desc
//@Struct Vector2 desc
//@Member (float64) x
//@Member (float64) y
type Vector2 struct {
	_x float64
	_y float64
}

//Initial desc
//@Method Initial desc: initialization vector2
//@Param (float64) x
//@Param (float64) y
func (slf *Vector2) Initial(x, y float64) {
	slf._x = x
	slf._y = y
}

//GetX desc
//@Method GetX desc: return x
//@Return (float64) x
func (slf *Vector2) GetX() float64 {
	return slf._x
}

//GetY desc
//@Method GetY desc: return y
//@Return (float64) y
func (slf *Vector2) GetY() float64 {
	return slf._y
}

//SetX desc
//@Method SetX desc: Setting x
//@Param (float64) x
func (slf *Vector2) SetX(x float64) {
	slf._x = x
}

//SetY desc
//@Method SetY desc: Setting y
//@Param (float64) y
func (slf *Vector2) SetY(y float64) {
	slf._y = y
}
