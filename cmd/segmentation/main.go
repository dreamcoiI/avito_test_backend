package main

import (
	"context"
	"fmt"
	"github.com/dreamcoiI/avito_test_backend/internal/app"
	config2 "github.com/dreamcoiI/avito_test_backend/internal/config"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	envFilePath := "config.env"

	config := config2.LoadConfigFromEnv(envFilePath)

	logger := new(zerolog.Logger)

	ctx, cancel := context.WithCancel(context.Background())

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	server := app.NewServer(config, ctx, logger)
	go func() {
		oscall := <-signalChannel
		logger.Info().Msg(fmt.Sprintf("Received system call: %+v", oscall))
		if err := server.Shutdown(); err != nil {
			logger.Err(err)
		}
		cancel()
	}()

	if err := server.Start(); err != nil {
		logger.Err(err)
	}
}
