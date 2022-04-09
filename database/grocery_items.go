package database

import (
	"context"

	"github.com/tatsu22/grocery/model"
)

func (db *Gorm) InsertGroceryItem(ctx context.Context, item model.GroceryItem) (model.GroceryItem, error) {
	tx := db.DB.WithContext(ctx)
	result := tx.Create(&item)
	return item, result.Error
}

func (db *Gorm) UpdateGroceryItem(ctx context.Context, item model.GroceryItem) (model.GroceryItem, error) {
	tx := db.DB.WithContext(ctx)
	result := tx.Save(&item)
	return item, result.Error
}

func (db *Gorm) SearchGroceryItems(ctx context.Context, name string) ([]model.GroceryItem, error) {
	var results []model.GroceryItem
	tx := db.DB.WithContext(ctx)
	result := tx.Where("name ILIKE ?", wrapString(name)).Find(&results)

	return results, result.Error
}

func (db *Gorm) GetGroceryItem(ctx context.Context, name, unit string) (model.GroceryItem, error) {
	tx := db.DB.WithContext(ctx)
	item := model.GroceryItem{
		Name: name,
		Unit: unit,
	}
	result := tx.Find(&item)
	return item, result.Error
}

func (db *Gorm) GetAllGroceryItems(ctx context.Context) ([]model.GroceryItem, error) {
	var groceryItems []model.GroceryItem

	tx := db.DB.WithContext(ctx)
	result := tx.Find(&groceryItems)

	return groceryItems, result.Error
}
