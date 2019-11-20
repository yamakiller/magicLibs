package mq

import (
	stan "github.com/nats-io/go-nats-streaming"
)

//NatsStreamClient desc
//@struct NatsStreamClient
type NatsStreamClient struct {
	_c          stan.Conn
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
