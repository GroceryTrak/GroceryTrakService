package repository

import (
	"github.com/GroceryTrak/GroceryTrakService/internal/dtos"
	"github.com/GroceryTrak/GroceryTrakService/internal/models"
	"gorm.io/gorm"
)

type ItemRepositoryImpl struct {
	db *gorm.DB
}

type ItemRepository interface {
	GetItem(id uint) (dtos.ItemResponse, error)
	CreateItem(req dtos.ItemRequest) (dtos.ItemResponse, error)
	UpdateItem(id uint, req dtos.ItemRequest) (dtos.ItemResponse, error)
	DeleteItem(id uint) error
	SearchItems(keyword string) (dtos.ItemsResponse, error)
}

func NewItemRepository(db *gorm.DB) ItemRepository {
	return &ItemRepositoryImpl{db: db}
}

func (r *ItemRepositoryImpl) GetItem(id uint) (dtos.ItemResponse, error) {
	var item models.Item
	if err := r.db.First(&item, "id = ?", id).Error; err != nil {
		return dtos.ItemResponse{}, err
	}

	return dtos.ItemResponse{
		ID:            item.ID,
		Name:          item.Name,
		Image:         item.Image,
		SpoonacularID: item.SpoonacularID,
	}, nil
}

func (r *ItemRepositoryImpl) CreateItem(req dtos.ItemRequest) (dtos.ItemResponse, error) {
	item := models.Item{
		Name:          req.Name,
		Image:         req.Image,
		SpoonacularID: req.SpoonacularID,
	}

	if err := r.db.Create(&item).Error; err != nil {
		return dtos.ItemResponse{}, err
	}

	return dtos.ItemResponse{
		ID:            item.ID,
		Name:          item.Name,
		Image:         item.Image,
		SpoonacularID: item.SpoonacularID,
	}, nil
}

func (r *ItemRepositoryImpl) UpdateItem(id uint, req dtos.ItemRequest) (dtos.ItemResponse, error) {
	var item models.Item
	if err := r.db.First(&item, "id = ?", id).Error; err != nil {
		return dtos.ItemResponse{}, err
	}

	item.Name = req.Name
	item.Image = req.Image
	item.SpoonacularID = req.SpoonacularID

	if err := r.db.Save(&item).Error; err != nil {
		return dtos.ItemResponse{}, err
	}

	return dtos.ItemResponse{
		ID:            item.ID,
		Name:          item.Name,
		Image:         item.Image,
		SpoonacularID: item.SpoonacularID,
	}, nil
}

func (r *ItemRepositoryImpl) DeleteItem(id uint) error {
	return r.db.Delete(&models.Item{}, "id = ?", id).Error
}

func (r *ItemRepositoryImpl) SearchItems(keyword string) (dtos.ItemsResponse, error) {
	var items []models.Item
	searchTerm := "%" + keyword + "%"

	result := r.db.Where("name LIKE ?", searchTerm).Find(&items)
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
