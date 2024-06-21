package dbredis

import (
	"context"
	"fmt"

	"github.com/arfan21/vocagame/config"
	"github.com/arfan21/vocagame/pkg/logger"
	"github.com/redis/go-redis/v9"
)

func New() (*redis.Client, error) {
	fmt.Println("redis", config.Get().Redis)

	redisUrl := fmt.Sprintf("rediss://%s:%s@%s:%s", config.Get().Redis.Username, config.Get().Redis.Password, config.Get().Redis.URL, config.Get().Redis.Port)
	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		logger.Log(context.Background()).Error().Err(err).Msg("failed to parse redis url")
		return nil, err
	}

	client := redis.NewClient(opt)

	err = client.Ping(context.Background()).Err()
	if err != nil {
		logger.Log(context.Background()).Error().Err(err).Msg("failed to ping redis")
		return nil, err
	}

	logger.Log(context.Background()).Info().Msg("dbredis: connection established")

	return client, nil
}
