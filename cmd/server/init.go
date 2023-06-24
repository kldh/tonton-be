package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ldhk/tonton-be/pkg/openapi"
	"github.com/ldhk/tonton-be/pkg/telemetry/logging"
	"github.com/ldhk/tonton-be/pkg/telemetry/logging/zap"
	"github.com/ldhk/tonton-be/pkg/telemetry/monitor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *Server) init() {
	s.initLogging()

	checkErr := func(name string, err error) {
		if err != nil {
			s.l.Fatalf("init %s failed: %+v", name, err)
		}
	}

	checkErr("module", s.initModule())
	s.initHTTP()
}

func (s *Server) start() {
	s.l.Infof("http.server: listening on %s", s.http.Addr)
	if err := s.http.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		s.l.Infof("http server: ListenAndServe failed: %v", err)
		return
	}
}

func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.l.Info("http server: shutting down")
	if err := s.http.Shutdown(ctx); err != nil {
		s.l.Warnf("http server: shutdown failed: %v", err)
	}
}

func (s *Server) initLogging() {
	c := zap.LocalConfig()
	if s.config.Env != "local" {
		c = zap.ReleaseConfig()
	}

	l := zap.New(c)
	logging.SetDefaultLogger(l)

	s.l = logging.FromContext(context.Background())
}

func (s *Server) initHTTP() {
	gin.SetMode(gin.ReleaseMode)

	s.l.Info("http.server: initializing")
	e := gin.New()
	e.GET("/ping", func(c *gin.Context) {
		c.Writer.WriteHeader(200)
		c.Writer.WriteString("pong")
	})
	e.GET("/health", func(c *gin.Context) { c.Writer.WriteHeader(200) })
	e.GET("/info", func(c *gin.Context) { c.Writer.WriteHeader(200) })
	e.GET("/metrics", gin.WrapH(promhttp.Handler()))
	monitor.Gin(e, s.config.Tracing.ServiceName)
	e.Use(gin.Recovery())

	s.module.openAPI.Route(e)

	s.http = &http.Server{
		Addr:              ":8080",
		Handler:           e,
		ReadHeaderTimeout: 60 * time.Second,
	}
}

func (s *Server) initModule() error {
	openApi, err := openapi.InitModule(openapi.Config{})
	if err != nil {
		return err
	}
	s.module.openAPI = openApi
	return nil
}
