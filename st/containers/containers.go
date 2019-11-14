package containers

//Container desc
//@struct Container desc: is base interface that all data structures
type Container interface {
	IsEmpty() bool
	Size() int
	Clear()
	Values() []interface{}
}
