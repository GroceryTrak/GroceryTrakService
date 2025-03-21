package repository

import (
	"github.com/GroceryTrak/GroceryTrakService/config"
	"github.com/GroceryTrak/GroceryTrakService/internal/dtos"
	"github.com/GroceryTrak/GroceryTrakService/internal/models"
)

func GetAllUserItems(userID uint) (dtos.UserItemsResponse, error) {
	var userItems []models.UserItem
	if err := config.DB.Where("user_id = ?", userID).Find(&userItems).Error; err != nil {
		return dtos.UserItemsResponse{}, err
	}

	var userItemResponses []dtos.UserItemResponse
	for _, userItem := range userItems {
		var item models.Item
		if err := config.DB.First(&item, "id = ?", userItem.ItemID).Error; err != nil {
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

func GetUserItem(itemID, userID uint) (dtos.UserItemResponse, error) {
	var userItem models.UserItem
	if err := config.DB.First(&userItem, "item_id = ? AND user_id = ?", itemID, userID).Error; err != nil {
		return dtos.UserItemResponse{}, err
	}

	var item models.Item
	if err := config.DB.First(&item, "id = ?", itemID).Error; err != nil {
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

func CreateUserItem(req dtos.UserItemRequest, userID uint) (dtos.UserItemResponse, error) {
	userItem := models.UserItem{
		UserID: userID,
		ItemID: req.ItemID,
		Amount: req.Amount,
		Unit:   req.Unit,
	}

	if err := config.DB.Create(&userItem).Error; err != nil {
		return dtos.UserItemResponse{}, err
	}

	var item models.Item
	if err := config.DB.First(&item, "id = ?", userItem.ItemID).Error; err != nil {
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

func UpdateUserItem(req dtos.UserItemRequest, itemID, userID uint) (dtos.UserItemResponse, error) {
	var userItem models.UserItem
	if err := config.DB.First(&userItem, "item_id = ? AND user_id = ?", itemID, userID).Error; err != nil {
		return dtos.UserItemResponse{}, err
	}

	userItem.Amount = req.Amount
	userItem.Unit = req.Unit

	if err := config.DB.Save(&userItem).Error; err != nil {
		return dtos.UserItemResponse{}, err
	}

	var item models.Item
	if err := config.DB.First(&item, "id = ?", userItem.ItemID).Error; err != nil {
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

func DeleteUserItem(itemID, userID uint) error {
	if err := config.DB.Delete(&models.UserItem{}, "item_id = ? AND user_id = ?", itemID, userID).Error; err != nil {
		return err
	}
	return nil
}

func SearchUserItems(query dtos.UserItemQuery, userID uint) (dtos.UserItemsResponse, error) {
	var userItems []models.UserItem
	searchTerm := "%" + query.Name + "%"

	result := config.DB.Where("user_id = ?", userID).Find(&userItems)
	if result.Error != nil {
		return dtos.UserItemsResponse{}, result.Error
	}

	var userItemResponses []dtos.UserItemResponse
	for _, userItem := range userItems {
		var item models.Item
		if err := config.DB.Where("id = ? AND name LIKE ?", userItem.ItemID, searchTerm).First(&item).Error; err != nil {
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
