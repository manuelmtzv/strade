package api

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strade/internal/config"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
)

type Application struct {
	Config config.APIConfig
	Router http.Handler
	Wg     sync.WaitGroup
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

	shutdownErrorStream := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		app.Logger.Infof("received signal: %s, shutting down server...", s)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			shutdownErrorStream <- err
			return
		}

		app.Logger.Infow("completing background tasks", "addr", app.Config.Addr)

		app.Wg.Wait()
		shutdownErrorStream <- nil
	}()

	app.Logger.Infow("starting server", "addr", app.Config.Addr)

	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownErrorStream
	if err != nil {
		return err
	}

	app.Logger.Infow("stopped server", "addr", app.Config.Addr)

	return nil
}
