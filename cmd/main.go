package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/human9001/teams/config"
	container "github.com/human9001/teams/internal/di"
	"github.com/human9001/teams/pkg/closer"
)

type App struct {
	diContainer *container.Container
	httpServer  *http.Server
}

func main() {
	if err := config.LoadEnvFile(); err != nil {
		slog.Info("не удалось загрузить .env файл", "warn", err)
	}
	configPath := config.ResolveConfigPath()

	config.MustLoad(configPath)

	di, err := container.NewDIContainer()
	if err != nil {
		slog.Error("ошибка инициализации DI контейнера", "error", err)
		os.Exit(1)
	}
	app := App{diContainer: di}
	app.initHTTPServer()
	err = app.run()
	if err != nil {
		os.Exit(1)
	}
}

func (a *App) run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		errCh <- a.runHTTPServer()
	}()

	var runErr error
	select {
	case runErr = <-errCh:
	case <-ctx.Done():
		slog.Info("получен сигнал завершения, начинаем graceful shutdown")
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := closer.CloseAll(shutdownCtx); err != nil {
		slog.Error("ошибка при завершении работы", "error", err)
		if runErr == nil {
			runErr = err
		}
	}

	return runErr
}

func (a *App) initHTTPServer() {
	a.httpServer = &http.Server{
		Addr:              config.AppConfig().HTTP.Address(),
		Handler:           a.diContainer.Router,
		ReadHeaderTimeout: config.AppConfig().HTTP.Options.HttpReadHeaderTimeout,
		ReadTimeout:       config.AppConfig().HTTP.Options.HttpReadTimeout,
		WriteTimeout:      config.AppConfig().HTTP.Options.HttpWriteTimeout,
		IdleTimeout:       config.AppConfig().HTTP.Options.HttpIdleTimeout,
	}

	closer.Add("HTTP server", func(_ context.Context) error {
		shutdownCtx, stop := context.WithTimeout(
			context.Background(),
			config.AppConfig().HTTP.Options.HttpShutdownTimeout,
		)
		defer stop()
		return a.httpServer.Shutdown(shutdownCtx)
	})
}

func (a *App) runHTTPServer() error {
	slog.Info("🚀 HTTP-сервер запущен", "address", config.AppConfig().HTTP.Address())
	err := a.httpServer.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		slog.Info("завершение HHTP сервера ")
		return nil
	}
	return err
}
