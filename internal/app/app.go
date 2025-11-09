package app

import (
	"context"
	"denet-test-task/config"
	v1 "denet-test-task/internal/api/v1"
	"denet-test-task/internal/repo"
	"denet-test-task/internal/services"
	"denet-test-task/pkg/hasher"
	"denet-test-task/pkg/httpserver"
	"denet-test-task/pkg/logctx"
	"denet-test-task/pkg/migrator"
	"denet-test-task/pkg/postgres"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-chi/chi/v5"
)

// ensureSSLMode adds sslmode=disable to the URL if sslmode parameter is not present
func ensureSSLMode(url string) string {
	if strings.Contains(url, "sslmode=") {
		return url
	}

	separator := "?"
	if strings.Contains(url, "?") {
		separator = "&"
	}

	return url + separator + "sslmode=disable"
}

func Run(configPath string) {
	// Configuration
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		slog.Error("config error", "err", err)
		os.Exit(1)
	}

	// Logger
	configureLogging()
	// root context logger
	ctx := context.Background()
	ctx = logctx.WithLogger(ctx, slog.With("app", "go-task-tracker"))
	log := logctx.FromContext(ctx)

	// DB
	log.Info("Initializing postgres...")
	dbURL := ensureSSLMode(cfg.PG.URL)
	pg, err := postgres.New(dbURL, postgres.MaxPoolSize(cfg.PG.MaxPoolSize))
	if err != nil {
		log.Error("app - Run - pgdb.NewServices", "err", err)
	}
	defer pg.Close()

	// Migrations (golang-migrate)
	log.Info("Running DB migrations...")
	if err := migrator.Up(dbURL, "./../../migrations", log); err != nil {
		log.Error("app - Run - migrator.Up", "err", err)
		os.Exit(1)
	}

	// Repositories
	log.Info("Initializing repositories...")
	repositories := repo.NewRepositories(pg)

	// Services dependencies
	log.Info("Initializing services...")
	deps := services.ServicesDependencies{
		Repos: repositories,
		// GDrive:   gdrive.New(cfg.WebAPI.GDriveJSONFilePath),
		Hasher:   hasher.NewSHA1Hasher(cfg.Hasher.Salt),
		SignKey:  cfg.JWT.SignKey,
		TokenTTL: cfg.JWT.TokenTTL,
	}
	services, err := services.NewServices(ctx, deps)
	if err != nil {
		log.Error("app - Run - services.NewServices", "err", err)
		os.Exit(1)
	}

	// Handlers
	log.Info("Initializing handlers and routes...")
	r := chi.NewRouter()
	v1.NewRouter(r, services)

	// HTTP server
	log.Info("Starting http server...")
	log.Debug("Server starting", "port", cfg.HTTP.Port)
	httpServer := httpserver.New(r, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	log.Info("Configuring graceful shutdown...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal", "signal", s.String())
	case err = <-httpServer.Notify():
		log.Error("app - Run - httpServer.Notify", "err", err)
	}

	// Graceful shutdown
	log.Info("Shutting down...")
	err = httpServer.Shutdown()
	if err != nil {
		log.Error("app - Run - httpServer.Shutdown", "err", err)
	}
}
