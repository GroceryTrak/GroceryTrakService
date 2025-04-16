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
		Ingredients:   make([]models.RecipeItem, len(req.Ingredients)),
		Nutrients:     make([]models.RecipeNutrient, len(req.Nutrients)),
		Instructions:  make([]models.RecipeInstruction, len(req.Instructions)),
	}

	for i, item := range req.Ingredients {
		recipe.Ingredients[i] = models.RecipeItem{
			ItemID: item.ItemID,
			Amount: item.Amount,
			Unit:   item.Unit,
		}
	}

	for i, n := range req.Nutrients {
		recipe.Nutrients[i] = models.RecipeNutrient{
			Name:                n.Name,
			Amount:              n.Amount,
			Unit:                n.Unit,
			PercentOfDailyNeeds: n.PercentOfDailyNeeds,
		}
	}

	for i, inst := range req.Instructions {
		recipe.Instructions[i] = models.RecipeInstruction{
			Number: inst.Number,
			Step:   inst.Step,
		}
	}

	if err := r.db.Create(&recipe).Error; err != nil {
		return nil, err
	}

	if err := r.db.Preload("Ingredients.Item").Preload("Nutrients").Preload("Instructions").First(&recipe, recipe.ID).Error; err != nil {
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

	if err := r.db.Where("recipe_id = ?", id).Delete(&models.RecipeNutrient{}).Error; err != nil {
		return nil, err
	}

	if err := r.db.Where("recipe_id = ?", id).Delete(&models.RecipeInstruction{}).Error; err != nil {
		return nil, err
	}

	modelNutrients := make([]models.RecipeNutrient, len(req.Nutrients))
	for i, n := range req.Nutrients {
		modelNutrients[i] = models.RecipeNutrient{
			RecipeID:            id,
			Name:                n.Name,
			Amount:              n.Amount,
			Unit:                n.Unit,
			PercentOfDailyNeeds: n.PercentOfDailyNeeds,
		}
	}

	modelInstructions := make([]models.RecipeInstruction, len(req.Instructions))
	for i, inst := range req.Instructions {
		modelInstructions[i] = models.RecipeInstruction{
			RecipeID: id,
			Number:   inst.Number,
			Step:     inst.Step,
		}
	}

	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Save(&recipe).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Create(&modelNutrients).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Create(&modelInstructions).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

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

func (r *RecipeRepositoryImpl) DeleteRecipe(id uint) error {
	if err := r.db.Where("recipe_id = ?", id).Delete(&models.RecipeNutrient{}).Error; err != nil {
		return err
	}
	if err := r.db.Where("recipe_id = ?", id).Delete(&models.RecipeInstruction{}).Error; err != nil {
		return err
	}
	return r.db.Delete(&models.Recipe{}, id).Error
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

	if err := db.Preload("Ingredients.Item").Preload("Nutrients").Preload("Instructions").Find(&recipes).Error; err != nil {
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
