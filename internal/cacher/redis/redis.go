package redis

import (
	"context"
	"fmt"
	"kprg/internal/models"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(connection *models.Connection) (*Redis, error) {

	addr := fmt.Sprintf("%s:%s", connection.Host, connection.Port)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	redis :=
		Redis{
			client: client,
		}

	err := redis.ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &redis, nil

}

func (r *Redis) Has(key string) (bool, error) {

	ctx := context.Background()

	val, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check if key exists in redis: %w", err)
	}

	if val > 0 {
		return true, nil
	}

	return false, nil

}

func (r *Redis) Set(key string, value interface{}) error {

	ctx := context.Background()

	err := r.client.Set(ctx, key, value, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to set value to redis: %w", err)
	}

	return nil

}

func (r *Redis) Get(key string) (interface{}, error) {

	ctx := context.Background()

	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get value from redis: %w", err)
	}

	return val, nil

}

func (r *Redis) Delete(key string) error {

	ctx := context.Background()

	_, err := r.client.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to delete value from redis: %w", err)
	}

	return nil

}

func (r *Redis) Close() error {

	return r.client.Close()

}

func (r *Redis) ping() error {

	ctx := context.Background()

	return r.client.Ping(ctx).Err()

}
