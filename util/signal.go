package util

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

//SignalWatch desc
//@Struct SignalWatch desc: signal watch proccesser
type SignalWatch struct {
	_c chan os.Signal
	_e sync.WaitGroup
	_f func()
}

//Initial desc
//@Method Initial desc: Initialization signal watcher
//@Param (func()) Signal response back call function
func (slf *SignalWatch) Initial(f func()) {
	slf._f = f
	slf._c = make(chan os.Signal)
	slf._e = sync.WaitGroup{}
	signal.Notify(slf._c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
}

//Watch desc
//@Method Watch desc: start watch signal
func (slf *SignalWatch) Watch() {
	slf._e.Add(1)
	go func() {
		defer slf._e.Done()
		for s := range slf._c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				slf._f()
				return
			default:
				break
			}
		}
	}()
}

//Wait desc
//@Method Wait desc: wait signal watcher exit
func (slf *SignalWatch) Wait() {
	slf._e.Wait()
}
