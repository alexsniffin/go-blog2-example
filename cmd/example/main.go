package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/alexsniffin/go-blog2-example/internal/example/config"
	"github.com/alexsniffin/go-blog2-example/internal/example/logger"
	"github.com/alexsniffin/go-blog2-example/internal/example/models"
	"github.com/alexsniffin/go-blog2-example/internal/example/server"
)

const (
	configFileName = "example"
	prefix         = "EXAMPLE"
)

func main() {
	newCfg := models.Config{}
	err := config.NewConfig(configFileName, prefix, &newCfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	newLogger, err := logger.NewLogger(newCfg.Logger)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	newLogger.Info().Msg("setting up service")
	newServer, err := server.NewServer(newCfg, newLogger)
	if err != nil {
		newLogger.Panic().Err(err).Msg("failed to init server")
	}

	go newServer.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	stopped := <-stop
	newLogger.Info().Msg(stopped.String() + " signal received")
	newServer.Shutdown()
	newLogger.Info().Msg("service shutdown")
}
