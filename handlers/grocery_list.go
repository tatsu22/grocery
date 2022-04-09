package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tatsu22/grocery/database"
	"github.com/tatsu22/grocery/model"
)

type GroceryListResp struct {
	Items     []model.GroceryListItem
	TotalCost float32
}

type AddRecipeRequest struct {
	Name string
}

func AddItemToList(c echo.Context) error {
	ctx := c.Request().Context()
	db := c.Get(dbContextKey).(database.Gorm)
	item := new(model.GroceryListItem)

	if err := c.Bind(item); err != nil {
		return err
	}

	insertedItem, err := db.AddItemToList(ctx, *item)
	if err != nil {
		resp := model.RequestError{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Error adding item to list: %+v", err),
		}
		return c.JSON(resp.Status, resp)
	}

	return c.JSON(http.StatusCreated, insertedItem)
}

func AddRecipeItemsToList(c echo.Context) error {
	ctx := c.Request().Context()
	db := c.Get(dbContextKey).(database.Gorm)
	recipeReq := new(AddRecipeRequest)

	if err := c.Bind(recipeReq); err != nil {
		return err
	}
	name := recipeReq.Name

	groceryItems, err := db.AddRecipeItemsToList(ctx, name)
	if err != nil {
		resp := model.RequestError{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("error adding recipe items to list: %+v", err),
		}
		return c.JSON(resp.Status, resp)
	}
	return c.JSON(http.StatusOK, groceryItems)
}

func GetGroceryList(c echo.Context) error {
	ctx := c.Request().Context()
	db := c.Get(dbContextKey).(database.Gorm)

	groceryItems, err := db.GetList(ctx)
	if err != nil {
		resp := model.RequestError{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Error retrieving grocery list: %+v", err),
		}
		return c.JSON(resp.Status, resp)
	}

	resp := GroceryListResp{
		Items:     groceryItems,
		TotalCost: GetTotalCostFromList(groceryItems),
	}

	return c.JSON(http.StatusOK, resp)

}

func SubtractItemFromList(c echo.Context) error {
	ctx := c.Request().Context()
	db := c.Get(dbContextKey).(database.Gorm)

	item := new(model.GroceryListItem)

	if err := c.Bind(item); err != nil {
		return err
	}

	updatedItem, err := db.SubtractItemFromList(ctx, *item)
	if err != nil {
		resp := model.RequestError{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("error subtracting item from grocery list: %+v", err),
		}
		return c.JSON(resp.Status, resp)
	}
	return c.JSON(http.StatusOK, updatedItem)
}

func DeleteItemFromList(c echo.Context) error {
	ctx := c.Request().Context()
	db := c.Get(dbContextKey).(database.Gorm)

	item := new(model.GroceryListItem)

	if err := c.Bind(item); err != nil {
		return err
	}

	deletedItem, err := db.DeleteItemFromList(ctx, *item)
	if err != nil {
		resp := model.RequestError{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("error deleting item from grocery list: %+v", err),
		}
		return c.JSON(resp.Status, resp)
	}
	return c.JSON(http.StatusOK, deletedItem)
}

func DeleteGroceryList(c echo.Context) error {
	ctx := c.Request().Context()
	db := c.Get(dbContextKey).(database.Gorm)

	items, err := db.DeleteList(ctx)
	if err != nil {
		resp := model.RequestError{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("error deleting grocery list: %+v", err),
		}
		return c.JSON(resp.Status, resp)
	}
	return c.JSON(http.StatusOK, items)
}

func GetTotalCostFromList(list []model.GroceryListItem) float32 {
	var cost float32
	for _, listItem := range list {
		cost += listItem.Cost * listItem.Number
	}
	return cost
}
