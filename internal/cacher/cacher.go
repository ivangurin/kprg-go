package cacher

import (
	"fmt"
	"kprg/internal/cacher/redis"
	"kprg/internal/models"
)

type Cacher interface {
	Has(key string) (bool, error)
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Delete(key string) error
	Close() error
}

func NewCacher(connection *models.Connection) (Cacher, error) {

	var cacher Cacher
	var err error

	cacher, err = redis.NewRedis(connection)
	if err != nil {
		return nil, fmt.Errorf("failed to create cacher: %w", err)
	}

	return cacher, nil

}
