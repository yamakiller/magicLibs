package mmath

//NewVector3 doc
//@Method NewVector3 @Summary create vector3 object
//@Param  (float64) x
//@Param  (float64) y
//@Param  (float64) z
//@Return (*Vector3)
func NewVector3(x, y, z float64) *Vector3 {
	return &Vector3{_x: x, _y: y, _z: z}
}

//Vector3 doc
//@Struct Vector3 doc
//@Member (float64) x
//@Member (float64) y
//@Member (float64) z
type Vector3 struct {
	_x float64
	_y float64
	_z float64
}

//Initial doc
//@Method Initial @Summary initialization vector3
//@Param (float64) x
//@Param (float64) y
//@Param (float64) z
func (slf *Vector3) Initial(x, y, z float64) {
	slf._x = x
	slf._y = y
	slf._z = z
}

//GetX doc
//@Method GetX @Summary return x
//@Return (float64) x
func (slf *Vector3) GetX() float64 {
	return slf._x
}

//GetY doc
//@Method GetY @Summary return y
//@Return (float64) y
func (slf *Vector3) GetY() float64 {
	return slf._y
}

//GetZ doc
//@Method GetZ @Summary return z
//@Return (float64)
func (slf *Vector3) GetZ() float64 {
	return slf._z
}

//SetX doc
//@Method SetX @Summary Setting x
//@Param (float64) x
func (slf *Vector3) SetX(x float64) {
	slf._x = x
}

//SetY doc
//@Method SetY @Summary Setting y
//@Param (float64) y
func (slf *Vector3) SetY(y float64) {
	slf._y = y
}

//SetZ doc
//@Method SetZ @Summary Setting z
//@Param (float64) z
func (slf *Vector3) SetZ(z float64) {
	slf._z = z
}
