package migrator

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Up runs all pending migrations from the given migrationsDir against dbURL.
// migrationsDir is a local folder path (e.g. "./migrations").
func Up(dbURL, migrationsDir string, logger *slog.Logger) error {
	absPath, err := filepath.Abs(migrationsDir)
	if err != nil {
		return fmt.Errorf("migrator.Up: resolve migrations path: %w", err)
	}
	// file:// expects forward slashes
	absPath = strings.ReplaceAll(absPath, "\\", "/")
	sourceURL := "file://" + absPath

	m, err := migrate.New(sourceURL, dbURL)
	if err != nil {
		return fmt.Errorf("migrator.Up: create migrate instance: %w", err)
	}
	defer func() {
		_, _ = m.Close()
	}()

	if logger != nil {
		logger.Info("Running database migrations...", "source", sourceURL)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			if logger != nil {
				logger.Info("Database is up-to-date (no migrations to apply)")
			}
			return nil
		}
		return fmt.Errorf("migrator.Up: apply migrations: %w", err)
	}

	if logger != nil {
		logger.Info("Database migrations applied successfully")
	}
	return nil
}
