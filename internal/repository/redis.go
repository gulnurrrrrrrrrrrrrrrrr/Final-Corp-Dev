package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"quadlingo/internal/config"
	"quadlingo/internal/models"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedis(cfg *config.Config) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	redisClient = rdb

	// Проверяем подключение
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}

	return nil
}

func GetCachedLessons() ([]models.Lesson, error) {
	ctx := context.Background()

	// 1. Проверяем кэш Redis (TTL 10 минут)
	cached, err := redisClient.Get(ctx, "lessons:all").Result()
	if err == nil {
		var lessons []models.Lesson
		if json.Unmarshal([]byte(cached), &lessons) == nil {
			return lessons, nil
		}
	}

	// 2. MISS — получаем из БД
	lessons, err := GetAllLessons()
	if err != nil {
		return nil, err
	}

	// 3. Кэшируем в Redis на 10 минут
	data, _ := json.Marshal(lessons)
	redisClient.Set(ctx, "lessons:all", data, 10*time.Minute)

	return lessons, nil
}
