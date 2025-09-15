package utils

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func GetCacheRedis(ctx context.Context, rdb *redis.Client, key string, dest any) (bool, error) {
	cmd := rdb.Get(ctx, key)
	if cmd.Err() != nil {
		if cmd.Err() == redis.Nil {
			log.Printf("Key %s does not exist\n", key)
			return false, nil
		}
		log.Println("Redis Error\nCause:", cmd.Err().Error())
		return false, cmd.Err()
	}

	data, err := cmd.Bytes()
	if err != nil {
		log.Println("Redis Bytes Error\nCause:", err.Error())
		return false, err
	}

	if err := json.Unmarshal(data, dest); err != nil {
		log.Println("Redis Unmarshal Error\nCause:", err.Error())
		return false, err
	}

	return true, nil
}

func SetCacheRedis(ctx context.Context, rdb *redis.Client, key string, value any, time time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		log.Println("Redis Marshal Error\nCause:", err.Error())
		return err
	}
	if err := rdb.Set(ctx, key, data, time).Err(); err != nil {
		log.Println("Redis Set Error\nCause:", err.Error())
		return err
	}
	return nil
}
