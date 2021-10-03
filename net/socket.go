package net

const (
	//INVALIDSOCKET invalid socket value
	INVALIDSOCKET = int32(-1)
)

//InvalidSocket doc
//@Method InvalidSocket @Summary Whether the socket is valid
//@Param (int32) socket
//@Return (bool) valid: false, invalid : true
func InvalidSocket(sock int32) bool {
	if sock == 0 || sock == -1 {
		return true
	}
	return false
}
