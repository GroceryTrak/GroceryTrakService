package repository

import (
	"strconv"
	"strings"

	"github.com/GroceryTrak/GroceryTrakService/internal/dtos"
	"github.com/GroceryTrak/GroceryTrakService/internal/models"
	"gorm.io/gorm"
)

type RecipeRepository interface {
	GetRecipe(id uint) (*dtos.RecipeResponse, error)
	CreateRecipe(req dtos.RecipeRequest) (*dtos.RecipeResponse, error)
	UpdateRecipe(id uint, req dtos.RecipeRequest) (*dtos.RecipeResponse, error)
	DeleteRecipe(id uint) error
	SearchRecipes(query dtos.RecipeQuery) (dtos.RecipesResponse, error)
}

type RecipeRepositoryImpl struct {
	db *gorm.DB
}

func NewRecipeRepository(db *gorm.DB) RecipeRepository {
	return &RecipeRepositoryImpl{db: db}
}

func (r *RecipeRepositoryImpl) GetRecipe(id uint) (*dtos.RecipeResponse, error) {
	var recipe models.Recipe
	if err := r.db.Preload("Ingredients.Item").First(&recipe, id).Error; err != nil {
		return nil, err
	}
	ingredients := make([]dtos.RecipeItemResponse, len(recipe.Ingredients))
	for i, item := range recipe.Ingredients {
		ingredients[i] = dtos.RecipeItemResponse{
			Item: dtos.ItemResponse{
				ID:            item.Item.ID,
				Name:          item.Item.Name,
				Image:         item.Item.Image,
				SpoonacularID: item.Item.SpoonacularID,
			},
			Amount: item.Amount,
			Unit:   item.Unit,
		}
	}

	return &dtos.RecipeResponse{
		ID:          recipe.ID,
		Title:       recipe.Title,
		ReadyTime:   recipe.ReadyTime,
		CookingTime: recipe.CookingTime,
		PrepTime:    recipe.PrepTime,
		Image:       recipe.Image,
		KCal:        recipe.KCal,
		Vegan:       recipe.Vegan,
		Vegetarian:  recipe.Vegetarian,
		Ingredients: ingredients,
	}, nil
}

func (r *RecipeRepositoryImpl) CreateRecipe(req dtos.RecipeRequest) (*dtos.RecipeResponse, error) {
	recipe := models.Recipe{
		Title:       req.Title,
		ReadyTime:   req.ReadyTime,
		CookingTime: req.CookingTime,
		PrepTime:    req.PrepTime,
		Image:       req.Image,
		KCal:        req.KCal,
		Vegan:       req.Vegan,
		Vegetarian:  req.Vegetarian,
		Ingredients: make([]models.RecipeItem, len(req.Ingredients)),
	}

	for i, item := range req.Ingredients {
		recipe.Ingredients[i] = models.RecipeItem{
			ItemID: item.ItemID,
			Amount: item.Amount,
			Unit:   item.Unit,
		}
	}

	if err := r.db.Create(&recipe).Error; err != nil {
		return nil, err
	}

	if err := r.db.Preload("Ingredients.Item").First(&recipe, recipe.ID).Error; err != nil {
		return nil, err
	}

	ingredients := make([]dtos.RecipeItemResponse, len(recipe.Ingredients))
	for i, item := range recipe.Ingredients {
		ingredients[i] = dtos.RecipeItemResponse{
			Item: dtos.ItemResponse{
				ID:            item.Item.ID,
				Name:          item.Item.Name,
				Image:         item.Item.Image,
				SpoonacularID: item.Item.SpoonacularID,
			},
			Amount: item.Amount,
			Unit:   item.Unit,
		}
	}

	return &dtos.RecipeResponse{
		ID:          recipe.ID,
		Title:       recipe.Title,
		ReadyTime:   recipe.ReadyTime,
		CookingTime: recipe.CookingTime,
		PrepTime:    recipe.PrepTime,
		Image:       recipe.Image,
		KCal:        recipe.KCal,
		Vegan:       recipe.Vegan,
		Vegetarian:  recipe.Vegetarian,
		Ingredients: ingredients,
	}, nil
}

func (r *RecipeRepositoryImpl) UpdateRecipe(id uint, req dtos.RecipeRequest) (*dtos.RecipeResponse, error) {
	var recipe models.Recipe
	if err := r.db.Preload("Ingredients.Item").First(&recipe, id).Error; err != nil {
		return nil, err
	}

	recipe.Title = req.Title
	recipe.ReadyTime = req.ReadyTime
	recipe.CookingTime = req.CookingTime
	recipe.PrepTime = req.PrepTime
	recipe.Image = req.Image
	recipe.KCal = req.KCal
	recipe.Vegan = req.Vegan
	recipe.Vegetarian = req.Vegetarian

	if err := r.db.Save(&recipe).Error; err != nil {
		return nil, err
	}

	ingredients := make([]dtos.RecipeItemResponse, len(recipe.Ingredients))
	for i, item := range recipe.Ingredients {
		ingredients[i] = dtos.RecipeItemResponse{
			Item: dtos.ItemResponse{
				ID:            item.Item.ID,
				Name:          item.Item.Name,
				Image:         item.Item.Image,
				SpoonacularID: item.Item.SpoonacularID,
			},
			Amount: item.Amount,
			Unit:   item.Unit,
		}
	}

	return &dtos.RecipeResponse{
		ID:          recipe.ID,
		Title:       recipe.Title,
		ReadyTime:   recipe.ReadyTime,
		CookingTime: recipe.CookingTime,
		PrepTime:    recipe.PrepTime,
		Image:       recipe.Image,
		KCal:        recipe.KCal,
		Vegan:       recipe.Vegan,
		Vegetarian:  recipe.Vegetarian,
		Ingredients: ingredients,
	}, nil
}

func (r *RecipeRepositoryImpl) DeleteRecipe(id uint) error {
	if err := r.db.Delete(&models.Recipe{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *RecipeRepositoryImpl) SearchRecipes(query dtos.RecipeQuery) (dtos.RecipesResponse, error) {
	var recipes []models.Recipe
	db := r.db

	if query.Title != "" {
		db = db.Where("LOWER(title) LIKE ?", "%"+strings.ToLower(query.Title)+"%")
	}

	validDiets := map[string]string{"vegan": "vegan", "vegetarian": "vegetarian"}
	if dietField, exists := validDiets[strings.ToLower(query.Diet)]; exists {
		db = db.Where(dietField+" = ?", true)
	}

	var ingredientIDs []uint
	for _, id := range query.Ingredients {
		if num, err := strconv.ParseUint(id, 10, 32); err == nil {
			ingredientIDs = append(ingredientIDs, uint(num))
		}
	}

	if len(ingredientIDs) > 0 {
		db = db.Joins("JOIN recipe_items ri ON recipes.id = ri.recipe_id").
			Where("ri.item_id IN ?", ingredientIDs).
			Distinct("recipes.*") // Prevent duplicates
	}

	var dietCounts []dtos.DietCount
	if err := db.Model(&models.Recipe{}).
		Select("vegan, vegetarian, COUNT(*) as count").
		Group("vegan, vegetarian").
		Scan(&dietCounts).Error; err != nil {
		return dtos.RecipesResponse{}, err
	}

	var totalCount int64
	for _, dc := range dietCounts {
		totalCount += dc.Count
	}

	if err := db.Preload("Ingredients.Item").Find(&recipes).Error; err != nil {
		return dtos.RecipesResponse{}, err
	}

	recipeResponses := make([]dtos.RecipeResponse, len(recipes))
	for i, recipe := range recipes {
		ingredients := make([]dtos.RecipeItemResponse, len(recipe.Ingredients))
		for j, item := range recipe.Ingredients {
			ingredients[j] = dtos.RecipeItemResponse{
				Item: dtos.ItemResponse{
					ID:            item.Item.ID,
					Name:          item.Item.Name,
					Image:         item.Item.Image,
					SpoonacularID: item.Item.SpoonacularID,
				},
				Amount: item.Amount,
				Unit:   item.Unit,
			}
		}
		recipeResponses[i] = dtos.RecipeResponse{
			ID:          recipe.ID,
			Title:       recipe.Title,
			ReadyTime:   recipe.ReadyTime,
			CookingTime: recipe.CookingTime,
			PrepTime:    recipe.PrepTime,
			Image:       recipe.Image,
			KCal:        recipe.KCal,
			Vegan:       recipe.Vegan,
			Vegetarian:  recipe.Vegetarian,
			Ingredients: ingredients,
		}
	}

	return dtos.RecipesResponse{
		Recipes:    recipeResponses,
		Count:      int(totalCount),
		DietCounts: dietCounts,
	}, nil
}
