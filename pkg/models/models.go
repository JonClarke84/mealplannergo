package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ShoppingListItem represents a single item in a shopping list
type ShoppingListItem struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"ID,omitempty"`
	IDHex  string             `bson:"IDHex,omitempty" json:"IDHex,omitempty"`
	Item   string             `bson:"Item" json:"Item"`
	Ticked bool               `bson:"Ticked" json:"Ticked"`
}

// ShoppingList represents the entire shopping list
type ShoppingList struct {
	ShoppingList []ShoppingListItem
	ID           primitive.ObjectID
	SortOrder    []primitive.ObjectID
}

// ShoppingListDocument represents how the shopping list is stored in MongoDB
type ShoppingListDocument struct {
	ID           primitive.ObjectID   `bson:"_id" json:"id,omitempty"`
	ShoppingList []ShoppingListItem   `bson:"ShoppingList" json:"ShoppingList"`
	SortOrder    []primitive.ObjectID `bson:"SortOrder" json:"SortOrder"`
}

// Meal represents a single meal for a day
type Meal struct {
	Day  string
	Meal string
}

// MealPlan represents a collection of meals for a period
type MealPlan struct {
	Meals []Meal
}

// PageData represents the data structure passed to the HTML template
type PageData struct {
	MealPlan     []Meal
	ShoppingList []ShoppingListItem
}

// Order represents the position of an item in the shopping list
type Order struct {
	ID       string `json:"id"`
	Position int    `json:"position"`
}

// OrderUpdate represents a collection of orders for updating positions
type OrderUpdate struct {
	Order []Order `json:"order"`
}