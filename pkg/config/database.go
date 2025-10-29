package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// O OpenDB abre uma conexão GORM PostgreSQL usando variáveis ​​de ambiente (carregadas de .env pelo aplicativo).
// Variáveis ​​de ambiente obrigatórias: POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB
// Variáveis ​​de ambiente opcionais com padrões: POSTGRES_HOST=localhost, POSTGRES_PORT=5432
func OpenDB() (*gorm.DB, error) {
	_ = godotenv.Load()

	host := getenvOr("POSTGRES_HOST", "localhost")
	port := getenvOr("POSTGRES_PORT", "5432")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	if user == "" || password == "" || dbname == "" {
		return nil, fmt.Errorf("missing required postgres env vars (POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB)")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		host, port, user, password, dbname,
	)

	// configure GORM logger (info by default)
	newLogger := gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gormlogger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  gormlogger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func getenvOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
