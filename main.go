package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tatsu22/grocery/database"
	"github.com/tatsu22/grocery/model"
)

func main() {
	fmt.Println("vim-go")

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	db, err := database.New()
	if err != nil {
		panic(err)
	}

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	e.POST("/groceryitems", func(c echo.Context) error {
		ctx := c.Request().Context()
		item := new(model.GroceryItem)

		if err := c.Bind(item); err != nil {
			return err
		}

		insertedItem, err := db.InsertGroceryItem(ctx, *item)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, insertedItem)
	})

	e.GET("/groceryitems", func(c echo.Context) error {
		name := c.QueryParam("name")
		ctx := c.Request().Context()
		items, err := db.SearchGroceryItems(ctx, name)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, items)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
