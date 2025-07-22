package db

import (
	"context"
	"os"
	"rinha-de-backend-2025/internal/logger"

	"github.com/redis/go-redis/v9"
)

var (
	DB          *redis.Client
	Ctx         = context.Background()
	redisLogger = logger.GetLogger("database")
)

func StartDB() error {
	redisLogger.Info("Initialize connection DB")

	urlconn := os.Getenv("REDIS_URL")
	DB = redis.NewClient(&redis.Options{
		Addr: urlconn,
	})

	_, err := DB.Ping(Ctx).Result()
	if err != nil {
		redisLogger.Fatalf("failed to connect to Redis: %v", err)
		return err
	}

	redisLogger.Info("DB Connected!!!")
	return nil
}
