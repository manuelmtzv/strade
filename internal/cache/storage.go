package cache

import "github.com/go-redis/redis/v8"

type Storage struct {
	client *redis.Client
}

func NewStorage(client *redis.Client) Storage {
	return Storage{
		client: client,
	}
}
