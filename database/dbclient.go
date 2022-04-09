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

func New() (Gorm, error) {
	dsn := "host=localhost user=postgres password=secret dbname=grocery port=5432 sslmode=disable"
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
