package model

type Recipe struct {
	Name       string
	Directions string
	Picture    string
}

type RecipeIngredient struct {
	Name   string
	Recipe string
	Number int
}
