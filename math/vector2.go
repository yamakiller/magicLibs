package math

//NewVector2 desc
//@method NewVector2 desc: create vector2 object
//@param  (float64) x
//@param  (float64) y
//@return (*Vector2)
func NewVector2(x, y float64) *Vector2 {
	return &Vector2{_x: x, _y: y}
}

//Vector2 desc
//@struct Vector2 desc
//@member (float64) x
//@member (float64) y
type Vector2 struct {
	_x float64
	_y float64
}

//Initial desc
//@method Initial desc: initialization vector2
//@param (float64) x
//@param (float64) y
func (slf *Vector2) Initial(x, y float64) {
	slf._x = x
	slf._y = y
}

//GetX desc
//@method GetX desc: return x
//@return (float64) x
func (slf *Vector2) GetX() float64 {
	return slf._x
}

//GetY desc
//@method GetY desc: return y
//@return (float64) y
func (slf *Vector2) GetY() float64 {
	return slf._y
}

//SetX desc
//@method SetX desc: Setting x
//@param (float64) x
func (slf *Vector2) SetX(x float64) {
	slf._x = x
}

//SetY desc
//@method SetY desc: Setting y
//@param (float64) y
func (slf *Vector2) SetY(y float64) {
	slf._y = y
}
