package listener

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

//WSSListener Web Socket Listener
type WSSListener struct {
	Listener *websocket.Upgrader
	_wg      sync.WaitGroup
}

//Accept 接受链接
func (slf *WSSListener) Accept(w http.ResponseWriter, r *http.Request) (*WSSConn, error) {
	c, err := slf.Listener.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	slf._wg.Add(1)
	return &WSSConn{Conn: c, _wg: &slf._wg}, nil
}

//Wait 等待所有客户端结束
func (slf *WSSListener) Wait() {
	slf._wg.Wait()
}

//WSSConn WebSocket Accept client
type WSSConn struct {
	*websocket.Conn
	_wg *sync.WaitGroup
}

//Close ...
func (slf *WSSConn) Close() error {
	if slf._wg != nil {
		slf._wg.Done()
	}
	slf._wg = nil
	return slf.Conn.Close()
}
