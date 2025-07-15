package nats

import (
	"github.com/nats-io/nats.go"
)

type Subscriber struct {
	conn         *nats.Conn
	subscription *nats.Subscription
}

func NewSubscriber() *Subscriber {
	return &Subscriber{}
}

func (s *Subscriber) Connect() error {
	var err error
	s.conn, err = nats.Connect(default_host)
	if err != nil {
		logger.Fatalf("An error occurred while trying to connect NATS: %v", err)
		return err
	}

	logger.Info("Subscriber connected to NATS")
	return nil
}
