package application

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/techies/orders-api/helpers"
	"go.uber.org/zap"
)

type App struct {
	Router http.Handler
	DB     *sql.DB
	errors *AppError
	Config *Config
	Logger *zap.Logger
}

func NewApp(cfg *Config) *App {
	logger := helpers.Get()
	db, err := sql.Open("sqlite3", cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("failed to open sqlite db", zap.Error(err))
	}

	// SQLite recommended pragmas
	if _, err := db.Exec(`
	       PRAGMA foreign_keys = ON;
	       PRAGMA journal_mode = WAL;
	       PRAGMA synchronous = NORMAL;
       `); err != nil {
		logger.Fatal("failed to set sqlite pragmas", zap.Error(err))
	}

	app := &App{
		DB:     db,
		errors: &AppError{},
		Config: cfg,
		Logger: logger,
	}

	app.loadRoutes()
	return app
}

func (a *App) Start(ctx context.Context) error {
	if a.Router == nil {
		return errors.New("router is nil")
	}
	if a.DB == nil {
		return errors.New("database is nil")
	}

	// SQLite health check
	if err := a.DB.PingContext(ctx); err != nil {
		a.Logger.Error("failed to connect to sqlite", zap.Error(err))
		return fmt.Errorf("failed to connect to sqlite: %w", err)
	}
	a.Logger.Info("connected to sqlite successfully")

	defer func() {
		if err := a.DB.Close(); err != nil {
			a.Logger.Warn("sqlite close error", zap.Error(err))
		}
	}()

	port := 3000
	shutdownTimeout := 10 * time.Second
	if a.Config != nil {
		port = a.Config.Port
		if a.Config.ShutdownTimeout > 0 {
			shutdownTimeout = a.Config.ShutdownTimeout
		}
	}
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: a.Router,
	}

	errCh := make(chan error, 1)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("http server error: %w", err)

	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			a.Logger.Error("server shutdown failed", zap.Error(err))
			return fmt.Errorf("server shutdown failed: %w", err)
		}

		a.Logger.Info("server shutdown gracefully")
		return nil
	}
}
