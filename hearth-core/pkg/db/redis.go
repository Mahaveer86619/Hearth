package db

import (
	"context"
	"fmt"
	"time"

	"github.com/Mahaveer86619/Hearth/pkg/constants"
	"github.com/Mahaveer86619/Hearth/pkg/logger"
	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func InitRedis(url string) error {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return fmt.Errorf("invalid redis url: %w", err)
	}

	RDB = redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RDB.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	logger.Info("Redis", "Connected to Redis")
	return nil
}

func PublishLog(ctx context.Context, msg []byte) error {
	return RDB.Publish(ctx, string(constants.RedisChannelLiveLogs), msg).Err()
}
