package db

import (
	"context"
	"fmt"
	"os"
	"rinha-de-backend-2025/internal/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	dbLogger = logger.GetLogger("[PGSQL]")
	PGDB     *pgxpool.Pool
)

func StartPG() error {
	dbLogger.Info("Initialize connection PG")
	pgUser := os.Getenv("POSTGRES_USER")
	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	pgHost := os.Getenv("POSTGRES_HOST")
	pgPort := os.Getenv("POSTGRES_PORT")
	pgDb := os.Getenv("POSTGRES_DB")

	DBURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", pgUser, pgPassword, pgHost, pgPort, pgDb)

	poolconn, err := pgxpool.New(context.Background(), DBURL)
	if err != nil {
		dbLogger.Errorf("Failed to connect to PostgreSQL: %v", err)
		return err
	}

	PGDB = poolconn

	dbLogger.Info("PG Connected successfully!")
	return nil
}

func GetDB() *pgxpool.Pool {
	return PGDB
}
