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
	// Путь к файлу с переменными окружения
	envFilePath := "path/to/your/env/file.env"

	// Загрузка конфигурации из файла с переменными окружения
	config := config2.LoadConfigFromEnv(envFilePath)

	// Создание логгера
	logger := new(zerolog.Logger)

	// Инициализация контекста и отмены
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Обработка сигналов завершения
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

	// Создание сервера

	// Запуск сервера
	if err := server.Start(); err != nil {
		logger.Err(err)
	}
}
