package middle

type Middleware interface {
	Error(error)
	Debug(error)
}
