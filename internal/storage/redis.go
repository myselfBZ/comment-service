package storage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisStore struct {
	client *redis.Client
}

var NotFound = redis.Nil

const COMMENTKEY = "comment-%s"

func InitRedisClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		DB:       0,
		Password: "",
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return client, err
}

func NewRedistStore(c *redis.Client) *RedisStore {
	return &RedisStore{c}
}

func (r *RedisStore) GetByID(id string, ctx context.Context) (*Comment, error) {
	result, err := r.client.Get(ctx, fmt.Sprintf(COMMENTKEY, id)).Result()
	if err != nil {
		return nil, err
	}
	var comment Comment
	if err := json.Unmarshal([]byte(result), &comment); err != nil {
		return nil, err
	}
	return &comment, err
}

func (r *RedisStore) Create(c *Comment, ctx context.Context) error {
	encoded, err := json.Marshal(c)
	if err != nil {
		return err
	}
	_, err = r.client.Set(ctx, fmt.Sprintf(COMMENTKEY, c.ID), encoded, 0).Result()
	return err
}
