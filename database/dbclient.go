package database

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/tatsu22/grocery/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Gorm struct {
	DB *gorm.DB
}

type DBClient interface {
	InsertGroceryItem(ctx context.Context, item model.GroceryItem) (model.GroceryItem, error)
	GetAllGroceryItems(ctx context.Context) ([]model.GroceryItem, error)
	SearchGroceryItems(ctx context.Context) ([]model.GroceryItem, error)
}

func (db *Gorm) InsertGroceryItem(ctx context.Context, item model.GroceryItem) (model.GroceryItem, error) {
	tx := db.DB.WithContext(ctx)
	result := tx.Create(&item)
	return item, result.Error
}

func (db *Gorm) SearchGroceryItems(ctx context.Context, name string) ([]model.GroceryItem, error) {
	var results []model.GroceryItem
	tx := db.DB.WithContext(ctx)
	result := tx.Where("name ILIKE ?", wrapString(name)).Find(&results)

	return results, result.Error
}

func (db *Gorm) GetAllGroceryItems(ctx context.Context) ([]model.GroceryItem, error) {
	var groceryItems []model.GroceryItem

	tx := db.DB.WithContext(ctx)
	result := tx.Find(&groceryItems)

	return groceryItems, result.Error
}

func New() (Gorm, error) {
	dsn := "host=localhost user=postgres password=password dbname=grocery port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		logrus.WithError(err).Errorf("error connecting to DB")
		return Gorm{}, err
	}

	return Gorm{
		DB: db,
	}, nil
}

func wrapString(str string) string {
	return "%" + str + "%"
}
