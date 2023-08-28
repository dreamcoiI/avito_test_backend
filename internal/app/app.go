package app

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dreamcoiI/avito_test_backend/api"
	"github.com/dreamcoiI/avito_test_backend/api/middleware"
	"github.com/dreamcoiI/avito_test_backend/internal/config"
	"github.com/dreamcoiI/avito_test_backend/internal/handlers"
	service2 "github.com/dreamcoiI/avito_test_backend/internal/service"
	"github.com/dreamcoiI/avito_test_backend/internal/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

type Server struct {
	config     config.Config
	context    context.Context
	server     *http.Server
	storage    *sql.DB
	middleware *zerolog.Logger
}

func NewServer(config config.Config, context context.Context, middleware *zerolog.Logger) *Server {
	server := new(Server)
	server.config = config
	server.context = context
	server.middleware = middleware
	return server
}

func (app *Server) Start() error {
	app.middleware.Info().Msg("Starting server")
	app.middleware.Info().Msg(app.config.GetDBString())

	var err error

	//app.storage, err = pgxpool.Connect(app.context, app.config.GetDBString())
	app.storage, err = sql.Open("pgx", app.config.GetDBString())

	if err != nil {
		fmt.Println("Failed to connect to the database")
		app.middleware.Fatal().Err(err).Msg("Failed to connect to the database")
	}

	storageInstance := storage.NewStorage(app.storage)
	service := service2.NewService(storageInstance)
	orderHandler := handlers.NewHandler(service)

	router := api.ConfigureRouters(orderHandler)
	router.Use(middleware.LogRequest)

	app.server = &http.Server{
		Addr:    "0.0.0.0:" + app.config.Port,
		Handler: router,
	}

	app.middleware.Info().Msg("Server started.")

	err = app.server.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		app.middleware.Err(err).Msg("Failed while serving")
		return err
	}

	return nil
}

func (app *Server) Shutdown() error {
	app.middleware.Info().Msg("Server stopped")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	app.storage.Close()

	if err := app.server.Shutdown(ctxShutdown); err != nil {
		app.middleware.Err(err)
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	app.middleware.Info().Msg("Server shutdown completed")

	return nil
}
