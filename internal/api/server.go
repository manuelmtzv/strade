package api

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strade/internal/config"
	"syscall"
	"time"

	"go.uber.org/zap"
)

type Application struct {
	Config config.APIConfig
	Router http.Handler
	Logger *zap.SugaredLogger
}

func NewApplication(config config.APIConfig, logger *zap.SugaredLogger) *Application {
	return &Application{
		Config: config,
		Logger: logger,
	}
}

func (app *Application) SetRouter(router http.Handler) {
	app.Router = router
}

func (app *Application) ServeHTTP() error {
	server := &http.Server{
		Addr:         app.Config.Addr,
		Handler:      app.Router,
		WriteTimeout: time.Second * 45,
		ReadTimeout:  time.Second * 20,
		IdleTimeout:  time.Minute,
	}

	shutdownErr := make(chan error, 1)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit
		app.Logger.Infof("received signal: %s, shutting down server...", s)

		ctx, cancel := context.WithTimeout(context.Background(), app.Config.ShutdownTimeout)
		defer cancel()

		shutdownErr <- server.Shutdown(ctx)
	}()

	app.Logger.Infow("starting server", "addr", app.Config.Addr)

	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	if err := <-shutdownErr; err != nil {
		return err
	}

	app.Logger.Infow("stopped server", "addr", app.Config.Addr)
	return nil
}
