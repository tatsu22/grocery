package model

type GroceryItem struct {
	Name    string `gorm:"primaryKey"`
	Cost    float32
	Picture string
	Unit    string `gorm:"primaryKey"`
}

type GroceryListItem struct {
	Name   string `gorm:"primaryKey"`
	Unit   string `gorm:"primaryKey"`
	Cost   float32
	Number float32
}

func (GroceryListItem) TableName() string {
	return "grocery_list"
}
