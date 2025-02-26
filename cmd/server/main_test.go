package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JonClarke84/mealplannergo/pkg/db/tests"
	"github.com/JonClarke84/mealplannergo/pkg/handlers"
	"github.com/JonClarke84/mealplannergo/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// setupTestServer creates a test HTTP server with mocked database
func setupTestServer(t *testing.T) (*httptest.Server, *tests.MockDB) {
	mockDB := new(tests.MockDB)
	
	h := handlers.New(mockDB)
	// Override template path to use test template
	h.TemplatePath = "./testdata/test_template.html"
	
	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.HomeHandler)
	mux.HandleFunc("/meal", h.MealHandler)
	mux.HandleFunc("/shopping-list", h.ShoppingListHandler)
	mux.HandleFunc("/shopping-list/tick", h.ShoppingListTickHandler)
	mux.HandleFunc("/shopping-list/sort", h.ShoppingListSortHandler)
	mux.HandleFunc("/shopping-list/edit", h.ShoppingListEditHandler)
	
	server := httptest.NewServer(mux)
	
	return server, mockDB
}

func TestHomeRoute(t *testing.T) {
	server, mockDB := setupTestServer(t)
	defer server.Close()
	
	// Setup mock data
	shoppingList := []models.ShoppingListItem{
		{
			ID:     primitive.NewObjectID(),
			IDHex:  "123",
			Item:   "Test Item",
			Ticked: false,
		},
	}
	
	mealPlan := models.MealPlan{
		Meals: []models.Meal{
			{
				Day:  "Monday",
				Meal: "Test Meal",
			},
		},
	}
	
	mockDB.On("GetShoppingList").Return(shoppingList, nil)
	mockDB.On("GetMealPlan").Return(mealPlan, nil)
	
	resp, err := http.Get(server.URL + "/")
	assert.NoError(t, err)
	defer resp.Body.Close()
	
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockDB.AssertExpectations(t)
}

func TestMealRoute(t *testing.T) {
	server, mockDB := setupTestServer(t)
	defer server.Close()
	
	mockDB.On("UpdateMeal", "Monday", "New Meal").Return(nil)
	
	formData := "Monday=New+Meal"
	resp, err := http.Post(server.URL+"/meal", "application/x-www-form-urlencoded", strings.NewReader(formData))
	assert.NoError(t, err)
	defer resp.Body.Close()
	
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockDB.AssertExpectations(t)
}

func TestShoppingListAddRoute(t *testing.T) {
	server, mockDB := setupTestServer(t)
	defer server.Close()
	
	newItem := models.ShoppingListItem{
		ID:     primitive.NewObjectID(),
		IDHex:  "123",
		Item:   "New Item",
		Ticked: false,
	}
	
	mockDB.On("AddShoppingListItem", "New Item").Return(newItem, nil)
	
	formData := "item=New+Item"
	resp, err := http.Post(server.URL+"/shopping-list", "application/x-www-form-urlencoded", strings.NewReader(formData))
	assert.NoError(t, err)
	defer resp.Body.Close()
	
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockDB.AssertExpectations(t)
}

func TestShoppingListTickRoute(t *testing.T) {
	server, mockDB := setupTestServer(t)
	defer server.Close()
	
	updatedItem := models.ShoppingListItem{
		ID:     primitive.NewObjectID(),
		IDHex:  "123",
		Item:   "Test Item",
		Ticked: true,
	}
	
	mockDB.On("TickShoppingListItem", "123", true).Return(updatedItem, nil)
	
	formData := "123=on"
	resp, err := http.Post(server.URL+"/shopping-list/tick", "application/x-www-form-urlencoded", strings.NewReader(formData))
	assert.NoError(t, err)
	defer resp.Body.Close()
	
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockDB.AssertExpectations(t)
}

func TestShoppingListSortRoute(t *testing.T) {
	server, mockDB := setupTestServer(t)
	defer server.Close()
	
	orders := []models.Order{
		{
			ID:       "123",
			Position: 1,
		},
		{
			ID:       "456",
			Position: 2,
		},
	}
	
	mockDB.On("SortShoppingList", mock.MatchedBy(func(o []models.Order) bool {
		return len(o) == 2 && o[0].ID == "123" && o[1].ID == "456"
	})).Return(nil)
	
	orderUpdate := models.OrderUpdate{Order: orders}
	jsonData, _ := json.Marshal(orderUpdate)
	
	resp, err := http.Post(
		server.URL+"/shopping-list/sort",
		"application/json",
		strings.NewReader(string(jsonData)),
	)
	assert.NoError(t, err)
	defer resp.Body.Close()
	
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockDB.AssertExpectations(t)
}

func TestShoppingListEditRoute(t *testing.T) {
	server, mockDB := setupTestServer(t)
	defer server.Close()
	
	updatedItem := models.ShoppingListItem{
		ID:     primitive.NewObjectID(),
		IDHex:  "123",
		Item:   "Updated Item",
		Ticked: false,
	}
	
	mockDB.On("UpdateShoppingListItem", "123", "Updated Item").Return(updatedItem, nil)
	
	formData := "123=Updated+Item"
	resp, err := http.Post(server.URL+"/shopping-list/edit", "application/x-www-form-urlencoded", strings.NewReader(formData))
	assert.NoError(t, err)
	defer resp.Body.Close()
	
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockDB.AssertExpectations(t)
}