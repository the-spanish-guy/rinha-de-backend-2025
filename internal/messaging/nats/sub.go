package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"rinha-de-backend-2025/internal/db"
	"rinha-de-backend-2025/internal/types"

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

func (s *Subscriber) Subscribe() error {
	// Validar se existe conexão ativa?

	var err error
	s.subscription, err = s.conn.Subscribe("pub.payments", func(msg *nats.Msg) {
		s.handleMessage(msg)
	})

	if err != nil {
		return fmt.Errorf("falha ao se inscrever no subject: %w", err)
	}

	logger.Debugf("Inscrito no subject: %s", "pub.payments")
	return nil
}

func (s *Subscriber) handleMessage(msg *nats.Msg) {
	message, err := types.MessageFromJSON(msg.Data)
	if err != nil {
		logger.Errorf("Falha ao deserializar mensagem: %v", err)
		return
	}

	logger.Infof("Mensagem recebida: %s", message.Content)
	// parse the message to a struct
	payment := types.PaymentsRequest{}
	err = json.Unmarshal([]byte(message.Content), &payment)
	if err != nil {
		logger.Errorf("Falha ao deserializar mensagem: %v", err)
		return
	}

	logger.Debug("format data to save on db")

	paymentDB := types.Payments{
		CorrelationId: payment.CorrelationId,
		Amount:        payment.Amount,
		Status:        "pending",
	}

	logger.Debugf("paymentDB: %v", paymentDB)

	// processar a mensagem
	// criar um /service/process-payments ? para lidar com a regra de negocio
	pgdb := db.GetDB()
	if pgdb == nil {
		logger.Error("Database connection is nil")
		return
	}

	_, err = pgdb.Exec(context.Background(),
		"INSERT INTO payments (correlation_id, amount, status) VALUES ($1, $2, $3)",
		paymentDB.CorrelationId, paymentDB.Amount, paymentDB.Status)

	if err != nil {
		logger.Errorf("Falha ao inserir payment no banco: %v", err)
		return
	}

	logger.Info("Payment salvo com sucesso no banco de dados")
}
