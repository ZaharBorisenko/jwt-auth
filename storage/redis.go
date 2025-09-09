package storage

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(addr string) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	return &RedisClient{client: rdb}
}

func (r *RedisClient) AddToBlackList(ctx context.Context, token string, expiration time.Duration) error {
	return r.client.Set(ctx, "blacklist:"+token, "true", expiration).Err()
}

func (r *RedisClient) IsInBlacklist(ctx context.Context, token string) (bool, error) {
	result, err := r.client.Exists(ctx, "blacklist:"+token).Result()
	return result > 0, err
}

func (r *RedisClient) GetBlackListKeys(ctx context.Context) (map[string]map[string]interface{}, error) {
	keys, err := r.client.Keys(ctx, "blacklist:*").Result()
	if err != nil {
		return nil, err
	}

	result := make(map[string]map[string]interface{})

	for _, key := range keys {
		value, err := r.client.Get(ctx, key).Result()
		if err != nil {
			continue
		}
		ttl, err := r.client.TTL(ctx, key).Result()
		if err != nil {
			continue
		}

		result[key] = map[string]interface{}{
			"value": value,
			"TTL":   ttl.Hours(),
		}

	}

	return result, nil
}

func (r *RedisClient) SaveVerificationCode(ctx context.Context, email string, code string) error {
	key := fmt.Sprintf("verify:%s", email)
	return r.client.Set(ctx, key, code, 15*time.Minute).Err()
}

func (r *RedisClient) GetVerificationCode(ctx context.Context, email string) (string, error) {
	key := fmt.Sprintf("verify:%s", email)
	code, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("verification code not found or expired")
	}
	return code, nil
}

func (r *RedisClient) DeleteVerificationCode(ctx context.Context, email string) error {
	key := fmt.Sprintf("verify:%s", email)
	return r.client.Del(ctx, key).Err()
}
