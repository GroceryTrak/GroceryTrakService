package repository

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/GroceryTrak/GroceryTrakService/internal/dtos"
	"github.com/GroceryTrak/GroceryTrakService/internal/models"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

type UserItemRepository interface {
	GetAllUserItems(userID uint) (dtos.UserItemsResponse, error)
	GetUserItem(itemID, userID uint) (dtos.UserItemResponse, error)
	CreateUserItem(req dtos.UserItemRequest, userID uint) (dtos.UserItemResponse, error)
	UpdateUserItem(req dtos.UserItemRequest, itemID, userID uint) (dtos.UserItemResponse, error)
	DeleteUserItem(itemID, userID uint) error
	SearchUserItems(query dtos.UserItemQuery, userID uint) (dtos.UserItemsResponse, error)
	PredictUserItems(items []string, userID uint) (dtos.UserItemsResponse, error)
	DetectUserItems(imageData []byte, userID uint, apiKey string) (dtos.UserItemsResponse, error)
}

type UserItemRepositoryImpl struct {
	db    *gorm.DB
	queue ItemQueueRepository
}

func NewUserItemRepository(db *gorm.DB, queue ItemQueueRepository) UserItemRepository {
	return &UserItemRepositoryImpl{
		db:    db,
		queue: queue,
	}
}

func (r *UserItemRepositoryImpl) GetAllUserItems(userID uint) (dtos.UserItemsResponse, error) {
	var userItems []models.UserItem
	if err := r.db.Where("user_id = ?", userID).Find(&userItems).Error; err != nil {
		return dtos.UserItemsResponse{}, err
	}

	var userItemResponses []dtos.UserItemResponse
	for _, userItem := range userItems {
		var item models.Item
		if err := r.db.First(&item, "id = ?", userItem.ItemID).Error; err != nil {
			return dtos.UserItemsResponse{}, err
		}

		userItemResponses = append(userItemResponses, dtos.UserItemResponse{
			Item: dtos.ItemResponse{
				ID:            item.ID,
				Name:          item.Name,
				Image:         item.Image,
				SpoonacularID: item.SpoonacularID,
			},
			Amount: userItem.Amount,
			Unit:   userItem.Unit,
		})
	}

	return dtos.UserItemsResponse{
		UserItems: userItemResponses,
	}, nil
}

func (r *UserItemRepositoryImpl) GetUserItem(itemID, userID uint) (dtos.UserItemResponse, error) {
	var userItem models.UserItem
	if err := r.db.First(&userItem, "item_id = ? AND user_id = ?", itemID, userID).Error; err != nil {
		return dtos.UserItemResponse{}, err
	}

	var item models.Item
	if err := r.db.First(&item, "id = ?", itemID).Error; err != nil {
		return dtos.UserItemResponse{}, err
	}

	return dtos.UserItemResponse{
		Item: dtos.ItemResponse{
			ID:            item.ID,
			Name:          item.Name,
			Image:         item.Image,
			SpoonacularID: item.SpoonacularID,
		},
		Amount: userItem.Amount,
		Unit:   userItem.Unit,
	}, nil
}

func (r *UserItemRepositoryImpl) CreateUserItem(req dtos.UserItemRequest, userID uint) (dtos.UserItemResponse, error) {
	userItem := models.UserItem{
		UserID: userID,
		ItemID: req.ItemID,
		Amount: req.Amount,
		Unit:   req.Unit,
	}

	if err := r.db.Create(&userItem).Error; err != nil {
		return dtos.UserItemResponse{}, err
	}

	var item models.Item
	if err := r.db.First(&item, "id = ?", userItem.ItemID).Error; err != nil {
		return dtos.UserItemResponse{}, err
	}

	return dtos.UserItemResponse{
		Item: dtos.ItemResponse{
			ID:            item.ID,
			Name:          item.Name,
			Image:         item.Image,
			SpoonacularID: item.SpoonacularID,
		},
		Amount: userItem.Amount,
		Unit:   userItem.Unit,
	}, nil
}

func (r *UserItemRepositoryImpl) UpdateUserItem(req dtos.UserItemRequest, itemID, userID uint) (dtos.UserItemResponse, error) {
	var userItem models.UserItem
	if err := r.db.First(&userItem, "item_id = ? AND user_id = ?", itemID, userID).Error; err != nil {
		return dtos.UserItemResponse{}, err
	}

	userItem.Amount = req.Amount
	userItem.Unit = req.Unit

	if err := r.db.Save(&userItem).Error; err != nil {
		return dtos.UserItemResponse{}, err
	}

	var item models.Item
	if err := r.db.First(&item, "id = ?", userItem.ItemID).Error; err != nil {
		return dtos.UserItemResponse{}, err
	}

	return dtos.UserItemResponse{
		Item: dtos.ItemResponse{
			ID:            item.ID,
			Name:          item.Name,
			Image:         item.Image,
			SpoonacularID: item.SpoonacularID,
		},
		Amount: userItem.Amount,
		Unit:   userItem.Unit,
	}, nil
}

func (r *UserItemRepositoryImpl) DeleteUserItem(itemID, userID uint) error {
	if err := r.db.Delete(&models.UserItem{}, "item_id = ? AND user_id = ?", itemID, userID).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserItemRepositoryImpl) SearchUserItems(query dtos.UserItemQuery, userID uint) (dtos.UserItemsResponse, error) {
	var userItems []models.UserItem
	searchTerm := "%" + query.Name + "%"

	result := r.db.Where("user_id = ?", userID).Find(&userItems)
	if result.Error != nil {
		return dtos.UserItemsResponse{}, result.Error
	}

	var userItemResponses []dtos.UserItemResponse
	for _, userItem := range userItems {
		var item models.Item
		if err := r.db.Where("id = ? AND name LIKE ?", userItem.ItemID, searchTerm).First(&item).Error; err != nil {
			continue
		}

		userItemResponses = append(userItemResponses, dtos.UserItemResponse{
			Item: dtos.ItemResponse{
				ID:            item.ID,
				Name:          item.Name,
				Image:         item.Image,
				SpoonacularID: item.SpoonacularID,
			},
			Amount: userItem.Amount,
			Unit:   userItem.Unit,
		})
	}

	return dtos.UserItemsResponse{
		UserItems: userItemResponses,
	}, nil
}

func (r *UserItemRepositoryImpl) PredictUserItems(items []string, userID uint) (dtos.UserItemsResponse, error) {
	var userItemResponses []dtos.UserItemResponse

	for _, class := range items {
		if class == "" {
			continue
		}

		var existingItem models.Item
		if err := r.db.Where("name LIKE ?", "%"+class+"%").First(&existingItem).Error; err == nil {
		} else {
			newItem := models.Item{
				Name:  class,
				Image: "",
			}
			if err := r.db.Create(&newItem).Error; err != nil {
				return dtos.UserItemsResponse{}, err
			}
			existingItem = newItem
		}

		userItem := models.UserItem{
			UserID: userID,
			ItemID: existingItem.ID,
		}

		var existingUserItem models.UserItem
		if err := r.db.Where("user_id = ? AND item_id = ?", userItem.UserID, userItem.ItemID).First(&existingUserItem).Error; err != nil {
			if err := r.db.Create(&userItem).Error; err != nil {
				return dtos.UserItemsResponse{}, err
			}
		}

		userItemResponses = append(userItemResponses, dtos.UserItemResponse{
			Item: dtos.ItemResponse{
				ID:            existingItem.ID,
				Name:          existingItem.Name,
				Image:         existingItem.Image,
				SpoonacularID: existingItem.SpoonacularID,
			},
			Amount: userItem.Amount,
			Unit:   userItem.Unit,
		})
	}

	return dtos.UserItemsResponse{
		UserItems: userItemResponses,
	}, nil
}

func (r *UserItemRepositoryImpl) DetectUserItems(imageData []byte, userID uint, apiKey string) (dtos.UserItemsResponse, error) {
	client := openai.NewClient(apiKey)
	imageBase64 := base64.StdEncoding.EncodeToString(imageData)
	prompt := `You are a grocery item detector. Analyze the image and identify all grocery items. For each item, provide:
1. The name of the item, uppercase first letter
2. The amount (as a whole number or decimal number with two decimal places)
3. The unit of measurement (e.g., unit, kg, g, lb, oz, etc.)

Return ONLY a JSON array of objects with these exact keys: name, amount, unit.
Example response:
[{"name":"Apple","amount":1,"unit":"unit"},{"name":"Milk","amount":1.5,"unit":"L"}]

Do not include any other text, explanations, or formatting. Return only the JSON array.`

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "gpt-4.1-mini",
			Messages: []openai.ChatCompletionMessage{
				{
					Role: "user",
					MultiContent: []openai.ChatMessagePart{
						{
							Type: "text",
							Text: prompt,
						},
						{
							Type: "image_url",
							ImageURL: &openai.ChatMessageImageURL{
								URL: fmt.Sprintf("data:image/jpeg;base64,%s", imageBase64),
							},
						},
					},
				},
			},
		},
	)
	if err != nil {
		return dtos.UserItemsResponse{}, err
	}

	var detectedItems []struct {
		Name   string  `json:"name"`
		Amount float64 `json:"amount"`
		Unit   string  `json:"unit"`
	}

	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &detectedItems); err != nil {
		return dtos.UserItemsResponse{}, err
	}

	var userItemResponses []dtos.UserItemResponse

	for _, detectedItem := range detectedItems {
		var item models.Item
		if err := r.db.Where("name = ?", detectedItem.Name).First(&item).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				item = models.Item{
					Name:          detectedItem.Name,
					Image:         "",
					SpoonacularID: 0,
				}
				if err := r.db.Create(&item).Error; err != nil {
					return dtos.UserItemsResponse{}, err
				}

				if r.queue != nil {
					queueItem := models.QueueItem{
						ItemID:    item.ID,
						Name:      item.Name,
						CreatedAt: time.Now(),
						Priority:  models.DefaultPriority,
					}
					if err := r.queue.AddItem(context.Background(), queueItem); err != nil {
						log.Printf("Failed to add item to enrichment queue: %v", err)
					}
				}
			} else {
				return dtos.UserItemsResponse{}, err
			}
		}

		userItem := models.UserItem{
			UserID: userID,
			ItemID: item.ID,
			Amount: float32(detectedItem.Amount),
			Unit:   detectedItem.Unit,
		}

		var existingUserItem models.UserItem
		if err := r.db.Where("user_id = ? AND item_id = ?", userID, item.ID).First(&existingUserItem).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := r.db.Create(&userItem).Error; err != nil {
					return dtos.UserItemsResponse{}, err
				}
			} else {
				return dtos.UserItemsResponse{}, err
			}
		} else {
			existingUserItem.Amount = float32(detectedItem.Amount)
			existingUserItem.Unit = detectedItem.Unit
			if err := r.db.Save(&existingUserItem).Error; err != nil {
				return dtos.UserItemsResponse{}, err
			}
			userItem = existingUserItem
		}

		userItemResponses = append(userItemResponses, dtos.UserItemResponse{
			Item: dtos.ItemResponse{
				ID:            item.ID,
				Name:          item.Name,
				Image:         item.Image,
				SpoonacularID: item.SpoonacularID,
			},
			Amount: userItem.Amount,
			Unit:   userItem.Unit,
		})
	}

	return dtos.UserItemsResponse{
		UserItems: userItemResponses,
	}, nil
}
