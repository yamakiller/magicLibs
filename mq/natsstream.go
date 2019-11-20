package mq

import (
	stan "github.com/nats-io/go-nats-streaming"
)

//INatPubAck desc
//@interface INatPubAck desc: nats stream Publish ack function
//@function onAck(string, error)
type INatPubAck interface {
	onAck(string, error)
}

//INatSubCall desc
//@interface INatSubCall desc: nats stream Subscribe Call function
//@function   onRecive*(*stan.Msg)
type INatSubCall interface {
	onRecive(msg *stan.Msg)
}

//NatsStreamClient desc
//@struct NatsStreamClient
type NatsStreamClient struct {
	_c          stan.Conn
	_ack        INatPubAck
	_sub        INatSubCall
	_isShutdown bool
	_clusterID  string
	_clientID   string

	AutoReConnectLimt int
	PingInterval      int
	PingMaxOut        int
	ConnectTimeout    int
}

//Connect desc
//@method Connect desc: Connect to the NatsSteam server
//@param (string) nats server cluster ID
//@param (string) nats client ID
//@return (error) nats connect fail error
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

	return nil
}

//Publish desc
//@method Publish desc: Post message
//@param  (string) Subscription name
//@param  ([]byte) Post data
//@return (error)
func (slf *NatsStreamClient) Publish(name string, data []byte) error {
	return slf._c.Publish(name, data)
}

//PublishAsync desc
//@method PublishAsync desc: Async Post message
//@param  (string) Subscription name
//@param  ([]byte) Post data
//@return (string)
//@return (error)
func (slf *NatsStreamClient) PublishAsync(name string, data []byte) (string, error) {
	return slf._c.PublishAsync(name, data, slf._ack.onAck)
}

//Subscribe desc
//@method Subscribe desc: Recvie message
//@param (string) Subscription name
//@param (...stan.SubscriptionOption) sub option
func (slf *NatsStreamClient) Subscribe(name string, opts ...stan.SubscriptionOption) {
	slf._c.Subscribe(name, slf._sub.onRecive, opts...)
}
