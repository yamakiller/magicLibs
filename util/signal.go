package util

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

//SignalWatch desc
//@struct SignalWatch desc: signal watch proccesser
type SignalWatch struct {
	_c chan os.Signal
	_e sync.WaitGroup
	_f func()
}

//WithCall desc
//@method WithCall desc: with signal response back call function
//@param (func()) call function
func (slf *SignalWatch) WithCall(f func()) {
	slf._f = f
}

//Initial desc
//@method Initial desc: Initialization signal watcher
//@param (func()) Signal response back call function
func (slf *SignalWatch) Initial(f func()) {
	slf._f = f
	slf._c = make(chan os.Signal)
	slf._e = sync.WaitGroup{}

	signal.Notify(slf._c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
}

//Watch desc
//@method Watch desc: start watch signal
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
//@method Wait desc: wait signal watcher exit
func (slf *SignalWatch) Wait() {
	slf._e.Wait()
}
