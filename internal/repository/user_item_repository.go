package repository

import (
	"github.com/GroceryTrak/GroceryTrakService/internal/dtos"
	"github.com/GroceryTrak/GroceryTrakService/internal/models"
	"gorm.io/gorm"
)

type UserItemRepository interface {
	GetAllUserItems(userID uint) (dtos.UserItemsResponse, error)
	GetUserItem(itemID, userID uint) (dtos.UserItemResponse, error)
	CreateUserItem(req dtos.UserItemRequest, userID uint) (dtos.UserItemResponse, error)
	UpdateUserItem(req dtos.UserItemRequest, itemID, userID uint) (dtos.UserItemResponse, error)
	DeleteUserItem(itemID, userID uint) error
	SearchUserItems(query dtos.UserItemQuery, userID uint) (dtos.UserItemsResponse, error)
	PredictUserItems(items []string, userID uint) ([]models.UserItem, error)
}

type UserItemRepositoryImpl struct {
	db *gorm.DB
}

func NewUserItemRepository(db *gorm.DB) UserItemRepository {
	return &UserItemRepositoryImpl{db: db}
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

func (r *UserItemRepositoryImpl) PredictUserItems(items []string, userID uint) ([]models.UserItem, error) {
	var result []models.UserItem

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
				return nil, err
			}
			existingItem = newItem
		}

		userItem := models.UserItem{
			UserID: userID,
			ItemID: existingItem.ID,
		}
		result = append(result, userItem)

		var existingUserItem models.UserItem
		if err := r.db.Where("user_id = ? AND item_id = ?", userItem.UserID, userItem.ItemID).First(&existingUserItem).Error; err != nil {
			if err := r.db.Create(&userItem).Error; err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}
