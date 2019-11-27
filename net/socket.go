package net

const (
	//INVALIDSOCKET invalid socket value
	INVALIDSOCKET = int32(-1)
)

//InvalidSocket desc
//@method InvalidSocket desc: Whether the socket is valid
//@param (int32) socket
//@return (bool) valid: false, invalid : true
func InvalidSocket(sock int32) bool {
	if sock == 0 || sock == -1 {
		return true
	}
	return false
}
