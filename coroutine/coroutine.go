package coroutine

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type coState int32
type coFunc func([]interface{})

const (
	coIdle    = int32(0)
	coRun     = int32(1)
	coClosing = int32(2)
	coDeath   = int32(3)
)

var (
	// ErrPoolStoped Coroutine pool is closed
	ErrPoolStoped = errors.New("coroutine pool stoped")
)

type task struct {
	cb     coFunc
	params []interface{}
}

func (slf *coObject) run(cop *CoroutinePool) {
	go func() {
		defer cop.wait.Done()
		for {
			for {
				select {
				case <-slf.q:
					slf.state = coDeath
					return
				case t := <-cop.taskQueue:
					atomic.CompareAndSwapInt32(&slf.state, coIdle, coRun)
					t.cb(t.params)
					t.params = nil
					atomic.CompareAndSwapInt32(&slf.state, coRun, coIdle)
				}
			}
		}
	}()
}

type coObject struct {
	state int32
	last  time.Duration
	q     chan int
	id    int
}

var (
	oneCoroutinePool     sync.Once
	defaultCoroutinePool *CoroutinePool
)

//Instance desc
//@method Instance desc: coroutine pool instance
//@return (*CoroutinePool)
func Instance() *CoroutinePool {
	oneCoroutinePool.Do(func() {
		defaultCoroutinePool = &CoroutinePool{}
	})
	return defaultCoroutinePool
}

//CoroutinePool desc
//@struct CoroutinePool desc: Coroutine pool
type CoroutinePool struct {
	taskLimit int
	maxNum    int
	minNum    int

	taskQueue chan task
	runing    int32
	curr      int32
	seq       int32
	cos       []coObject
	wait      sync.WaitGroup
	quit      chan int
}

func (slf *CoroutinePool) scheduer() {
	defer slf.wait.Done()
	for {
		select {
		case <-slf.quit:
			goto scheduer_end
		default:
		}

		now := time.Now().Unix()
		for _, v := range slf.cos {
			if v.state == coRun ||
				v.state == coDeath ||
				v.state == coClosing {
				continue
			}

			if v.state == coIdle && ((time.Duration(now) - v.last) > time.Second*60) {
				if atomic.CompareAndSwapInt32(&v.state, int32(coIdle), int32(coClosing)) {
					close(v.q)
				}
			}
		}

		time.Sleep(time.Millisecond * 1000)
	}
scheduer_end:
	for i := 0; i < slf.maxNum; i++ {
		if slf.cos[i].state != coDeath && slf.cos[i].state != coClosing {
			atomic.StoreInt32(&slf.cos[i].state, coClosing)
			close(slf.cos[i].q)
		}
	}
}

//Start desc
//@method Start desc: Start the coroutine pool
//@param (int) coroutine max of number
//@param (int) coroutine min of number
//@param (int) coroutine task max limit
func (slf *CoroutinePool) Start(max int, min int, taskmax int) {
	slf.maxNum = max
	slf.minNum = min
	slf.taskLimit = taskmax

	if slf.maxNum == 0 {
		slf.maxNum = 65535
	}
	atomic.StoreInt32(&slf.seq, 1)

	now := time.Now().Unix()
	slf.taskQueue = make(chan task, slf.taskLimit)
	slf.cos = make([]coObject, slf.maxNum)
	for k, v := range slf.cos {
		v.id = k + 1
		v.state = coDeath
		v.last = time.Duration(now)
	}

	if slf.maxNum > 0 && slf.maxNum > slf.minNum {
		slf.wait.Add(1)
		go slf.scheduer()
	}

	for i := 0; i < slf.minNum; i++ {
		slf.cos[i].state = coIdle
		slf.startOne(i)
	}
}

func (slf *CoroutinePool) startOne(idx int) {
	slf.wait.Add(1)
	atomic.AddInt32(&slf.seq, 1)
	atomic.AddInt32(&slf.curr, 1)

	slf.cos[idx].run(slf)
}

//Stop desc
//@method Stop desc: Stop the coroutine pool
func (slf *CoroutinePool) Stop() {
	close(slf.quit)
	slf.wait.Wait()
}

//Go desc
//@method Go desc: Running coroutine
//@param  (func(params []interface{})) call function
//@param  (...interface{}) call args
//@return (error)
func (slf *CoroutinePool) Go(f func(params []interface{}), params ...interface{}) error {

	select {
	case <-slf.quit:
		return ErrPoolStoped
	default:
	}
	runing := atomic.LoadInt32(&slf.runing)
	curr := atomic.LoadInt32(&slf.curr)

	if runing >= curr && curr < int32(slf.maxNum) {
		for i := 0; i < slf.maxNum; i++ {
			hash := ((int32(i) + slf.seq) % int32(slf.maxNum))
			if atomic.CompareAndSwapInt32(&slf.cos[hash].state, coDeath, coIdle) {
				slf.seq = int32(hash) + slf.seq
				slf.startOne(int(hash))
				break
			}
		}
	}

	select {
	case <-slf.quit:
		return ErrPoolStoped
	case slf.taskQueue <- task{f, params}:
		return nil
	}
}