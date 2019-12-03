package net

import (
	"context"
	"errors"
	"sync"
	"time"

	"google.golang.org/grpc"
)

type grpcIdleConn struct {
	_c *grpc.ClientConn
	_t time.Time
}

var (
	errGRPCSetClosed   = errors.New("pool is closed")
	errGRPCSetInvalid  = errors.New("invalid config")
	errGRPCSetRejected = errors.New("connection is nil. rejecting")
	errGRPCSetTargets  = errors.New("targets server is empty")
)

// GRPCClientSet : GRPC 客户端集
type GRPCClientSet struct {
	IdleTimeout time.Duration
	MinConn     int
	MaxConn     int
	ConnTimeout time.Duration
	_cs         chan *grpcIdleConn
	_factory    func() (*grpc.ClientConn, error)
	_close      func(*grpc.ClientConn) error
	_mx         sync.Mutex
}

//Initial desc
//@Method Initial desc:Initialization grpc client sets
//@Param (string) target address
//@Param (...grpc.DialOption) grpc client options
func (slf *GRPCClientSet) Initial(target string, opts ...grpc.DialOption) error {
	slf._cs = make(chan *grpcIdleConn, slf.MaxConn)
	slf._factory = func() (*grpc.ClientConn, error) {
		ctx, cancel := context.WithTimeout(context.Background(), slf.ConnTimeout)
		defer cancel()

		return grpc.DialContext(ctx, target)
	}
	slf._close = func(c *grpc.ClientConn) error { return c.Close() }

	for i := 0; i < slf.MinConn; i++ {
		conn, err := slf._factory()
		if err != nil {
			slf.Close()
			return nil
		}

		slf._cs <- &grpcIdleConn{_c: conn, _t: time.Now()}
	}

	return nil
}

//Close desc
//@Method Close desc: Closing grpc client set
func (slf *GRPCClientSet) Close() {
	slf._mx.Lock()
	conns := slf._cs
	slf._cs = nil
	slf._factory = nil
	closeFun := slf._close
	slf._close = nil
	slf._mx.Unlock()

	if conns == nil {
		return
	}

	close(conns) //？修改退出 未必可以全部删除
	for wrapConn := range conns {
		closeFun(wrapConn._c)
	}
}

//Invoke desc
//@Method Invoke desc: Invoke grpc method
func (slf *GRPCClientSet) Invoke(method string, args, reply interface{}) error { //优化参数设置
	conn, err := slf.getConn()
	if err != nil {
		return err
	}
	defer slf.putConn(conn)

	ctx, cancel := context.WithTimeout(context.Background(), slf.ConnTimeout)
	defer cancel()

	return conn.Invoke(ctx, method, args, reply)
}

func (slf *GRPCClientSet) getConn() (*grpc.ClientConn, error) {
	slf._mx.Lock()
	conns := slf._cs
	slf._mx.Unlock()

	if conns == nil {
		return nil, errGRPCSetClosed
	}
	for {
		select {
		case wrapConn := <-conns:
			if wrapConn == nil {
				return nil, errGRPCSetClosed
			}
			//判断是否超时，超时则丢弃
			if timeout := slf.IdleTimeout; timeout > 0 {
				if wrapConn._t.Add(timeout).Before(time.Now()) {
					//丢弃并关闭该链接
					slf._close(wrapConn._c)
					continue
				}
			}
			return wrapConn._c, nil
		default:
			conn, err := slf._factory()
			if err != nil {
				return nil, err
			}

			return conn, nil
		}
	}
}

func (slf *GRPCClientSet) putConn(conn *grpc.ClientConn) error {
	if conn == nil {
		return errGRPCSetRejected
	}

	slf._mx.Lock()
	defer slf._mx.Unlock()

	if slf._cs == nil {
		return slf._close(conn)
	}

	select {
	case slf._cs <- &grpcIdleConn{_c: conn, _t: time.Now()}:
		return nil
	default:
		//连接池已满，直接关闭该链接
		return slf._close(conn)
	}
}

//GetIdleCount desc
//@Method GetIdleCount desc: Number of idle connections
//@Return (int)
func (slf *GRPCClientSet) GetIdleCount() int {
	slf._mx.Lock()
	conns := slf._cs
	slf._mx.Unlock()
	return len(conns)
}
