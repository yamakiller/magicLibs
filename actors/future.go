package actors

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrTimeout    = errors.New("future timeout")
	ErrDeadLetter = errors.New("future: dead letter")
)

type Future struct {
	_cond   *sync.Cond
	_done   bool
	_result interface{}
	_err    error
	_t      *time.Timer
}

func (f *Future) wait() {
	f._cond.L.Lock()
	for !f._done {
		f._cond.Wait()
	}
	f._cond.L.Unlock()
}

func (f *Future) Result() (interface{}, error) {
	f.wait()

	return f._result, f._err
}

func (f *Future) Wait() error {
	f.wait()

	return f._err
}
