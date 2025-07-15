package messaging

import (
	"rinha-de-backend-2025/internal/config"
	"rinha-de-backend-2025/internal/messaging/nats"
)

func SetupMessaging() (pub *nats.Publisher, sub *nats.Subscriber) {
	logger := config.GetLogger("MSG")

	// start nats server
	ns := nats.NewServer(logger)
	if err := ns.Start(); err != nil {
		logger.Fatal("An error occurred when NATS server starting", err.Error())
	}

	// create subscribe
	sub = nats.NewSubscriber()
	// connect sub to nats server
	if err := sub.Connect(); err != nil {
		logger.Fatal("An error occurred when create subscriber server", err.Error())
	}

	//create publisher
	pub = nats.NewPublisher()
	// connect pub to nats server
	if err := pub.Connect(); err != nil {
		logger.Fatal("An error occurred when create publisher server", err.Error())
	}

	return pub, sub
}
