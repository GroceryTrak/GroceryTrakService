package repository

import (
	"github.com/GroceryTrak/GroceryTrakService/config"
	"github.com/GroceryTrak/GroceryTrakService/internal/models"
	"github.com/GroceryTrak/GroceryTrakService/internal/templates"
)

func GetItem(id uint) (templates.ItemResponse, error) {
	var item models.Item
	if err := config.DB.First(&item, "id = ?", id).Error; err != nil {
		return templates.ItemResponse{}, err
	}

	return templates.ItemResponse{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
	}, nil
}

func CreateItem(req templates.ItemRequest) (templates.ItemResponse, error) {
	item := models.Item{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := config.DB.Create(&item).Error; err != nil {
		return templates.ItemResponse{}, err
	}

	return templates.ItemResponse{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
	}, nil
}

func UpdateItem(id uint, req templates.ItemRequest) (templates.ItemResponse, error) {
	var item models.Item
	if err := config.DB.First(&item, "id = ?", id).Error; err != nil {
		return templates.ItemResponse{}, err
	}

	item.Name = req.Name
	item.Description = req.Description

	if err := config.DB.Save(&item).Error; err != nil {
		return templates.ItemResponse{}, err
	}

	return templates.ItemResponse{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
	}, nil
}

func DeleteItem(id uint) error {
	return config.DB.Delete(&models.Item{}, "id = ?", id).Error
}

func SearchItems(keyword string) (templates.ItemsResponse, error) {
	var items []models.Item
	searchTerm := "%" + keyword + "%"

	result := config.DB.Where("name LIKE ? OR description LIKE ?", searchTerm, searchTerm).Find(&items)
	if result.Error != nil {
		return templates.ItemsResponse{}, result.Error
	}

	// Convert models.Item to templates.ItemResponse
	var itemResponses []templates.ItemResponse
	for _, item := range items {
		itemResponses = append(itemResponses, templates.ItemResponse{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
		})
	}

	return templates.ItemsResponse{
		Items: itemResponses,
	}, nil
}
