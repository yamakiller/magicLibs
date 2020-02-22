package contractor

//Conv ID预约器接口
type Conv interface {
	//预约一个ID
	Reserve(uint64) (uint32, uint64, error)
	Confirm([]byte) bool
	Authorized(id uint32) bool
	Close(uint32)
}
