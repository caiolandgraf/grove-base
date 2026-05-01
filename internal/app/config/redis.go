package config

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gomodule/redigo/redis"
)

func InitRedis() (*redis.Pool, error) {
	pool := &redis.Pool{
		MaxIdle:     10,
		MaxActive:   100,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(
				"tcp",
				fmt.Sprintf("%s:%d", Env.RedisHost, Env.RedisPort),
			)
			if err != nil {
				return nil, err
			}

			if Env.RedisPassword != "" {
				if _, err := c.Do("AUTH", Env.RedisPassword); err != nil {
					_ = c.Close()
					return nil, err
				}
			}

			if _, err := c.Do("SELECT", 0); err != nil {
				_ = c.Close()
				return nil, err
			}

			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	// Test connection
	conn := pool.Get()
	defer func() {
		_ = conn.Close()
	}()

	if _, err := conn.Do("PING"); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	slog.Info("Redis connected successfully",
		"host", Env.RedisHost,
		"port", Env.RedisPort,
	)

	return pool, nil
}
