package files

//FileHandle desc
//@struct FileHandle desc: file handle
//@member ([]byte) file data
//@member (int) file size
type FileHandle struct {
	_data []byte
	_size int
}

//GetBytes desc
//@method GetBytes desc: Returns data
//@return ([]byte)
func (slf *FileHandle) GetBytes() []byte {
	return slf._data
}

//GetSize desc
//@method GetSize desc: Returns data size
//@return (int) size
func (slf *FileHandle) GetSize() int {
	return slf._size
}
