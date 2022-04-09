package model

type Recipe struct {
	Name       string `gorm:"primaryKey"`
	Directions string
	Picture    string
}

type RecipeIngredient struct {
	Name   string `gorm:"column:grocery_item;primaryKey"`
	Recipe string `gorm:"primaryKey"`
	Number float32
	Unit   string `gorm:"primaryKey"`
}
