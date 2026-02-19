package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenRepository struct {
	client *redis.Client
}

func NewTokenRepository(addr string) (*TokenRepository, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &TokenRepository{client: client}, nil
}

func (r *TokenRepository) StoreRefreshToken(empID int, token string, expires time.Duration) error {
	return r.client.Set(context.Background(), token, empID, expires).Err()
}

func (r *TokenRepository) GetEmpIDByRefreshToken(token string) (int, error) {
	result := r.client.Get(context.Background(), token)
	if err := result.Err(); err != nil {
		return 0, err
	}
	return result.Int()
}

func (r *TokenRepository) DeleteRefreshToken(token string) error {
	return r.client.Del(context.Background(), token).Err()
}
