package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/GroceryTrak/GroceryTrakService/internal/models"
	"github.com/redis/go-redis/v9"
)

const (
	queueKey       = "item_enrichment_queue"
	apiCreditsKey  = "spoonacular_api_credits"
	defaultCredits = 150
)

type ItemQueueRepository interface {
	AddItem(ctx context.Context, item models.QueueItem) error
	GetNextBatch(ctx context.Context, batchSize int) ([]models.QueueItem, error)
	RemoveItem(ctx context.Context, item models.QueueItem) error
	CheckAPICredits(ctx context.Context) (int, error)
	DecrementAPICredits(ctx context.Context) error
}

type ItemQueueRepositoryImpl struct {
	redis *redis.Client
}

func NewItemQueueRepository(redis *redis.Client) ItemQueueRepository {
	return &ItemQueueRepositoryImpl{
		redis: redis,
	}
}

func (r *ItemQueueRepositoryImpl) AddItem(ctx context.Context, item models.QueueItem) error {
	itemBytes, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	score := float64(item.Priority)*1000000 + float64(item.CreatedAt.Unix())
	if err := r.redis.ZAdd(ctx, queueKey, redis.Z{
		Score:  score,
		Member: itemBytes,
	}).Err(); err != nil {
		return fmt.Errorf("failed to add item to queue: %w", err)
	}

	return nil
}

func (r *ItemQueueRepositoryImpl) GetNextBatch(ctx context.Context, batchSize int) ([]models.QueueItem, error) {
	// Get items from queue
	result, err := r.redis.ZRange(ctx, queueKey, 0, int64(batchSize-1)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get items from queue: %w", err)
	}

	var items []models.QueueItem
	for _, itemStr := range result {
		var item models.QueueItem
		if err := json.Unmarshal([]byte(itemStr), &item); err != nil {
			return nil, fmt.Errorf("failed to unmarshal item: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *ItemQueueRepositoryImpl) RemoveItem(ctx context.Context, item models.QueueItem) error {
	itemBytes, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	if err := r.redis.ZRem(ctx, queueKey, itemBytes).Err(); err != nil {
		return fmt.Errorf("failed to remove item from queue: %w", err)
	}

	return nil
}

func (r *ItemQueueRepositoryImpl) CheckAPICredits(ctx context.Context) (int, error) {
	// Get current credits
	credits, err := r.redis.Get(ctx, apiCreditsKey).Int()
	if err == redis.Nil {
		// Initialize credits if not set
		if err := r.redis.Set(ctx, apiCreditsKey, defaultCredits, 24*time.Hour).Err(); err != nil {
			return 0, fmt.Errorf("failed to initialize API credits: %w", err)
		}
		return defaultCredits, nil
	} else if err != nil {
		return 0, fmt.Errorf("failed to get API credits: %w", err)
	}

	return credits, nil
}

func (r *ItemQueueRepositoryImpl) DecrementAPICredits(ctx context.Context) error {
	if err := r.redis.Decr(ctx, apiCreditsKey).Err(); err != nil {
		return fmt.Errorf("failed to decrement API credits: %w", err)
	}
	return nil
}
