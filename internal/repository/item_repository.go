package repository

import (
	"github.com/GroceryTrak/GroceryTrakService/config"
	"github.com/GroceryTrak/GroceryTrakService/internal/dtos"
	"github.com/GroceryTrak/GroceryTrakService/internal/models"
)

func GetItem(id uint) (dtos.ItemResponse, error) {
	var item models.Item
	if err := config.DB.First(&item, "id = ?", id).Error; err != nil {
		return dtos.ItemResponse{}, err
	}

	return dtos.ItemResponse{
		ID:            item.ID,
		Name:          item.Name,
		Image:         item.Image,
		SpoonacularID: item.SpoonacularID,
	}, nil
}

func CreateItem(req dtos.ItemRequest) (dtos.ItemResponse, error) {
	item := models.Item{
		Name:          req.Name,
		Image:         req.Image,
		SpoonacularID: req.SpoonacularID,
	}

	if err := config.DB.Create(&item).Error; err != nil {
		return dtos.ItemResponse{}, err
	}

	return dtos.ItemResponse{
		ID:            item.ID,
		Name:          item.Name,
		Image:         item.Image,
		SpoonacularID: item.SpoonacularID,
	}, nil
}

func UpdateItem(id uint, req dtos.ItemRequest) (dtos.ItemResponse, error) {
	var item models.Item
	if err := config.DB.First(&item, "id = ?", id).Error; err != nil {
		return dtos.ItemResponse{}, err
	}

	item.Name = req.Name
	item.Image = req.Image
	item.SpoonacularID = req.SpoonacularID

	if err := config.DB.Save(&item).Error; err != nil {
		return dtos.ItemResponse{}, err
	}

	return dtos.ItemResponse{
		ID:            item.ID,
		Name:          item.Name,
		Image:         item.Image,
		SpoonacularID: item.SpoonacularID,
	}, nil
}

func DeleteItem(id uint) error {
	return config.DB.Delete(&models.Item{}, "id = ?", id).Error
}

func SearchItems(keyword string) (dtos.ItemsResponse, error) {
	var items []models.Item
	searchTerm := "%" + keyword + "%"

	result := config.DB.Where("name LIKE ?", searchTerm).Find(&items)
	if result.Error != nil {
		return dtos.ItemsResponse{}, result.Error
	}

	var itemResponses []dtos.ItemResponse
	for _, item := range items {
		itemResponses = append(itemResponses, dtos.ItemResponse{
			ID:            item.ID,
			Name:          item.Name,
			Image:         item.Image,
			SpoonacularID: item.SpoonacularID,
		})
	}

	return dtos.ItemsResponse{
		Items: itemResponses,
	}, nil
}
