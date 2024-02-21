package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const defaultReconnectionTimeout = 5

func NewRedisClient(host string, port string, logger *zap.Logger) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", host, port)
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx := context.Background()

	_, err := client.Ping(ctx).Result()
	if err == nil {
		return client, nil
	}

	logger.Error("error when starting redis server", zap.Error(err))

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	timeoutExceeded := time.After(time.Second * time.Duration(defaultReconnectionTimeout))

	for {
		select {
		case <-timeoutExceeded:
			return nil, fmt.Errorf("cache connection failed after %d timeout:%s", defaultReconnectionTimeout, err.Error())
		case <-ticker.C:
			err := client.Ping(ctx).Err()
			if err == nil {
				return client, nil
			}

			logger.Error("error when starting redis server", zap.Error(err))
		}
	}

}
