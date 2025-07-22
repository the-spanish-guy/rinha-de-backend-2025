package nats

import (
	"rinha-de-backend-2025/internal/logger"
	"time"

	"github.com/nats-io/nats-server/v2/server"
)

var log = logger.GetLogger("[NATS]")
var default_host = "nats://localhost:4222"

type Server struct {
	natsServer *server.Server
	log        *logger.Logger
}

func NewServer(l *logger.Logger) *Server {
	return &Server{
		log: l,
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
		ns.log.Errorf("An error occurred when create NATS server")
		return err
	}

	ns.log.Infof("Starting NATS server on %s:%d", opts.Host, opts.Port)

	go ns.natsServer.Start()

	if !ns.natsServer.ReadyForConnections(5 * time.Second) {
		return err
	}

	ns.log.Info("Server NATS initiated")
	return nil
}
