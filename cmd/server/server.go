package server

import (
	"net/http"

	"github.com/ldhk/tonton-be/pkg/telemetry/logging"
)

type Config struct {
	Env     string
	Tracing struct {
		ServiceName string
	}
}

type Server struct {
	config Config

	http *http.Server

	l logging.Logger

	queue struct{}

	client struct{}

	module struct{}

	database struct{}
}

func NewServer(c Config) *Server {
	return &Server{config: c}
}

func (s *Server) Start() {
	s.init()
	s.start()
}
