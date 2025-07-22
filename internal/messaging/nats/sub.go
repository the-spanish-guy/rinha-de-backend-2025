package nats

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"rinha-de-backend-2025/internal/config"
	"rinha-de-backend-2025/internal/db"
	"rinha-de-backend-2025/internal/types"
	"strings"
	"time"

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
		log.Fatalf("An error occurred while trying to connect NATS: %v", err)
		return err
	}

	log.Info("Subscriber connected to NATS")
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

	log.Debugf("Inscrito no subject: %s", "pub.payments")
	return nil
}

func (s *Subscriber) handleMessage(msg *nats.Msg) {
	log.Info("Receive message")
	message, err := types.MessageFromJSON(msg.Data)
	if err != nil {
		log.Errorf("Falha ao deserializar mensagem: %v", err)
		return
	}
	log.Infof("Formatted message: %s", message.Content)

	pm := config.NewProcessorManager()
	activeHost := pm.GetActiveProcessor()
	processorType := getProcessorType(activeHost)

	payment := types.PaymentsRequest{}
	err = json.Unmarshal([]byte(message.Content), &payment)
	if err != nil {
		log.Errorf("Falha ao deserializar mensagem: %v", err)
		return
	}

	payment.RequestedAt = time.Now()
	client := &http.Client{
		Timeout: 1500 * time.Millisecond,
	}

	request, _ := json.Marshal(payment)
	_, errPayments := client.Post(activeHost+"/payments", "application/json", bytes.NewBuffer(request))
	if errPayments != nil {
		log.Errorf("Error on POST /payments %s", errPayments)
		// TODO: implementar retry
	}

	paymentDB := types.Payments{
		CorrelationId: payment.CorrelationId,
		Amount:        payment.Amount,
		RequestedAt:   payment.RequestedAt,
		Status:        "PENDING",
		Processor:     processorType,
	}

	pgdb := db.GetDB()
	if pgdb == nil {
		log.Error("Database connection is nil")
		return
	}

	_, err = pgdb.Exec(context.Background(),
		"INSERT INTO payments (correlation_id, amount, status, processor, requested_at) VALUES ($1, $2, $3, $4, $5)",
		paymentDB.CorrelationId, paymentDB.Amount, paymentDB.Status, paymentDB.Processor, paymentDB.RequestedAt)

	if err != nil {
		log.Errorf("Falha ao inserir payment no banco: %v", err)
		return
	}

	log.Info("Payment salvo com sucesso no banco de dados")
}

func getProcessorType(host string) string {
	res := "DEFAULT"

	if strings.Contains(host, "fallback") {
		res = "FALLBACK"
	}

	return res
}
