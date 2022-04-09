package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tatsu22/grocery/database"
	"github.com/tatsu22/grocery/handlers"
)

const (
	dbContextKey = "__db"
)

func dbMiddleware(db database.Gorm) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(dbContextKey, db)
			return next(c)
		}
	}
}

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

	e.Use(dbMiddleware(db))

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	// Grocery Items
	e.POST("/groceryitems", handlers.PostGroceryItem)
	e.GET("/groceryitems", handlers.GetGroceryItems)

	// Grocery List
	e.POST("/grocerylist/addItem", handlers.AddItemToList)
	e.POST("/grocerylist/subtractItem", handlers.SubtractItemFromList)
	e.DELETE("/grocerylist/deleteItem", handlers.DeleteItemFromList)
	e.GET("/grocerylist", handlers.GetGroceryList)
	e.POST("/grocerylist/addRecipe", handlers.AddRecipeItemsToList)
	e.DELETE("/grocerylist", handlers.DeleteGroceryList)

	// Recipes
	e.POST("/recipes", handlers.PostRecipe)
	e.GET("/recipes", handlers.GetRecipe)
	e.DELETE("/recipes", handlers.DeleteRecipe)

	e.Logger.Fatal(e.Start(":8080"))
}
