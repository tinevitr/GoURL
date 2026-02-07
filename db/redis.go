package database

import (
	"context"
	"fmt"
	"log"
	"time"
	"gourl/config"

	"github.com/go-redis/redis/v8"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

type RedisClient struct {
	Client *redis.Client
}

func InitRedis() (*RedisClient, error) {
	cfg := config.LoadConfig()

	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %v", err)
	}

	client := redis.NewClient(opt)

	// Test connection
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	log.Println("Successfully connected to Redis")
	return &RedisClient{Client: client}, nil
}

func (r *RedisClient) Close() error {
	return r.Client.Close()
}

// Helper functions
func (r *RedisClient) SaveURL(slug, originalURL string) error {
	expiration := config.URLExpiration
	
	// Save URL mapping
	err := r.Client.Set(ctx, "gourl:url:"+slug, originalURL, expiration).Err()
	if err != nil {
		return err
	}

	// Save stats
	now := time.Now().Unix()
	stats := map[string]interface{}{
		"original_url":  originalURL,
		"created_at":    now,
		"last_accessed": now,
		"click_count":   0,
	}

	statsKey := "gourl:stats:" + slug
	err = r.Client.HSet(ctx, statsKey, stats).Err()
	if err != nil {
		return err
	}

	// Set expiration for stats (longer than URL)
	return r.Client.Expire(ctx, statsKey, config.StatsTTL).Err()
}

func (r *RedisClient) GetURL(slug string) (string, error) {
	return r.Client.Get(ctx, "gourl:url:"+slug).Result()
}

func (r *RedisClient) IncrementClickCount(slug string) error {
	statsKey := "gourl:stats:" + slug
	
	// Increment counter
	err := r.Client.HIncrBy(ctx, statsKey, "click_count", 1).Err()
	if err != nil {
		return err
	}
	
	// Update last accessed time
	return r.Client.HSet(ctx, statsKey, "last_accessed", time.Now().Unix()).Err()
}

func (r *RedisClient) GetStats(slug string) (map[string]string, error) {
	return r.Client.HGetAll(ctx, "gourl:stats:"+slug).Result()
}

func (r *RedisClient) GetAllStats() ([]map[string]string, error) {
	var cursor uint64
	var stats []map[string]string
	
	for {
		var keys []string
		var err error
		keys, cursor, err = r.Client.Scan(ctx, cursor, "gourl:stats:*", 100).Result()
		if err != nil {
			return nil, err
		}

		for _, key := range keys {
			stat, err := r.Client.HGetAll(ctx, key).Result()
			if err == nil {
				// Extract slug from key
				slug := key[12:] // Remove "gourl:stats:"
				stat["slug"] = slug
				stats = append(stats, stat)
			}
		}

		if cursor == 0 {
			break
		}
	}

	return stats, nil
}

func (r *RedisClient) SlugExists(slug string) (bool, error) {
	exists, err := r.Client.Exists(ctx, "gourl:url:"+slug).Result()
	return exists == 1, err
}
