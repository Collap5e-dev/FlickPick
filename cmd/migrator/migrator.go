package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/Collap5e-dev/FlickPick/internal/config"
)

const (
	ActionUp   = "up"
	ActionDown = "down"
)

func main() {
	action := ActionUp
	if len(os.Args) > 1 {
		action = os.Args[1]
	}
	flag.Parse()

	err := RunMigrator(context.Background(), action)
	if err != nil {
		fmt.Printf("\nSTOP with error: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("DONE")
}

func RunMigrator(ctx context.Context, action string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config.Load() err: %w", err)
	}

	err = DoMigration(cfg, action)
	if err != nil {
		return fmt.Errorf("DoMigration() err: %w", err)
	}
	return err
}

func DoMigration(cfg *config.Config, action string) error {
	return doMigration(cfg.MigrationPath, cfg.Db.Dsn(), action)
}

func doMigration(source string, dbDsn string, action string) error {
	sourceURL := "file://" + source

	m, err := migrate.New(sourceURL, dbDsn)
	if err != nil {
		return fmt.Errorf("cant migrate.New, err: %w", err)
	}

	switch action {
	case ActionUp:
		fmt.Println("Up in process...")
		err = m.Up()
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	case ActionDown:
		fmt.Println("Down in process...")
		err = m.Steps(-1)
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	default:
		return fmt.Errorf("supported actions: [%v]", []string{ActionUp, ActionDown})
	}

	return nil
}
