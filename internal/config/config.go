package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Db struct {
	Host   string `envconfig:"host" json:"host"`
	Port   int64  `envconfig:"port" json:"port"`
	User   string `envconfig:"user" json:"user"`
	Pass   string `envconfig:"pass" json:"pass"`
	DbName string `envconfig:"dbname" json:"dbname"`
}

func (d Db) Dsn() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		d.User,
		d.Pass,
		d.Host,
		d.Port,
		d.DbName,
	)
}

func (d Db) PgFormat() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		d.Host,
		d.Port,
		d.User,
		d.Pass,
		d.DbName,
	)
}

type Config struct {
	Db            Db
	Port          int64
	MigrationPath string `envconfig:"migration_path" json:"migration_path"`
	SecretKey     string `envconfig:"secret_key" json:"secret_key"`
}

func Load() (*Config, error) {
	err := initDotEnv()
	if err != nil {
		return nil, fmt.Errorf("initDotEnv err: %w", err)
	}
	var cfg Config
	err = envconfig.Process("app", &cfg)
	if err != nil {
		return nil, fmt.Errorf("envconfig.Process err: %w", err)
	}
	return &cfg, nil
}

func initDotEnv() error {
	if !fileExists(".env") {
		return nil
	}

	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("gotenv.Load err: %w", err)
	}
	return nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	// If there's an error that isn't about the file not existing,
	// it might be another kind of issue.
	return false
}
