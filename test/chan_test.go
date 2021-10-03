package test

import (
	"errors"
	"sync"
	"testing"
	"time"
)

type control struct {
	_closed chan bool
	_queue  chan interface{}
	_wwg    sync.WaitGroup
	_rwg    sync.WaitGroup
}

func (slf *control) Start() {
	slf._wwg.Add(1)
	slf._rwg.Add(1)
	go slf.writeServe()
	go slf.readServe()
}

func (slf *control) writeServe() {
	defer func() {
		slf._rwg.Wait()
		close(slf._queue)
		slf._wwg.Done()
	}()

	for {
		select {
		case <-slf._closed:
			goto exit
		case msg := <-slf._queue:
			if msg != nil {

			}
		}
	}
exit:
}

func (slf *control) readServe() {
	defer func() {
		slf._rwg.Done()
	}()

	for {
		select {
		case <-slf._closed:
			goto exit
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
exit:
}

func (slf *control) SendTo(msg interface{}) error {
	select {
	case <-slf._closed:
		return errors.New("closed")
	default:
	}

	slf._queue <- msg
	return nil
}

func (slf *control) Close() error {
	close(slf._closed)
	slf._rwg.Wait()
	slf._wwg.Wait()
	return nil
}

func tSend(c *control, w *sync.WaitGroup) {
	defer w.Done()

	for {
		if err := c.SendTo(1); err != nil {
			break
		}
		time.Sleep(1 * time.Millisecond)
	}
}

func TestSim(t *testing.T) {
	var gw sync.WaitGroup

	total := 0
	for {

		total++
		if total > 100 {
			break
		}

		c := &control{
			_closed: make(chan bool, 1),
			_queue:  make(chan interface{}, 64),
		}

		c.Start()
		gw.Add(1)
		go tSend(c, &gw)

		iw := 0
		for {
			iw++
			if iw > 100 {
				break
			}
			time.Sleep(1 * time.Millisecond)
		}

		c.Close()
		gw.Wait()
	}

	await := make(chan bool, 1)

	select {
	case <-await:
	default:
		close(await)
	}

	select {
	case <-await:
	default:
		close(await)
	}
}
