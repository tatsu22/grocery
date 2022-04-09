package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/tatsu22/grocery/database"
	"github.com/tatsu22/grocery/model"
)

type RecipeAdditionRequest struct {
	Name        string
	Directions  string
	Ingredients []RecipeIngredientRequest
	Picture     string
}

type RecipeIngredientRequest struct {
	Name   string
	Number float32
	Unit   string
}

type RecipeAdditionResponse struct {
	Recipe      model.Recipe
	Ingredients []model.RecipeIngredient
}

type RecipeSearchResponse struct {
	Recipes []RecipeAdditionRequest
}

func PostRecipe(c echo.Context) error {
	ctx := c.Request().Context()
	db := c.Get(dbContextKey).(database.Gorm)
	req := new(RecipeAdditionRequest)

	if err := c.Bind(req); err != nil {
		resp := model.RequestError{
			Status:  http.StatusInternalServerError,
			Message: "Unable to unmarshall request",
		}
		return c.JSON(resp.Status, resp)
	}
	if req.Name == "" || req.Directions == "" || req.Ingredients == nil || len(req.Ingredients) == 0 {
		resp := model.RequestError{
			Status:  http.StatusBadRequest,
			Message: "Must have name, directions, and ingredients in recipe",
		}
		return c.JSON(resp.Status, resp)
	}

	recipe := model.Recipe{
		Name:       strings.ReplaceAll(req.Name, " ", "_"),
		Directions: req.Directions,
		Picture:    req.Picture,
	}

	ingredients := make([]model.RecipeIngredient, len(req.Ingredients))
	for i, ing := range req.Ingredients {
		newIng := model.RecipeIngredient{
			Name:   ing.Name,
			Number: ing.Number,
			Recipe: recipe.Name,
			Unit:   ing.Unit,
		}
		ingredients[i] = newIng
	}

	recipe, ingredients, err := db.InsertFullRecipe(ctx, recipe, ingredients)
	if err != nil {
		resp := model.RequestError{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Unable to save recipe to database: +%v", err),
		}
		return c.JSON(resp.Status, resp)
	}

	resp := RecipeAdditionResponse{
		Recipe:      recipe,
		Ingredients: ingredients,
	}

	return c.JSON(http.StatusCreated, resp)
}

func GetRecipe(c echo.Context) error {
	db := c.Get(dbContextKey).(database.Gorm)

	name := c.QueryParam("name")
	ctx := c.Request().Context()

	recipes, err := db.SearchRecipes(ctx, name)
	if err != nil {
		resp := model.RequestError{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Error searching in database for recipes: %+v", err),
		}
		return c.JSON(resp.Status, resp)
	}

	foundRecipes := make([]RecipeAdditionRequest, len(recipes))

	for i, rec := range recipes {
		ings, err := db.GetIngredientsByRecipe(ctx, rec.Name)
		if err != nil {
			resp := model.RequestError{
				Status:  http.StatusInternalServerError,
				Message: fmt.Sprintf("Error searching in database for recipes: %+v", err),
			}
			return c.JSON(resp.Status, resp)
		}
		ingreds := make([]RecipeIngredientRequest, len(ings))
		for i, ing := range ings {
			ingreds[i] = RecipeIngredientRequest{
				Name:   ing.Name,
				Number: ing.Number,
				Unit:   ing.Unit,
			}
		}
		foundRecipes[i] = RecipeAdditionRequest{
			Name:        rec.Name,
			Directions:  rec.Directions,
			Picture:     rec.Picture,
			Ingredients: ingreds,
		}

	}

	resp := RecipeSearchResponse{
		Recipes: foundRecipes,
	}

	return c.JSON(http.StatusOK, resp)
}

func DeleteRecipe(c echo.Context) error {
	db := c.Get(dbContextKey).(database.Gorm)

	name := c.QueryParam("name")
	ctx := c.Request().Context()

	rec, err := db.DeleteRecipe(ctx, name)
	if err != nil {
		resp := model.RequestError{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("error deleting recipe from database: %+v", err),
		}
		return c.JSON(resp.Status, resp)
	}

	return c.JSON(http.StatusOK, rec)
}
