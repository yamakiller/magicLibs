package mq

import (
	stan "github.com/nats-io/go-nats-streaming"
)

//NatsStreamClient desc
//@struct NatsStreamClient
type NatsStreamClient struct {
	c          stan.Conn
	isShutdown bool
	clusterID  string
	clientID   string

	AutoReConnectLimt int
	PingInterval      int
	PingMaxOut        int
	ConnectTimeout    int
}

// Connect : xx
func (slf *NatsStreamClient) Connect(clusterID string, clientID string) error {
	slf.clusterID = clusterID
	slf.clientID = clientID
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
