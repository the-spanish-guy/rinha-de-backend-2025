package nats

import (
	"rinha-de-backend-2025/internal/config"
	"time"

	"github.com/nats-io/nats-server/v2/server"
)

var logger = config.GetLogger("[NATS]")
var default_host = "nats://localhost:4222"

type Server struct {
	natsServer *server.Server
	logger     *config.Logger
}

func NewServer(l *config.Logger) *Server {
	return &Server{
		logger: l,
	}
}

func (ns *Server) Start() error {
	opts := &server.Options{
		ServerName: "payments_server",
		DontListen: false,
	}

	var err error
	ns.natsServer, err = server.NewServer(opts)

	if err != nil {
		ns.logger.Errorf("An error occurred when create NATS server")
		return err
	}

	ns.logger.Infof("Starting NATS server on %s:%d", opts.Host, opts.Port)

	go ns.natsServer.Start()

	if !ns.natsServer.ReadyForConnections(5 * time.Second) {
		return err
	}

	ns.logger.Info("Server NATS initiated")
	return nil
}
