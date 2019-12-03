package containers

//Container desc
//@Struct Container desc: is base interface that all data structures
type Container interface {
	IsEmpty() bool
	Size() int
	Clear()
	Values() []interface{}
}
