package boxs

import (
	"reflect"
	"sync"

	"github.com/yamakiller/magicLibs/actors"
	"github.com/yamakiller/magicLibs/actors/messages"
)

// SpawnBox create an box
func SpawnBox(pid *actors.PID) *Box {
	if pid != nil {
		return &Box{
			_pid:           pid,
			_events:        make(map[interface{}][]Method),
			_started:       make(chan bool, 1),
			_stopped:       make(chan bool, 1),
			_stoppingAfter: nil,
			_stoppedAfter:  nil,
		}
	}

	return &Box{
		_events:        make(map[interface{}][]Method),
		_started:       make(chan bool, 1),
		_stopped:       make(chan bool, 1),
		_stoppingAfter: nil,
		_stoppedAfter:  nil,
	}
}

// Box container for executing logic
type Box struct {
	_pid           *actors.PID
	_events        map[interface{}][]Method
	_evmutx        sync.Mutex
	_context       Context
	_started       chan bool
	_stopped       chan bool
	_stoppingAfter func(context *Context)
	_stoppedAfter  func(context *Context)
}

// GetPID Returns pid
func (slf *Box) GetPID() *actors.PID {
	return slf._pid
}

// WithPID setting pid
func (slf *Box) WithPID(pid *actors.PID) {
	slf._pid = pid
}

// WithStoppingAfterCallback You can set the handler during the execution stop process
func (slf *Box) WithStoppingAfterCallback(cb func(context *Context)) {
	slf._stoppingAfter = cb
}

// WithStoppedAfterCallback setting stopped after callback If this callback function does not exist,
// it will be automatically destroyed. Otherwise, you need to call Destory at the corresponding location
// to manually destroy it.
func (slf *Box) WithStoppedAfterCallback(cb func(context *Context)) {
	slf._stoppedAfter = cb
}

// StartedWait wait box started
func (slf *Box) StartedWait() {
	select {
	case <-slf._started:
		break
	}
}

// Shutdown shutdown box
func (slf *Box) Shutdown() {
	slf._pid.Stop()
}

// WaitStopped wait stopped
func (slf *Box) StoppedWait() {
	select {
	case <-slf._stopped:
	}
}

// Destory resources to be released
func (slf *Box) Destory() {
	if slf._started != nil {
		close(slf._started)
		slf._started = nil
	}

	if slf._stopped != nil {
		close(slf._stopped)
		slf._stopped = nil
	}

	slf._events = nil
	slf._context.Context = nil
	slf._context._funs = nil
	slf._stoppingAfter = nil
	slf._stoppedAfter = nil
}

// Register register event
func (slf *Box) Register(key interface{}, args ...Method) {
	var ms []Method
	ms = append(ms, args...)
	slf._evmutx.Lock()
	defer slf._evmutx.Unlock()
	slf._events[key] = ms
}

// Receive event receive proccess
func (slf *Box) Receive(context *actors.Context) {
	slf._context.Context = context
	message := context.Message()
	switch msg := message.(type) {
	case *actors.Pack:
		message = msg.Message
	default:
	}

	var after Method
	switch message.(type) {
	case *messages.Started:
		after = slf.onStartedAfter
	case *messages.Stopping:
		after = slf.onStoppingAfter
	case *messages.Stopped:
		after = slf.onStoppedAfter
	default:
	}

	slf._evmutx.Lock()
	if f, ok := slf._events[reflect.TypeOf(message)]; ok && len(f) > 0 {
		slf._evmutx.Unlock()
		slf._context._funs = f
		slf._context._idx = 0
		slf._context._funs[0](&slf._context)
		goto end
	}
	slf._evmutx.Unlock()

	if after != nil {
		goto end
	}

	slf.onError(&slf._context)
end:
	if after != nil {
		//default event before function
		after(&slf._context)
	}
}

func (slf *Box) onStartedAfter(context *Context) {
	slf._started <- true
}

func (slf *Box) onStoppingAfter(context *Context) {
	if slf._stoppingAfter != nil {
		slf._stoppingAfter(context)
	}
}

func (slf *Box) onStoppedAfter(context *Context) {
	slf._stopped <- true
	if slf._stoppedAfter != nil {
		slf._stoppedAfter(context)
	} else { // 自我销毁
		slf.Destory()
	}
}

func (slf *Box) onError(context *Context) {
	context.Error("Box %+v message is undefined", context.Message())
}
