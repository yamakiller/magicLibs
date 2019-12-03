package containers

//Container doc
//@Struct Container @Summary is base interface that all data structures
type Container interface {
	IsEmpty() bool
	Size() int
	Clear()
	Values() []interface{}
}
