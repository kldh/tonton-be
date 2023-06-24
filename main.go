package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ldhk/tonton-be/cmd/server"
	"github.com/ldhk/tonton-be/pkg/config"
)

func main() {
	// cp := os.Getenv("CONFIG_PATH")
	cp := "./config/dev.yaml"
	if cp == "" {
		log.Fatalf("Config path is empty")
	}

	var c server.Config
	if err := config.Load(cp, &c); err != nil {
		log.Fatalf("Load config from %q failed: %v", cp, err)
	}

	s := server.NewServer(c)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, os.Interrupt)
	go s.Start()

	<-shutdown
	s.Shutdown()
}
