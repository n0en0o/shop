package cache

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/n0en0o/marketplace/internal/basket/domain"
	"github.com/n0en0o/marketplace/internal/basket/domain/repositories"
	"github.com/redis/go-redis/v9"
)

const cacheTTL = 30 * time.Second

type RedisCartRepository struct {
	repo   repositories.CartRepository
	client *redis.Client
}

func NewRedisCartRepository(
	repo repositories.CartRepository,
	client *redis.Client,
) *RedisCartRepository {
	return &RedisCartRepository{
		repo:   repo,
		client: client,
	}
}

func (r *RedisCartRepository) Save(
	ctx context.Context,
	cart domain.ShoppingCart,
) (*domain.ShoppingCart, error) {

	result, err := r.repo.Save(ctx, cart)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(result)
	if err == nil {
		if err := r.client.Set(ctx, result.AccountName, data, cacheTTL).Err(); err != nil {
			log.Printf("cache set failed: %v", err)
		}
	}

	return result, nil
}

func (r *RedisCartRepository) Get(
	ctx context.Context,
	accountName string,
) (*domain.ShoppingCart, error) {

	cached, err := r.client.Get(ctx, accountName).Result()
	if err == nil && cached != "" {
		var cart domain.ShoppingCart
		if err := json.Unmarshal([]byte(cached), &cart); err == nil {
			return &cart, nil
		}
	}

	cart, err := r.repo.Get(ctx, accountName)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(cart); err == nil {
		if err := r.client.Set(ctx, accountName, data, cacheTTL).Err(); err != nil {
			log.Printf("cache set failed: %v", err)
		}
	}

	return cart, nil
}

func (r *RedisCartRepository) Remove(
	ctx context.Context,
	accountName string,
) (bool, error) {
	result, err := r.repo.Remove(ctx, accountName)
	if err != nil {
		return false, err
	}

	_ = r.client.Del(ctx, accountName).Err()

	return result, nil
}
