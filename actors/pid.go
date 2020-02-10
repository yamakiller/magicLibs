package actors

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

//PID 外部接口
type PID struct {
	ID      uint32
	_h      handle
	_parent *Core
}

//ToString 返回ID字符串
func (slf *PID) ToString() string {
	return fmt.Sprintf(".%08x", slf.ID)
}

func (slf *PID) ref() handle {
	p := (*handle)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&slf._h))))
	if p != nil {
		if l, ok := (*p).(*actorHandle); ok && atomic.LoadInt32(&l._death) == 1 {
			atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&slf._h)), nil)
		} else {
			return *p
		}
	}

	ref := slf._parent.getHandle(slf)
	if ref != nil {
		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&slf._h)), unsafe.Pointer(&ref))
	}

	return ref
}

//Post 发送消息
func (slf *PID) Post(message interface{}) {
	slf.postUsrMessage(message)
}

func (slf *PID) postUsrMessage(message interface{}) {
	ref := slf.ref()
	ref.postUsrMessage(slf, message)
	/*overload := ref.OverloadUsrMessage()
	if overload > 0 {
		logger.Warning(pid.ID, "mailbox overload :%d", overload)
	}*/
}

func (slf *PID) postSysMessage(message interface{}) {
	slf.ref().postSysMessage(slf, message)
}

//Stop 停止
func (slf *PID) Stop() {
	slf.ref().Stop(slf)
}
