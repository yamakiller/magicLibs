package mq

import (
	"time"

	stan "github.com/nats-io/go-nats-streaming"
)

//INatPubAck doc
//@Interface INatPubAck @Summary nats stream Publish ack function
//@function onAck(string, error)
type INatPubAck interface {
	onAck(string, error)
}

//INatSubCall doc
//@Interface INatSubCall @Summary nats stream Subscribe Call function
//@function   onRecive*(*stan.Msg)
type INatSubCall interface {
	onRecive(msg *stan.Msg)
}

//INatLostCall doc
//@Interface INatLostCall @Summary nats stream connection lost Call function
//@function onLost(stan.Conn, error)
type INatLostCall interface {
	onLost(stan.Conn, error)
}

//NatsStreamClient doc
//@Struct NatsStreamClient
type NatsStreamClient struct {
	_c          stan.Conn
	_ack        INatPubAck
	_sub        INatSubCall
	_lost       INatLostCall
	_isShutdown bool
	_clusterID  string
	_clientID   string

	PingInterval   int
	PingMaxOut     int
	ConnectTimeout int
}

//Connect doc
//@Method Connect @Summary Connect to the NatsSteam server
//@Param (string) nats server cluster ID
//@Param (string) nats client ID
//@Return (error) nats connect fail error
func (slf *NatsStreamClient) Connect(clusterID string, clientID string) error {
	slf._clusterID = clusterID
	slf._clientID = clientID
	if slf.PingInterval == 0 {
		slf.PingInterval = stan.DefaultPingInterval
	}

	if slf.PingMaxOut == 0 {
		slf.PingMaxOut = stan.DefaultPingMaxOut
	}

	if slf.ConnectTimeout == 0 {
		slf.ConnectTimeout = 2
	}

	c, e := stan.Connect(clusterID,
		clientID,
		stan.Pings(slf.PingInterval, slf.PingMaxOut),
		stan.ConnectWait(time.Duration(slf.ConnectTimeout)*time.Millisecond),
		stan.SetConnectionLostHandler(slf._lost.onLost))

	if e != nil {
		return e
	}

	slf._c = c

	return nil
}

//Publish doc
//@Method Publish @Summary Post message
//@Param  (string) Subscription name
//@Param  ([]byte) Post data
//@Return (error)
func (slf *NatsStreamClient) Publish(name string, data []byte) error {
	return slf._c.Publish(name, data)
}

//PublishAsync doc
//@Method PublishAsync @Summary Async Post message
//@Param  (string) Subscription name
//@Param  ([]byte) Post data
//@Return (string)
//@Return (error)
func (slf *NatsStreamClient) PublishAsync(name string, data []byte) (string, error) {
	return slf._c.PublishAsync(name, data, slf._ack.onAck)
}

//Subscribe doc
//@Method Subscribe @Summary Recvie message
//@Param (string) Subscription name
//@Param (...stan.SubscriptionOption) sub option
func (slf *NatsStreamClient) Subscribe(name string, opts ...stan.SubscriptionOption) (stan.Subscription, error) {
	return slf._c.Subscribe(name, slf._sub.onRecive, opts...)
}

//QueueSubscribe doc
//@Method QueueSubscribe @Summary Recvie Queue message
//@Param (string) Subscription name
//@Param (string) Subscription group name
//@Param (...stan.SubscriptionOption) sub option
func (slf *NatsStreamClient) QueueSubscribe(name, qgroup string, opts ...stan.SubscriptionOption) (stan.Subscription, error) {
	return slf._c.QueueSubscribe(name, qgroup, slf._sub.onRecive, opts...)
}

//Close doc
//@Method Close @Summary Close connection
//@Return (error) close error returns
func (slf *NatsStreamClient) Close() error {
	return slf._c.Close()
}
