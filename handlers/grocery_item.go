package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tatsu22/grocery/database"
	"github.com/tatsu22/grocery/model"
)

const (
	dbContextKey = "__db"
)

func PostGroceryItem(c echo.Context) error {
	ctx := c.Request().Context()
	db := c.Get(dbContextKey).(database.Gorm)
	item := new(model.GroceryItem)

	if err := c.Bind(item); err != nil {
		return err
	}

	insertedItem, err := db.InsertGroceryItem(ctx, *item)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, insertedItem)
}

func GetGroceryItems(c echo.Context) error {
	db := c.Get(dbContextKey).(database.Gorm)

	name := c.QueryParam("name")
	ctx := c.Request().Context()

	items, err := db.SearchGroceryItems(ctx, name)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, items)
}
