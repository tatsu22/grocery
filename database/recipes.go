package database

import (
	"context"

	"github.com/tatsu22/grocery/model"
)

func (db *Gorm) InsertRecipe(ctx context.Context, rec model.Recipe) (model.Recipe, error) {
	tx := db.DB.WithContext(ctx)
	result := tx.Create(&rec)
	return rec, result.Error
}

func (db *Gorm) DeleteRecipe(ctx context.Context, name string) (model.Recipe, error) {
	tx := db.DB.WithContext(ctx)
	rec := model.Recipe{
		Name: name,
	}

	result := tx.Where("recipe = ?", name).Delete(&model.RecipeIngredient{})
	if result.Error != nil {
		return rec, result.Error
	}

	result = tx.Delete(&rec)
	return rec, result.Error
}

func (db *Gorm) InsertRecipeIngredients(ctx context.Context, ing []model.RecipeIngredient) ([]model.RecipeIngredient, error) {
	tx := db.DB.WithContext(ctx)
	result := tx.Create(&ing)
	return ing, result.Error
}

func (db *Gorm) InsertFullRecipe(ctx context.Context, rec model.Recipe, ing []model.RecipeIngredient) (model.Recipe, []model.RecipeIngredient, error) {
	tx := db.DB.WithContext(ctx)
	result := tx.Create(&rec)
	if result.Error != nil {
		return rec, ing, result.Error
	}
	result = tx.Create(&ing)
	if result.Error != nil {
		tx.Delete(&rec)
		return rec, ing, result.Error
	}
	return rec, ing, nil
}

func (db *Gorm) SearchRecipes(ctx context.Context, name string) ([]model.Recipe, error) {
	tx := db.DB.WithContext(ctx)

	var recipes []model.Recipe
	result := tx.Where("name ILIKE ?", wrapString(name)).Find(&recipes)
	return recipes, result.Error
}

func (db *Gorm) GetIngredientsByRecipe(ctx context.Context, name string) ([]model.RecipeIngredient, error) {
	tx := db.DB.WithContext(ctx)

	var ing []model.RecipeIngredient
	result := tx.Where("recipe = ?", name).Find(&ing)
	return ing, result.Error
}
