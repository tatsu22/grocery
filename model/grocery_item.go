package model

type GroceryItem struct {
	Name    string
	Cost    float32
	Picture string
	Unit    string
}

type GroceryList struct {
	GroceryItem string
	Number      int
}
