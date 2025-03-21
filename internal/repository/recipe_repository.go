package repository

import (
	"strconv"
	"strings"

	"github.com/GroceryTrak/GroceryTrakService/config"
	"github.com/GroceryTrak/GroceryTrakService/internal/dtos"
	"github.com/GroceryTrak/GroceryTrakService/internal/models"
)

func GetRecipe(id uint) (*dtos.RecipeResponse, error) {
	var recipe models.Recipe
	if err := config.DB.Preload("Ingredients.Item").First(&recipe, id).Error; err != nil {
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

func CreateRecipe(req dtos.RecipeRequest) (*dtos.RecipeResponse, error) {
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

	if err := config.DB.Create(&recipe).Error; err != nil {
		return nil, err
	}

	if err := config.DB.Preload("Ingredients.Item").First(&recipe, recipe.ID).Error; err != nil {
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

func UpdateRecipe(id uint, req dtos.RecipeRequest) (*dtos.RecipeResponse, error) {
	var recipe models.Recipe
	if err := config.DB.Preload("Ingredients.Item").First(&recipe, id).Error; err != nil {
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

	if err := config.DB.Save(&recipe).Error; err != nil {
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

func DeleteRecipe(id uint) error {
	if err := config.DB.Delete(&models.Recipe{}, id).Error; err != nil {
		return err
	}
	return nil
}

func SearchRecipes(query dtos.RecipeQuery) (dtos.RecipesResponse, error) {
	var recipes []models.Recipe

	if query.Title != "" {
		config.DB = config.DB.Where("LOWER(title) LIKE ?", "%"+strings.ToLower(query.Title)+"%")
	}

	// Handle diet filtering
	validDiets := map[string]string{"vegan": "vegan", "vegetarian": "vegetarian"}
	if dietField, exists := validDiets[strings.ToLower(query.Diet)]; exists {
		config.DB = config.DB.Where(dietField+" = ?", true)
	}

	// Handle ingredient filtering safely
	var ingredientIDs []uint
	for _, id := range query.Ingredients {
		if num, err := strconv.ParseUint(id, 10, 32); err == nil {
			ingredientIDs = append(ingredientIDs, uint(num))
		}
	}

	if len(ingredientIDs) > 0 {
		config.DB = config.DB.Joins("JOIN recipe_items ri ON recipes.id = ri.recipe_id").
			Where("ri.item_id IN ?", ingredientIDs).
			Distinct("recipes.*") // Prevent duplicates
	}

	if err := config.DB.Preload("Ingredients.Item").Find(&recipes).Error; err != nil {
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

	return dtos.RecipesResponse{Recipes: recipeResponses}, nil
}
