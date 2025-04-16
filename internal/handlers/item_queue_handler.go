package handlers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/GroceryTrak/GroceryTrakService/internal/clients"
	"github.com/GroceryTrak/GroceryTrakService/internal/dtos"
	"github.com/GroceryTrak/GroceryTrakService/internal/models"
	"github.com/GroceryTrak/GroceryTrakService/internal/repository"
)

type ItemQueueHandler struct {
	queue       repository.ItemQueueRepository
	spoonacular *clients.SpoonacularClient
	itemRepo    repository.ItemRepository
	batchSize   int
	interval    time.Duration
}

func NewItemQueueHandler(
	queue repository.ItemQueueRepository,
	spoonacular *clients.SpoonacularClient,
	itemRepo repository.ItemRepository,
) *ItemQueueHandler {
	return &ItemQueueHandler{
		queue:       queue,
		spoonacular: spoonacular,
		itemRepo:    itemRepo,
		batchSize:   10,
		interval:    1 * time.Minute,
	}
}

func (h *ItemQueueHandler) Start(ctx context.Context) error {
	ticker := time.NewTicker(h.interval)
	defer ticker.Stop()

	log.Println("Starting item queue handler...")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := h.processBatch(ctx); err != nil {
				log.Printf("Error processing batch: %v", err)
			}
		}
	}
}

func (h *ItemQueueHandler) processBatch(ctx context.Context) error {
	credits, err := h.queue.CheckAPICredits(ctx)
	if err != nil {
		return fmt.Errorf("failed to check API credits: %w", err)
	}
	if credits <= 0 {
		return fmt.Errorf("no API credits remaining")
	}

	items, err := h.queue.GetNextBatch(ctx, h.batchSize)
	if err != nil {
		return fmt.Errorf("failed to get batch: %w", err)
	}

	if len(items) == 0 {
		return nil
	}

	for _, item := range items {
		spoonacularItem, err := h.spoonacular.SearchIngredient(ctx, item.Name)
		if err != nil {
			log.Printf("Failed to search for item %s: %v", item.Name, err)
			continue
		}

		nutrients := make([]models.ItemNutrient, len(spoonacularItem.Nutrition.Nutrients))
		for i, n := range spoonacularItem.Nutrition.Nutrients {
			nutrients[i] = models.ItemNutrient{
				Name:                n.Name,
				Amount:              n.Amount,
				Unit:                n.Unit,
				PercentOfDailyNeeds: n.PercentOfDailyNeeds,
			}
		}

		updateReq := dtos.ItemRequest{
			Name:          item.Name,
			Image:         spoonacularItem.Image,
			SpoonacularID: uint(spoonacularItem.ID),
			Nutrients:     make([]dtos.ItemNutrientRequest, len(nutrients)),
		}

		for i, n := range nutrients {
			updateReq.Nutrients[i] = dtos.ItemNutrientRequest{
				Name:                n.Name,
				Amount:              n.Amount,
				Unit:                n.Unit,
				PercentOfDailyNeeds: n.PercentOfDailyNeeds,
			}
		}

		_, err = h.itemRepo.UpdateItem(item.ItemID, updateReq)
		if err != nil {
			log.Printf("Failed to update item %d: %v", item.ItemID, err)
			continue
		}

		if err := h.queue.RemoveItem(ctx, item); err != nil {
			log.Printf("Failed to remove item from queue: %v", err)
			continue
		}

		if err := h.queue.DecrementAPICredits(ctx); err != nil {
			log.Printf("Failed to decrement API credits: %v", err)
		}

		time.Sleep(100 * time.Millisecond)
	}

	return nil
}
