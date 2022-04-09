package database

import (
	"context"
	"strings"

	"github.com/tatsu22/grocery/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (db *Gorm) AddItemToList(ctx context.Context, item model.GroceryListItem) (model.GroceryListItem, error) {
	tx := db.DB.WithContext(ctx)
	result := tx.Create(&item)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			tx.Transaction(func(tx2 *gorm.DB) error {
				var existingItem model.GroceryListItem
				tx2.Where("name = ? AND unit = ?", item.Name, item.Unit).First(&existingItem)
				item.Number += existingItem.Number
				result = tx2.Save(&item)
				return nil
			})
		}
	}

	return item, result.Error
}

func (db *Gorm) SubtractItemFromList(ctx context.Context, item model.GroceryListItem) (model.GroceryListItem, error) {
	tx := db.DB.WithContext(ctx)
	var existingItem model.GroceryListItem
	err := tx.Transaction(func(tx2 *gorm.DB) error {
		result := tx2.Where("name = ? and unit = ?", item.Name, item.Unit).First(&existingItem)
		if result.Error != nil {
			return result.Error
		}
		existingItem.Number -= item.Number
		if existingItem.Number >= 0 {
			result = tx2.Delete(&existingItem)
		} else {
			result = tx2.Save(&existingItem)
		}
		return result.Error
	})
	return existingItem, err
}

func (db *Gorm) DeleteItemFromList(ctx context.Context, item model.GroceryListItem) (model.GroceryListItem, error) {
	tx := db.DB.WithContext(ctx)
	result := tx.Clauses(clause.Returning{}).Delete(&item)
	return item, result.Error
}

func (db *Gorm) DeleteList(ctx context.Context) ([]model.GroceryListItem, error) {
	tx := db.DB.WithContext(ctx)
	var items []model.GroceryListItem
	result := tx.Clauses(clause.Returning{}).Where("1 = 1").Delete(&items)
	return items, result.Error
}

func (db *Gorm) AddRecipeItemsToList(ctx context.Context, name string) ([]model.GroceryListItem, error) {
	ingredients, err := db.GetIngredientsByRecipe(ctx, name)
	if err != nil {
		return nil, err
	}

	addedItems := make([]model.GroceryListItem, len(ingredients))
	for i, ing := range ingredients {
		grItem, err := db.GetGroceryItem(ctx, ing.Name, ing.Unit)
		if err != nil {
			return nil, err
		}
		listItem := IngAndItemToListItem(ing, grItem)
		listItem, err = db.AddItemToList(ctx, listItem)
		if err != nil {
			return nil, err
		}
		addedItems[i] = listItem
	}
	return addedItems, nil
}

func (db *Gorm) GetList(ctx context.Context) ([]model.GroceryListItem, error) {
	tx := db.DB.WithContext(ctx)

	var groceryList []model.GroceryListItem
	result := tx.Find(&groceryList)
	return groceryList, result.Error
}

func IngAndItemToListItem(ing model.RecipeIngredient, item model.GroceryItem) model.GroceryListItem {
	return model.GroceryListItem{
		Name:   item.Name,
		Unit:   item.Unit,
		Number: ing.Number,
		Cost:   item.Cost,
	}
}
