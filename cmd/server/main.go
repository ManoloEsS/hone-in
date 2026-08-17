package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ManoloEsS/hone-in/internal/database"
	"github.com/ManoloEsS/hone-in/internal/recipe"
	"github.com/ManoloEsS/hone-in/internal/web"
)

const (
	defaultAddress    = ":8080"
	shutdownTimeout   = 10 * time.Second
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 10 * time.Second
	idleTimeout       = 60 * time.Second
	maxHeaderBytes    = 1 << 20
)

func newHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(response http.ResponseWriter, _ *http.Request) {
		response.WriteHeader(http.StatusOK)
		response.Write([]byte("healthcheck"))
	})

	return mux
}

func newServer(address string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              address,
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}
}

func serverAddress(configured string) string {
	if configured == "" {
		return defaultAddress
	}
	return configured
}

func runWithContext(ctx context.Context, address string, handler http.Handler, listen func(string, string) (net.Listener, error)) error {
	address = serverAddress(address)
	listener, err := listen("tcp", address)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", address, err)
	}

	return serve(ctx, newServer(address, handler), listener, shutdownTimeout)
}

func run() error {
	databasePath := os.Getenv("DB_PATH")
	if databasePath == "" {
		return errors.New("DB_PATH is required")
	}

	db, err := database.Open(databasePath)
	if err != nil {
		return fmt.Errorf("open application database: %w", err)
	}
	defer db.Close()

	handler, err := web.NewHandler(recipe.NewRepository(db))
	if err != nil {
		return fmt.Errorf("create application handler: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	return runWithContext(ctx, os.Getenv("HTTP_ADDR"), handler, net.Listen)
}

func main() {
	if err := run(); err != nil {
		slog.Error("server failed", slog.Any("error", err))
		os.Exit(1)
	}
}

func serve(ctx context.Context, server *http.Server, listener net.Listener, shutdownTimeout time.Duration) error {
	serveErrors := make(chan error, 1)
	go func() {
		serveErrors <- server.Serve(listener)
	}()

	select {
	case <-ctx.Done():
		shutdownContext, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownContext); err != nil {
			return err
		}

		err := <-serveErrors
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case err := <-serveErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}
