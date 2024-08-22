package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Collap5e-dev/FlickPick/internal/config"
)

func main() {
	err := ApplyDump(context.Background())
	if err != nil {
		fmt.Printf("\nSTOP with error: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("DONE")
}

func ApplyDump(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config.Load() err: %w", err)
	}

	psqlInfo := cfg.Db.Dsn()
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		return fmt.Errorf("Error connecting to the database: %w", err)
	}
	defer db.Close()

	dumpPath := strings.Replace(cfg.MigrationPath, "pg", "dump", 1)

	sqlFiles, err := getSQLFiles(dumpPath)
	if err != nil {
		return fmt.Errorf("Error getting SQL files: %w", err)
	}

	err = executeSQLFiles(ctx, db, sqlFiles)
	if err != nil {
		return fmt.Errorf("Error executing SQL files: %w", err)
	}

	return nil
}

func getSQLFiles(dumpPath string) ([]string, error) {
	files, err := os.ReadDir(dumpPath)
	if err != nil {
		return nil, fmt.Errorf("Error reading dump directory: %w", err)
	}

	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, filepath.Join(dumpPath, file.Name()))
		}
	}

	sort.Strings(sqlFiles)

	return sqlFiles, nil
}

func executeSQLFiles(ctx context.Context, db *sqlx.DB, sqlFiles []string) error {
	for _, file := range sqlFiles {
		err := executeSQLFile(ctx, db, file)
		if err != nil {
			return fmt.Errorf("Error executing file %s: %w", file, err)
		}
	}

	return nil
}

func executeSQLFile(ctx context.Context, db *sqlx.DB, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("Error reading file %s: %w", filePath, err)
	}

	_, err = db.ExecContext(ctx, string(content))
	if err != nil {
		return fmt.Errorf("Error executing SQL in file %s: %w", filePath, err)
	}

	fmt.Printf("Executed %s successfully\n", filePath)
	return nil
}
