package tests

import (
	"github.com/JonClarke84/mealplannergo/pkg/db"
	"github.com/JonClarke84/mealplannergo/pkg/models"
	"github.com/stretchr/testify/mock"
)

// MockDB is a mock implementation of the DBInterface for testing
type MockDB struct {
	mock.Mock
}

// Ensure MockDB implements DBInterface
var _ db.DBInterface = (*MockDB)(nil)

// GetShoppingList mocks the GetShoppingList method
func (m *MockDB) GetShoppingList() ([]models.ShoppingListItem, error) {
	args := m.Called()
	return args.Get(0).([]models.ShoppingListItem), args.Error(1)
}

// GetShoppingListItemFromIDHex mocks the GetShoppingListItemFromIDHex method
func (m *MockDB) GetShoppingListItemFromIDHex(IDHex string) (models.ShoppingListItem, error) {
	args := m.Called(IDHex)
	return args.Get(0).(models.ShoppingListItem), args.Error(1)
}

// UpdateMeal mocks the UpdateMeal method
func (m *MockDB) UpdateMeal(day string, meal string) error {
	args := m.Called(day, meal)
	return args.Error(0)
}

// AddShoppingListItem mocks the AddShoppingListItem method
func (m *MockDB) AddShoppingListItem(itemName string) (models.ShoppingListItem, error) {
	args := m.Called(itemName)
	return args.Get(0).(models.ShoppingListItem), args.Error(1)
}

// AddShoppingListIdToShoppingListOrder mocks the AddShoppingListIdToShoppingListOrder method
func (m *MockDB) AddShoppingListIdToShoppingListOrder(itemId string) error {
	args := m.Called(itemId)
	return args.Error(0)
}

// UpdateShoppingListItem mocks the UpdateShoppingListItem method
func (m *MockDB) UpdateShoppingListItem(itemId string, newItem string) (models.ShoppingListItem, error) {
	args := m.Called(itemId, newItem)
	return args.Get(0).(models.ShoppingListItem), args.Error(1)
}

// DeleteShoppingListItem mocks the DeleteShoppingListItem method
func (m *MockDB) DeleteShoppingListItem(itemIDHex string) error {
	args := m.Called(itemIDHex)
	return args.Error(0)
}

// TickShoppingListItem mocks the TickShoppingListItem method
func (m *MockDB) TickShoppingListItem(itemId string, ticked bool) (models.ShoppingListItem, error) {
	args := m.Called(itemId, ticked)
	return args.Get(0).(models.ShoppingListItem), args.Error(1)
}

// GetMealPlan mocks the GetMealPlan method
func (m *MockDB) GetMealPlan() (models.MealPlan, error) {
	args := m.Called()
	return args.Get(0).(models.MealPlan), args.Error(1)
}

// SortShoppingList mocks the SortShoppingList method
func (m *MockDB) SortShoppingList(newOrder []models.Order) error {
	args := m.Called(newOrder)
	return args.Error(0)
}

// Close mocks the Close method
func (m *MockDB) Close() {
	m.Called()
}