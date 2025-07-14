package db

import (
	"context"
	"os"
	"rinha-de-backend-2025/internal/config"

	"github.com/redis/go-redis/v9"
)

var (
	DB       *redis.Client
	Ctx      = context.Background()
	dbLogger = config.GetLogger("database")
)

func StartDB() error {
	dbLogger.Info("Initialize connection DB")

	urlconn := os.Getenv("REDIS_URL")
	DB = redis.NewClient(&redis.Options{
		Addr: urlconn,
	})

	_, err := DB.Ping(Ctx).Result()
	if err != nil {
		dbLogger.Fatalf("failed to connect to Redis: %v", err)
		return err
	}

	dbLogger.Info("DB Connected!!!")
	return nil
}
