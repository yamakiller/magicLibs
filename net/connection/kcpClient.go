package connection

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/yamakiller/magicLibs/mmath"
	"github.com/yamakiller/magicLibs/util"
	"github.com/yamakiller/mgokcp/mkcp"
)

//KCPSeria KCP序列化反序列化接口
type KCPSeria interface {
	UnSeria([]byte) (interface{}, error)
	Seria(interface{}, *mkcp.KCP) (int, error)
}

//KCPClient KCP(UDP)协议客户端
type KCPClient struct {
	WriteWaitQueue int
	ReadWaitQueue  int
	RecvWndSize    int32
	SendWndSize    int32
	NoDelay        int32
	Interval       int32
	Resend         int32
	Nc             int32
	RxMinRto       int32
	FastResend     int32
	S              KCPSeria
	E              Exception
	Mtu            int

	_c    *net.UDPConn
	_id   uint32
	_kcp  *mkcp.KCP
	_sync sync.Mutex
	_addr *net.UDPAddr

	_closed  chan bool
	_sdQueue chan interface{}
	_rdQueue chan []byte

	_buffer     []byte
	_wTotal     int
	_rTotal     int
	_lastActive int64
	_wwg        sync.WaitGroup
	_rwg        sync.WaitGroup
}

//WithID 临时使用后续删除
func (slf *KCPClient) WithID(id uint32) {
	slf._id = id
}

//Connect 连接服务器
func (slf *KCPClient) Connect(addr string, timeout time.Duration) error {

	udpAddr, err := net.ResolveUDPAddr("", addr)
	if err != nil {
		return err
	}

	if slf.Mtu == 0 {
		slf.Mtu = 1400
	}

	c, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return err
	}

	//握手获取ID

	slf._addr = udpAddr
	slf._sdQueue = make(chan interface{}, slf.WriteWaitQueue)
	slf._rdQueue = make(chan []byte, slf.ReadWaitQueue)
	if slf._closed == nil {
		slf._closed = make(chan bool, 1)
	}

	slf._buffer = make([]byte, mmath.Align(uint32(slf.Mtu), 4))
	slf._kcp = mkcp.New(slf._id, slf)
	slf._kcp.WithOutput(output)
	slf._kcp.WndSize(slf.SendWndSize, slf.RecvWndSize)
	slf._kcp.NoDelay(slf.NoDelay, slf.Interval, slf.Resend, slf.Nc)
	slf._kcp.SetMTU(int32(slf.Mtu))
	if slf.RxMinRto > 0 {
		slf._kcp.SetRxMinRto(slf.RxMinRto)
	}
	if slf.FastResend > 0 {
		slf._kcp.SetFastResend(slf.FastResend)
	}
	slf._c = c

	slf._wwg.Add(1)
	slf._rwg.Add(1)
	go slf.writeServe()
	go slf.readServe()

	return nil
}

func (slf *KCPClient) writeServe() {
	defer func() {
		slf._rwg.Wait()
		close(slf._sdQueue)
		slf._sync.Lock()
		mkcp.Free(slf._kcp)
		slf._kcp = nil
		slf._sync.Unlock()
	}()

	var current int64
	for {
	active:
		current = util.Timestamp()
		select {
		case <-slf._closed:
			goto exit
		case <-time.After(time.Duration(slf.Interval) * time.Millisecond):
		case msg, ok := <-slf._sdQueue:
			if !ok {
				goto exit
			}

			slf._sync.Lock()
			n, err := slf.S.Seria(msg, slf._kcp)
			slf._sync.Unlock()
			if err != nil {
				if slf.E != nil {
					slf.E.Error(err)
				}

				goto active
			}

			slf._wTotal += n
		}

		slf._sync.Lock()
		slf._kcp.Update(uint32(current & 0xFFFFFFFF))
		slf._sync.Unlock()
	}
exit:
}

func (slf *KCPClient) readServe() {
	defer func() {
		close(slf._rdQueue)
		slf._rwg.Done()
	}()

	for {
	active:
		select {
		case <-slf._closed:
			goto exit
		default:
			n, _, err := slf._c.ReadFromUDP(slf._buffer)
			if err != nil {
				if slf.E != nil {
					slf.E.Error(err)
				}
				goto active
			}

			slf._sync.Lock()
			slf._kcp.Input(slf._buffer, int32(n))
			slf._sync.Unlock()
			slf._rTotal += n
			slf._lastActive = util.Timestamp()

			for {
				slf._sync.Lock()
				n = int(slf._kcp.Recv(slf._buffer, int32(len(slf._buffer))))
				slf._sync.Unlock()
				if n < 0 {
					break
				}
				//需要修改池化
				tmpByte := make([]byte, n)
				copy(tmpByte, slf._buffer[:n])
				slf._rdQueue <- tmpByte
			}
		}
	}
exit:
}

//Parse 解析数据, 需要修改
func (slf *KCPClient) Parse() (interface{}, error) {
	select {
	case <-slf._closed:
		return nil, errors.New("closed")
	case data, ok := <-slf._rdQueue:
		if !ok {
			return nil, errors.New("closed")
		}

		msg, err := slf.S.UnSeria(data)
		if err != nil {
			return nil, err
		}

		return msg, nil
	}
}

//SendTo 发送数据
func (slf *KCPClient) SendTo(msg interface{}) error {
	select {
	case <-slf._closed:
		return errors.New("closed")
	default:
	}

	slf._sdQueue <- msg
	return nil
}

//Close 关闭
func (slf *KCPClient) Close() error {
	if slf._closed != nil {
		select {
		case <-slf._closed:
		default:
			close(slf._closed)
		}
	}

	err := slf._c.Close()
	slf._wwg.Wait()
	slf._rwg.Wait()

	return err
}

func output(buff []byte, user interface{}) int32 {
	client := user.(*KCPClient)
	client._c.Write(buff)
	fmt.Println(len(buff))
	return 0
}
