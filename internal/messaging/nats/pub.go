package nats

import (
	"fmt"
	"rinha-de-backend-2025/internal/types"

	"github.com/nats-io/nats.go"
)

type Publisher struct {
	conn *nats.Conn
}

func NewPublisher() *Publisher {
	return &Publisher{}
}

func (p *Publisher) Connect() error {
	var err error
	p.conn, err = nats.Connect(default_host)
	if err != nil {
		return fmt.Errorf("fAn error occurred while trying to connect NATS: %w", err)
	}

	log.Infof("Publisher connected to NATS %s", default_host)

	return nil
}

func (p *Publisher) PublishMessage(message *types.Message) error {
	data, err := message.ToJSON()
	if err != nil {
		log.Errorf("An error occurred trying parsing JSON")
		return err
	}

	err = p.conn.Publish("pub.payments", data)
	if err != nil {
		log.Errorf("Error on publish: %v", err)
		return err
	}

	log.Infof("Message published: %s", message.Content)
	return nil
}
