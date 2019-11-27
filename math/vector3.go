package math

//NewVector3 desc
//@method NewVector3 desc: create vector3 object
//@param  (float64) x
//@param  (float64) y
//@param  (float64) z
//@return (*Vector3)
func NewVector3(x, y, z float64) *Vector3 {
	return &Vector3{_x: x, _y: y, _z: z}
}

//Vector3 desc
//@struct Vector3 desc
//@member (float64) x
//@member (float64) y
//@member (float64) z
type Vector3 struct {
	_x float64
	_y float64
	_z float64
}

//Initial desc
//@method Initial desc: initialization vector3
//@param (float64) x
//@param (float64) y
//@param (float64) z
func (slf *Vector3) Initial(x, y, z float64) {
	slf._x = x
	slf._y = y
	slf._z = z
}

//GetX desc
//@method GetX desc: return x
//@return (float64) x
func (slf *Vector3) GetX() float64 {
	return slf._x
}

//GetY desc
//@method GetY desc: return y
//@return (float64) y
func (slf *Vector3) GetY() float64 {
	return slf._y
}

//GetZ desc
//@method GetZ desc: return z
//@return (float64)
func (slf *Vector3) GetZ() float64 {
	return slf._z
}

//SetX desc
//@method SetX desc: Setting x
//@param (float64) x
func (slf *Vector3) SetX(x float64) {
	slf._x = x
}

//SetY desc
//@method SetY desc: Setting y
//@param (float64) y
func (slf *Vector3) SetY(y float64) {
	slf._y = y
}

//SetZ desc
//@method SetZ desc: Setting z
//@param (float64) z
func (slf *Vector3) SetZ(z float64) {
	slf._z = z
}
