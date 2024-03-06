package redisdb

import (
	"context"
	"log/slog"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisDB struct {
	сlient     *redis.Client
	sessionTTL int
	logger     *slog.Logger
}

func NewRedisDB(logger *slog.Logger) *RedisDB {
	redis := &RedisDB{
		logger: logger,
	}
	if err := redis.setupConnection(); err != nil {
		logger.Error("redis.New: failed to setup connection", "err", err, os.Getenv("REDIS_CONN"))
		os.Exit(1)
	}

	ttl := os.Getenv("REFRESH_TOKEN_TTL")
	if ttl == "" {
		ttl = "0"
	}
	t, err := strconv.Atoi(ttl)
	if err != nil {
		logger.Error("failed to parse REFRESH_TOKEN_TTL", "err", err)
		os.Exit(1)
	}

	redis.sessionTTL = t

	logger.Info("redis.New: connection established")
	return redis
}

func (r *RedisDB) setupConnection() error {
	//url := os.Getenv("REDIS_CONN")
	options, err := redis.ParseURL("redis://localhost:6379")
	if err != nil {
		return err
	}

	cli := redis.NewClient(options)
	if err := cli.Ping(context.TODO()).Err(); err != nil {
		return err
	}
	r.сlient = cli
	return nil
}

func (r *RedisDB) Close() {
	r.сlient.Close()
}
