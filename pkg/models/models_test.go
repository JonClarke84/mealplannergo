package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestShoppingListItem(t *testing.T) {
	id := primitive.NewObjectID()
	item := ShoppingListItem{
		ID:     id,
		IDHex:  id.Hex(),
		Item:   "Test Item",
		Ticked: false,
	}

	assert.Equal(t, id, item.ID, "ID should match")
	assert.Equal(t, id.Hex(), item.IDHex, "IDHex should match")
	assert.Equal(t, "Test Item", item.Item, "Item should match")
	assert.False(t, item.Ticked, "Ticked should be false")
}

func TestShoppingList(t *testing.T) {
	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	listID := primitive.NewObjectID()

	items := []ShoppingListItem{
		{
			ID:     id1,
			IDHex:  id1.Hex(),
			Item:   "Item 1",
			Ticked: false,
		},
		{
			ID:     id2,
			IDHex:  id2.Hex(),
			Item:   "Item 2",
			Ticked: true,
		},
	}

	list := ShoppingList{
		ShoppingList: items,
		ID:           listID,
		SortOrder:    []primitive.ObjectID{id1, id2},
	}

	assert.Len(t, list.ShoppingList, 2, "ShoppingList should have 2 items")
	assert.Equal(t, listID, list.ID, "ID should match")
	assert.Len(t, list.SortOrder, 2, "SortOrder should have 2 items")
}

func TestMeal(t *testing.T) {
	meal := Meal{
		Day:  "Monday",
		Meal: "Pasta",
	}

	assert.Equal(t, "Monday", meal.Day, "Day should match")
	assert.Equal(t, "Pasta", meal.Meal, "Meal should match")
}

func TestMealPlan(t *testing.T) {
	meals := []Meal{
		{
			Day:  "Monday",
			Meal: "Pasta",
		},
		{
			Day:  "Tuesday",
			Meal: "Pizza",
		},
	}

	mealPlan := MealPlan{
		Meals: meals,
	}

	assert.Len(t, mealPlan.Meals, 2, "MealPlan should have 2 meals")
	assert.Equal(t, "Monday", mealPlan.Meals[0].Day, "First meal day should be Monday")
	assert.Equal(t, "Pizza", mealPlan.Meals[1].Meal, "Second meal should be Pizza")
}

func TestPageData(t *testing.T) {
	id := primitive.NewObjectID()
	
	meals := []Meal{
		{
			Day:  "Monday",
			Meal: "Pasta",
		},
	}

	items := []ShoppingListItem{
		{
			ID:     id,
			IDHex:  id.Hex(),
			Item:   "Pasta",
			Ticked: false,
		},
	}

	pageData := PageData{
		MealPlan:     meals,
		ShoppingList: items,
	}

	assert.Len(t, pageData.MealPlan, 1, "PageData should have 1 meal")
	assert.Len(t, pageData.ShoppingList, 1, "PageData should have 1 shopping list item")
}

func TestOrder(t *testing.T) {
	order := Order{
		ID:       "123",
		Position: 1,
	}

	assert.Equal(t, "123", order.ID, "ID should match")
	assert.Equal(t, 1, order.Position, "Position should match")
}

func TestOrderUpdate(t *testing.T) {
	orders := []Order{
		{
			ID:       "123",
			Position: 1,
		},
		{
			ID:       "456",
			Position: 2,
		},
	}

	orderUpdate := OrderUpdate{
		Order: orders,
	}

	assert.Len(t, orderUpdate.Order, 2, "OrderUpdate should have 2 orders")
	assert.Equal(t, "123", orderUpdate.Order[0].ID, "First order ID should match")
	assert.Equal(t, 2, orderUpdate.Order[1].Position, "Second order position should match")
}