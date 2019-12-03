package files

//FileHandle desc
//@Struct FileHandle desc: file handle
//@Member ([]byte) file data
//@Member (int) file size
type FileHandle struct {
	_data []byte
	_size int
}

//GetBytes desc
//@Method GetBytes desc: Returns data
//@Return ([]byte)
func (slf *FileHandle) GetBytes() []byte {
	return slf._data
}

//GetSize desc
//@Method GetSize desc: Returns data size
//@Return (int) size
func (slf *FileHandle) GetSize() int {
	return slf._size
}
