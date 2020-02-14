package boxs

import (
	"reflect"

	"github.com/yamakiller/magicLibs/actors"
	"github.com/yamakiller/magicLibs/actors/messages"
)

//SpawnBox create an box
func SpawnBox(pid actors.PID) *Box {
	return &Box{
		_pid:     pid,
		_events:  make(map[interface{}]interface{}),
		_started: make(chan bool),
		_stopped: make(chan bool),
	}
}

//Box container for executing logic
type Box struct {
	_pid     actors.PID
	_events  map[interface{}]interface{}
	_started chan bool
	_stopped chan bool
}

//GetPID Returns pid
func (slf *Box) GetPID() *actors.PID {
	return &slf._pid
}

//WithPID setting pid
func (slf *Box) WithPID(pid actors.PID) {
	slf._pid = pid
}

//StartedWait wait box started
func (slf *Box) StartedWait() {
	select {
	case <-slf._started:
		break
	}
}

//Shutdown shutdown box
func (slf *Box) Shutdown() {
	slf._pid.Stop()
}

//ShutdownWait Close the box and wait for resources to be released
func (slf *Box) ShutdownWait() {
	slf._pid.Stop()
	select {
	case <-slf._stopped:
	}
}

//Register register event
func (slf *Box) Register(key, value interface{}) {
	slf._events[key] = value
}

//Receive event receive proccess
func (slf *Box) Receive(context *actors.Context) {
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

	if f, ok := slf._events[reflect.TypeOf(message)]; ok {
		f.(Method)(context)
		goto end
	}

	if after != nil {
		goto end
	}

	slf.onError(context)
end:
	if after != nil {
		//default event before function
		after(context)
	}
}

func (slf *Box) onStartedAfter(context *actors.Context) {
	slf._started <- true
}

func (slf *Box) onStoppingAfter(context *actors.Context) {

}

func (slf *Box) onStoppedAfter(context *actors.Context) {
	close(slf._stopped)
}

func (slf *Box) onError(context *actors.Context) {
	context.Error("Box %+v message is undefined", context, context.Message())
}
