package handlers

import (
	"fmt"
	"net/http"
	"strings"

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

	item.Name = strings.Title(strings.ToLower(item.Name))
	item.Unit = strings.ToLower(item.Unit)

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

func DeleteGroceryItem(c echo.Context) error {
	db := c.Get(dbContextKey).(database.Gorm)
	ctx := c.Request().Context()

	item := new(model.GroceryItem)

	if err := c.Bind(item); err != nil {
		return err
	}

	deletedItem, err := db.DeleteGroceryItem(ctx, *item)
	if err != nil {
		resp := model.RequestError{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("error deleting item from store: %+v", err),
		}
		return c.JSON(resp.Status, resp)
	}
	return c.JSON(http.StatusOK, deletedItem)

}
