package db

import (
	"github.com/JonClarke84/mealplannergo/pkg/models"
)

// DBInterface defines the interface for database operations
// This allows for easy mocking in tests
type DBInterface interface {
	GetShoppingList() ([]models.ShoppingListItem, error)
	GetShoppingListItemFromIDHex(IDHex string) (models.ShoppingListItem, error)
	UpdateMeal(day string, meal string) error
	AddShoppingListItem(itemName string) (models.ShoppingListItem, error)
	AddShoppingListIdToShoppingListOrder(itemId string) error
	UpdateShoppingListItem(itemId string, newItem string) (models.ShoppingListItem, error)
	DeleteShoppingListItem(itemIDHex string) error
	TickShoppingListItem(itemId string, ticked bool) (models.ShoppingListItem, error)
	GetMealPlan() (models.MealPlan, error)
	SortShoppingList(newOrder []models.Order) error
	Close()
}