package db

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	DB  *redis.Client
	Ctx = context.Background()
)

func StartDB() error {
	urlconn := os.Getenv("REDIS_URLAS")
	DB = redis.NewClient(&redis.Options{
		Addr: urlconn,
	})

	_, err := DB.Ping(Ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %v", err)
	}

	fmt.Print("DB Connected!!!")
	return nil
}
