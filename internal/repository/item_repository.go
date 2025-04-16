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
	if err := r.db.Preload("Nutrients").First(&item, "id = ?", id).Error; err != nil {
		return dtos.ItemResponse{}, err
	}

	nutrients := make([]dtos.ItemNutrientResponse, len(item.Nutrients))
	for i, n := range item.Nutrients {
		nutrients[i] = dtos.ItemNutrientResponse{
			Name:                n.Name,
			Amount:              n.Amount,
			Unit:                n.Unit,
			PercentOfDailyNeeds: n.PercentOfDailyNeeds,
		}
	}

	return dtos.ItemResponse{
		ID:            item.ID,
		Name:          item.Name,
		Image:         item.Image,
		SpoonacularID: item.SpoonacularID,
		Nutrients:     nutrients,
	}, nil
}

func (r *ItemRepositoryImpl) CreateItem(req dtos.ItemRequest) (dtos.ItemResponse, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return dtos.ItemResponse{}, tx.Error
	}

	item := models.Item{
		Name:          req.Name,
		Image:         req.Image,
		SpoonacularID: req.SpoonacularID,
		Nutrients:     make([]models.ItemNutrient, len(req.Nutrients)),
	}

	for i, n := range req.Nutrients {
		item.Nutrients[i] = models.ItemNutrient{
			Name:                n.Name,
			Amount:              n.Amount,
			Unit:                n.Unit,
			PercentOfDailyNeeds: n.PercentOfDailyNeeds,
		}
	}

	if err := tx.Create(&item).Error; err != nil {
		tx.Rollback()
		return dtos.ItemResponse{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return dtos.ItemResponse{}, err
	}

	nutrientResponses := make([]dtos.ItemNutrientResponse, len(item.Nutrients))
	for i, n := range item.Nutrients {
		nutrientResponses[i] = dtos.ItemNutrientResponse{
			Name:                n.Name,
			Amount:              n.Amount,
			Unit:                n.Unit,
			PercentOfDailyNeeds: n.PercentOfDailyNeeds,
		}
	}

	return dtos.ItemResponse{
		ID:            item.ID,
		Name:          item.Name,
		Image:         item.Image,
		SpoonacularID: item.SpoonacularID,
		Nutrients:     nutrientResponses,
	}, nil
}

func (r *ItemRepositoryImpl) UpdateItem(id uint, req dtos.ItemRequest) (dtos.ItemResponse, error) {
	var item models.Item
	if err := r.db.First(&item, "id = ?", id).Error; err != nil {
		return dtos.ItemResponse{}, err
	}

	if err := r.db.Where("item_id = ?", id).Delete(&models.ItemNutrient{}).Error; err != nil {
		return dtos.ItemResponse{}, err
	}

	modelNutrients := make([]models.ItemNutrient, len(req.Nutrients))
	for i, n := range req.Nutrients {
		modelNutrients[i] = models.ItemNutrient{
			ItemID:              id,
			Name:                n.Name,
			Amount:              n.Amount,
			Unit:                n.Unit,
			PercentOfDailyNeeds: n.PercentOfDailyNeeds,
		}
	}

	item.Name = req.Name
	item.Image = req.Image
	item.SpoonacularID = req.SpoonacularID

	tx := r.db.Begin()
	if tx.Error != nil {
		return dtos.ItemResponse{}, tx.Error
	}

	if err := tx.Save(&item).Error; err != nil {
		tx.Rollback()
		return dtos.ItemResponse{}, err
	}

	if err := tx.Create(&modelNutrients).Error; err != nil {
		tx.Rollback()
		return dtos.ItemResponse{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return dtos.ItemResponse{}, err
	}

	if err := r.db.Preload("Nutrients").First(&item, id).Error; err != nil {
		return dtos.ItemResponse{}, err
	}

	nutrients := make([]dtos.ItemNutrientResponse, len(item.Nutrients))
	for i, n := range item.Nutrients {
		nutrients[i] = dtos.ItemNutrientResponse{
			Name:                n.Name,
			Amount:              n.Amount,
			Unit:                n.Unit,
			PercentOfDailyNeeds: n.PercentOfDailyNeeds,
		}
	}

	return dtos.ItemResponse{
		ID:            item.ID,
		Name:          item.Name,
		Image:         item.Image,
		SpoonacularID: item.SpoonacularID,
		Nutrients:     nutrients,
	}, nil
}

func (r *ItemRepositoryImpl) DeleteItem(id uint) error {
	if err := r.db.Where("item_id = ?", id).Delete(&models.ItemNutrient{}).Error; err != nil {
		return err
	}
	return r.db.Delete(&models.Item{}, "id = ?", id).Error
}

func (r *ItemRepositoryImpl) SearchItems(keyword string) (dtos.ItemsResponse, error) {
	var items []models.Item
	searchTerm := "%" + keyword + "%"

	result := r.db.Preload("Nutrients").Where("name LIKE ?", searchTerm).Find(&items)
	if result.Error != nil {
		return dtos.ItemsResponse{}, result.Error
	}

	var itemResponses []dtos.ItemResponse
	for _, item := range items {
		nutrients := make([]dtos.ItemNutrientResponse, len(item.Nutrients))
		for i, n := range item.Nutrients {
			nutrients[i] = dtos.ItemNutrientResponse{
				Name:                n.Name,
				Amount:              n.Amount,
				Unit:                n.Unit,
				PercentOfDailyNeeds: n.PercentOfDailyNeeds,
			}
		}

		itemResponses = append(itemResponses, dtos.ItemResponse{
			ID:            item.ID,
			Name:          item.Name,
			Image:         item.Image,
			SpoonacularID: item.SpoonacularID,
			Nutrients:     nutrients,
		})
	}

	return dtos.ItemsResponse{
		Items: itemResponses,
	}, nil
}
