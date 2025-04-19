package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/GroceryTrak/GroceryTrakService/internal/clients"
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
	db          *gorm.DB
	spoonacular *clients.SpoonacularClient
	itemQueue   ItemQueueRepository
}

func NewRecipeRepository(db *gorm.DB, spoonacular *clients.SpoonacularClient, itemQueue ItemQueueRepository) RecipeRepository {
	return &RecipeRepositoryImpl{
		db:          db,
		spoonacular: spoonacular,
		itemQueue:   itemQueue,
	}
}

func (r *RecipeRepositoryImpl) GetRecipe(id uint) (*dtos.RecipeResponse, error) {
	var recipe models.Recipe
	if err := r.db.Preload("Ingredients.Item").Preload("Nutrients").Preload("Instructions").First(&recipe, id).Error; err != nil {
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

	nutrients := make([]dtos.RecipeNutrientResponse, len(recipe.Nutrients))
	for i, n := range recipe.Nutrients {
		nutrients[i] = dtos.RecipeNutrientResponse{
			Name:                n.Name,
			Amount:              n.Amount,
			Unit:                n.Unit,
			PercentOfDailyNeeds: n.PercentOfDailyNeeds,
		}
	}

	instructions := make([]dtos.RecipeInstructionResponse, len(recipe.Instructions))
	for i, inst := range recipe.Instructions {
		instructions[i] = dtos.RecipeInstructionResponse{
			Number: inst.Number,
			Step:   inst.Step,
		}
	}

	return &dtos.RecipeResponse{
		ID:            recipe.ID,
		Title:         recipe.Title,
		Summary:       recipe.Summary,
		SpoonacularID: recipe.SpoonacularID,
		Instructions:  instructions,
		Servings:      recipe.Servings,
		ReadyTime:     recipe.ReadyTime,
		CookingTime:   recipe.CookingTime,
		PrepTime:      recipe.PrepTime,
		Image:         recipe.Image,
		KCal:          recipe.KCal,
		Vegan:         recipe.Vegan,
		Vegetarian:    recipe.Vegetarian,
		Ingredients:   ingredients,
		Nutrients:     nutrients,
	}, nil
}

func (r *RecipeRepositoryImpl) CreateRecipe(req dtos.RecipeRequest) (*dtos.RecipeResponse, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	recipe := models.Recipe{
		Title:         req.Title,
		Summary:       req.Summary,
		SpoonacularID: req.SpoonacularID,
		Servings:      req.Servings,
		ReadyTime:     req.ReadyTime,
		CookingTime:   req.CookingTime,
		PrepTime:      req.PrepTime,
		Image:         req.Image,
		KCal:          req.KCal,
		Vegan:         req.Vegan,
		Vegetarian:    req.Vegetarian,
	}

	if err := tx.Create(&recipe).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create ingredients
	ingredients := make([]models.RecipeItem, len(req.Ingredients))
	for i, item := range req.Ingredients {
		ingredients[i] = models.RecipeItem{
			RecipeID: recipe.ID,
			ItemID:   item.ItemID,
			Amount:   item.Amount,
			Unit:     item.Unit,
		}
	}

	if err := tx.Create(&ingredients).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create nutrients
	nutrients := make([]models.RecipeNutrient, len(req.Nutrients))
	for i, n := range req.Nutrients {
		nutrients[i] = models.RecipeNutrient{
			RecipeID:            recipe.ID,
			Name:                n.Name,
			Amount:              n.Amount,
			Unit:                n.Unit,
			PercentOfDailyNeeds: n.PercentOfDailyNeeds,
		}
	}

	if err := tx.Create(&nutrients).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create instructions
	instructions := make([]models.RecipeInstruction, len(req.Instructions))
	for i, inst := range req.Instructions {
		instructions[i] = models.RecipeInstruction{
			RecipeID: recipe.ID,
			Number:   inst.Number,
			Step:     inst.Step,
		}
	}

	if err := tx.Create(&instructions).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	if err := r.db.Preload("Ingredients.Item").Preload("Nutrients").Preload("Instructions").First(&recipe, recipe.ID).Error; err != nil {
		return nil, err
	}

	ingredientResponses := make([]dtos.RecipeItemResponse, len(recipe.Ingredients))
	for i, item := range recipe.Ingredients {
		ingredientResponses[i] = dtos.RecipeItemResponse{
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

	nutrientResponses := make([]dtos.RecipeNutrientResponse, len(recipe.Nutrients))
	for i, n := range recipe.Nutrients {
		nutrientResponses[i] = dtos.RecipeNutrientResponse{
			Name:                n.Name,
			Amount:              n.Amount,
			Unit:                n.Unit,
			PercentOfDailyNeeds: n.PercentOfDailyNeeds,
		}
	}

	instructionResponses := make([]dtos.RecipeInstructionResponse, len(recipe.Instructions))
	for i, inst := range recipe.Instructions {
		instructionResponses[i] = dtos.RecipeInstructionResponse{
			Number: inst.Number,
			Step:   inst.Step,
		}
	}

	return &dtos.RecipeResponse{
		ID:            recipe.ID,
		Title:         recipe.Title,
		Summary:       recipe.Summary,
		SpoonacularID: recipe.SpoonacularID,
		Instructions:  instructionResponses,
		Servings:      recipe.Servings,
		ReadyTime:     recipe.ReadyTime,
		CookingTime:   recipe.CookingTime,
		PrepTime:      recipe.PrepTime,
		Image:         recipe.Image,
		KCal:          recipe.KCal,
		Vegan:         recipe.Vegan,
		Vegetarian:    recipe.Vegetarian,
		Ingredients:   ingredientResponses,
		Nutrients:     nutrientResponses,
	}, nil
}

func (r *RecipeRepositoryImpl) UpdateRecipe(id uint, req dtos.RecipeRequest) (*dtos.RecipeResponse, error) {
	var recipe models.Recipe
	if err := r.db.Preload("Ingredients.Item").Preload("Nutrients").Preload("Instructions").First(&recipe, id).Error; err != nil {
		return nil, err
	}

	recipe.Title = req.Title
	recipe.Summary = req.Summary
	recipe.SpoonacularID = req.SpoonacularID
	recipe.Servings = req.Servings
	recipe.ReadyTime = req.ReadyTime
	recipe.CookingTime = req.CookingTime
	recipe.PrepTime = req.PrepTime
	recipe.Image = req.Image
	recipe.KCal = req.KCal
	recipe.Vegan = req.Vegan
	recipe.Vegetarian = req.Vegetarian

	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Delete existing related records
	if err := tx.Where("recipe_id = ?", id).Delete(&models.RecipeNutrient{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Where("recipe_id = ?", id).Delete(&models.RecipeInstruction{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Where("recipe_id = ?", id).Delete(&models.RecipeItem{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create new nutrients
	nutrients := make([]models.RecipeNutrient, len(req.Nutrients))
	for i, n := range req.Nutrients {
		nutrients[i] = models.RecipeNutrient{
			RecipeID:            id,
			Name:                n.Name,
			Amount:              n.Amount,
			Unit:                n.Unit,
			PercentOfDailyNeeds: n.PercentOfDailyNeeds,
		}
	}

	if err := tx.Create(&nutrients).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create new instructions
	instructions := make([]models.RecipeInstruction, len(req.Instructions))
	for i, inst := range req.Instructions {
		instructions[i] = models.RecipeInstruction{
			RecipeID: id,
			Number:   inst.Number,
			Step:     inst.Step,
		}
	}

	if err := tx.Create(&instructions).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create new ingredients
	ingredients := make([]models.RecipeItem, len(req.Ingredients))
	for i, item := range req.Ingredients {
		ingredients[i] = models.RecipeItem{
			RecipeID: id,
			ItemID:   item.ItemID,
			Amount:   item.Amount,
			Unit:     item.Unit,
		}
	}

	if err := tx.Create(&ingredients).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Save recipe changes
	if err := tx.Save(&recipe).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	if err := r.db.Preload("Ingredients.Item").Preload("Nutrients").Preload("Instructions").First(&recipe, id).Error; err != nil {
		return nil, err
	}

	ingredientResponses := make([]dtos.RecipeItemResponse, len(recipe.Ingredients))
	for i, item := range recipe.Ingredients {
		ingredientResponses[i] = dtos.RecipeItemResponse{
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

	nutrientResponses := make([]dtos.RecipeNutrientResponse, len(recipe.Nutrients))
	for i, n := range recipe.Nutrients {
		nutrientResponses[i] = dtos.RecipeNutrientResponse{
			Name:                n.Name,
			Amount:              n.Amount,
			Unit:                n.Unit,
			PercentOfDailyNeeds: n.PercentOfDailyNeeds,
		}
	}

	instructionResponses := make([]dtos.RecipeInstructionResponse, len(recipe.Instructions))
	for i, inst := range recipe.Instructions {
		instructionResponses[i] = dtos.RecipeInstructionResponse{
			Number: inst.Number,
			Step:   inst.Step,
		}
	}

	return &dtos.RecipeResponse{
		ID:            recipe.ID,
		Title:         recipe.Title,
		Summary:       recipe.Summary,
		SpoonacularID: recipe.SpoonacularID,
		Instructions:  instructionResponses,
		Servings:      recipe.Servings,
		ReadyTime:     recipe.ReadyTime,
		CookingTime:   recipe.CookingTime,
		PrepTime:      recipe.PrepTime,
		Image:         recipe.Image,
		KCal:          recipe.KCal,
		Vegan:         recipe.Vegan,
		Vegetarian:    recipe.Vegetarian,
		Ingredients:   ingredientResponses,
		Nutrients:     nutrientResponses,
	}, nil
}

func (r *RecipeRepositoryImpl) DeleteRecipe(id uint) error {
	if err := r.db.Where("recipe_id = ?", id).Delete(&models.RecipeNutrient{}).Error; err != nil {
		return err
	}

	if err := r.db.Where("recipe_id = ?", id).Delete(&models.RecipeInstruction{}).Error; err != nil {
		return err
	}

	if err := r.db.Where("recipe_id = ?", id).Delete(&models.RecipeItem{}).Error; err != nil {
		return err
	}

	return r.db.Delete(&models.Recipe{}, "id = ?", id).Error
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
	countQuery := r.db.Model(&models.Recipe{})
	if len(ingredientIDs) > 0 {
		countQuery = countQuery.Joins("JOIN recipe_items ri ON recipes.id = ri.recipe_id").
			Where("ri.item_id IN ?", ingredientIDs)
	}
	if err := countQuery.
		Select("vegan, vegetarian, COUNT(*) as count").
		Group("vegan, vegetarian").
		Scan(&dietCounts).Error; err != nil {
		return dtos.RecipesResponse{}, err
	}

	var totalCount int64
	for _, dc := range dietCounts {
		totalCount += dc.Count
	}

	if err := db.Preload("Ingredients.Item").Preload("Nutrients").Preload("Instructions").Find(&recipes).Error; err != nil {
		return dtos.RecipesResponse{}, err
	}

	// If no recipes found in database, search Spoonacular API
	if len(recipes) == 0 {
		// Get ingredient names for the API call
		var ingredientNames []string
		for _, id := range ingredientIDs {
			var item models.Item
			if err := r.db.First(&item, id).Error; err == nil {
				ingredientNames = append(ingredientNames, item.Name)
			}
		}

		// Build API URL
		apiURL := fmt.Sprintf("%s/recipes/complexSearch", r.spoonacular.GetBaseURL())
		params := map[string]string{
			"apiKey":                r.spoonacular.GetAPIKey(),
			"addRecipeInstructions": "true",
			"addRecipeNutrition":    "true",
			"number":                "2",
		}

		if len(ingredientNames) > 0 {
			params["includeIngredients"] = strings.Join(ingredientNames, ",")
		}

		if query.Diet != "" {
			params["diet"] = query.Diet
		}

		if query.Title != "" {
			params["titleMatch"] = query.Title
		}

		// Make API request
		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			return dtos.RecipesResponse{}, err
		}

		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()

		resp, err := r.spoonacular.GetClient().Do(req)
		if err != nil {
			return dtos.RecipesResponse{}, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return dtos.RecipesResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		var apiResponse struct {
			Results []struct {
				ID                 int     `json:"id"`
				Title              string  `json:"title"`
				Image              string  `json:"image"`
				ReadyInMinutes     int     `json:"readyInMinutes"`
				PreparationMinutes int     `json:"preparationMinutes"`
				CookingMinutes     int     `json:"cookingMinutes"`
				Servings           float32 `json:"servings"`
				Summary            string  `json:"summary"`
				Vegan              bool    `json:"vegan"`
				Vegetarian         bool    `json:"vegetarian"`
				Nutrition          struct {
					Nutrients []struct {
						Name                string  `json:"name"`
						Amount              float64 `json:"amount"`
						Unit                string  `json:"unit"`
						PercentOfDailyNeeds float64 `json:"percentOfDailyNeeds"`
					} `json:"nutrients"`
					Ingredients []struct {
						ID     int     `json:"id"`
						Name   string  `json:"name"`
						Amount float64 `json:"amount"`
						Unit   string  `json:"unit"`
					} `json:"ingredients"`
				} `json:"nutrition"`
				AnalyzedInstructions []struct {
					Steps []struct {
						Number int    `json:"number"`
						Step   string `json:"step"`
					} `json:"steps"`
				} `json:"analyzedInstructions"`
			} `json:"results"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
			return dtos.RecipesResponse{}, err
		}

		totalCount = int64(len(apiResponse.Results))

		// Process each recipe from API
		for _, apiRecipe := range apiResponse.Results {
			// Check if recipe already exists
			var existingRecipe models.Recipe
			if err := r.db.Where("spoonacular_id = ?", apiRecipe.ID).First(&existingRecipe).Error; err == nil {
				recipes = append(recipes, existingRecipe)
				continue
			}

			// Create new recipe
			recipe := models.Recipe{
				SpoonacularID: uint(apiRecipe.ID),
				Title:         apiRecipe.Title,
				Summary:       apiRecipe.Summary,
				Image:         apiRecipe.Image,
				Servings:      apiRecipe.Servings,
				ReadyTime:     int16(apiRecipe.ReadyInMinutes),
				CookingTime:   int16(apiRecipe.CookingMinutes),
				PrepTime:      int16(apiRecipe.PreparationMinutes),
				Vegan:         apiRecipe.Vegan,
				Vegetarian:    apiRecipe.Vegetarian,
			}

			// Create nutrients
			recipe.Nutrients = make([]models.RecipeNutrient, len(apiRecipe.Nutrition.Nutrients))
			for i, n := range apiRecipe.Nutrition.Nutrients {
				if n.Name == "Calories" {
					recipe.KCal = float32(n.Amount)
				}
				recipe.Nutrients[i] = models.RecipeNutrient{
					Name:                n.Name,
					Amount:              float64(n.Amount),
					Unit:                n.Unit,
					PercentOfDailyNeeds: float64(n.PercentOfDailyNeeds),
				}
			}

			// Create instructions
			if len(apiRecipe.AnalyzedInstructions) > 0 {
				recipe.Instructions = make([]models.RecipeInstruction, len(apiRecipe.AnalyzedInstructions[0].Steps))
				for i, step := range apiRecipe.AnalyzedInstructions[0].Steps {
					recipe.Instructions[i] = models.RecipeInstruction{
						Number: uint(step.Number),
						Step:   step.Step,
					}
				}
			}

			// Create ingredients
			recipe.Ingredients = make([]models.RecipeItem, len(apiRecipe.Nutrition.Ingredients))
			for i, ing := range apiRecipe.Nutrition.Ingredients {
				// Check if item exists
				var item models.Item
				if err := r.db.Where("spoonacular_id = ?", ing.ID).First(&item).Error; err != nil {
					// Create new item
					item = models.Item{
						Name:          ing.Name,
						SpoonacularID: uint(ing.ID),
					}
					if err := r.db.Create(&item).Error; err != nil {
						return dtos.RecipesResponse{}, err
					}

					// Add to item queue for enrichment
					queueItem := models.QueueItem{
						ItemID:    item.ID,
						Name:      item.Name,
						CreatedAt: time.Now(),
						Priority:  models.DefaultPriority,
					}
					if err := r.itemQueue.AddItem(context.Background(), queueItem); err != nil {
						return dtos.RecipesResponse{}, err
					}
				}

				recipe.Ingredients[i] = models.RecipeItem{
					ItemID: item.ID,
					Amount: float32(ing.Amount),
					Unit:   ing.Unit,
				}
			}

			// Save recipe
			if err := r.db.Create(&recipe).Error; err != nil {
				return dtos.RecipesResponse{}, err
			}

			recipes = append(recipes, recipe)
		}
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

		nutrients := make([]dtos.RecipeNutrientResponse, len(recipe.Nutrients))
		for j, n := range recipe.Nutrients {
			nutrients[j] = dtos.RecipeNutrientResponse{
				Name:                n.Name,
				Amount:              n.Amount,
				Unit:                n.Unit,
				PercentOfDailyNeeds: n.PercentOfDailyNeeds,
			}
		}

		instructions := make([]dtos.RecipeInstructionResponse, len(recipe.Instructions))
		for j, inst := range recipe.Instructions {
			instructions[j] = dtos.RecipeInstructionResponse{
				Number: inst.Number,
				Step:   inst.Step,
			}
		}

		recipeResponses[i] = dtos.RecipeResponse{
			ID:            recipe.ID,
			Title:         recipe.Title,
			Summary:       recipe.Summary,
			SpoonacularID: recipe.SpoonacularID,
			Instructions:  instructions,
			Servings:      recipe.Servings,
			ReadyTime:     recipe.ReadyTime,
			CookingTime:   recipe.CookingTime,
			PrepTime:      recipe.PrepTime,
			Image:         recipe.Image,
			KCal:          recipe.KCal,
			Vegan:         recipe.Vegan,
			Vegetarian:    recipe.Vegetarian,
			Ingredients:   ingredients,
			Nutrients:     nutrients,
		}
	}

	return dtos.RecipesResponse{
		Recipes:    recipeResponses,
		Count:      int(totalCount),
		DietCounts: dietCounts,
	}, nil
}
