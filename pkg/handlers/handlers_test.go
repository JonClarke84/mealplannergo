package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JonClarke84/mealplannergo/pkg/db/tests"
	"github.com/JonClarke84/mealplannergo/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNew(t *testing.T) {
	mockDB := new(tests.MockDB)
	handler := New(mockDB)
	
	assert.Equal(t, mockDB, handler.DB, "Handler should use the provided DB")
	assert.Equal(t, "./pkg/templates/index.html", handler.TemplatePath, "Handler should use the default template path")
}

func TestHomeHandler(t *testing.T) {
	// Setup mock DB
	mockDB := new(tests.MockDB)
	
	// Create test data
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
	
	// Set expectations on mock
	mockDB.On("GetShoppingList").Return(shoppingList, nil)
	mockDB.On("GetMealPlan").Return(mealPlan, nil)
	
	// Create handler with mock
	handler := &Handler{
		DB:           mockDB,
		TemplatePath: "../../cmd/server/testdata/test_template.html", // Use test template
	}
	
	// Create test request and response recorder
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	
	// Call handler
	handler.HomeHandler(w, req)
	
	// Assert expectations were met
	mockDB.AssertExpectations(t)
	
	// Check response
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	// Note: We can't easily test the template rendering without implementing a custom template
}

func TestHomeHandlerShoppingListError(t *testing.T) {
	mockDB := new(tests.MockDB)
	mockDB.On("GetShoppingList").Return([]models.ShoppingListItem{}, errors.New("database error"))
	
	handler := &Handler{
		DB:           mockDB,
		TemplatePath: "../../cmd/server/testdata/test_template.html", // Use test template
	}
	
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	
	handler.HomeHandler(w, req)
	
	mockDB.AssertExpectations(t)
	
	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestHomeHandlerMealPlanError(t *testing.T) {
	mockDB := new(tests.MockDB)
	
	shoppingList := []models.ShoppingListItem{
		{
			ID:     primitive.NewObjectID(),
			IDHex:  "123",
			Item:   "Test Item",
			Ticked: false,
		},
	}
	
	mockDB.On("GetShoppingList").Return(shoppingList, nil)
	mockDB.On("GetMealPlan").Return(models.MealPlan{}, errors.New("database error"))
	
	handler := &Handler{
		DB:           mockDB,
		TemplatePath: "../../cmd/server/testdata/test_template.html", // Use test template
	}
	
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	
	handler.HomeHandler(w, req)
	
	mockDB.AssertExpectations(t)
	
	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestMealHandler(t *testing.T) {
	mockDB := new(tests.MockDB)
	mockDB.On("UpdateMeal", "Monday", "New Meal").Return(nil)
	
	handler := &Handler{
		DB:           mockDB,
		TemplatePath: "../../cmd/server/testdata/test_template.html", // Use test template
	}
	
	// Create form data
	formData := "Monday=New+Meal"
	req := httptest.NewRequest("POST", "/meal", strings.NewReader(formData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	
	handler.MealHandler(w, req)
	
	mockDB.AssertExpectations(t)
	
	// Note: We can't easily test the template execution without implementing a custom template
}

func TestMealHandlerUpdateError(t *testing.T) {
	mockDB := new(tests.MockDB)
	mockDB.On("UpdateMeal", "Monday", "New Meal").Return(errors.New("update error"))
	
	handler := &Handler{
		DB:           mockDB,
		TemplatePath: "../../cmd/server/testdata/test_template.html", // Use test template
	}
	
	formData := "Monday=New+Meal"
	req := httptest.NewRequest("POST", "/meal", strings.NewReader(formData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	
	handler.MealHandler(w, req)
	
	mockDB.AssertExpectations(t)
	
	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestShoppingListHandler_Post(t *testing.T) {
	mockDB := new(tests.MockDB)
	
	newItem := models.ShoppingListItem{
		ID:     primitive.NewObjectID(),
		IDHex:  "123",
		Item:   "New Item",
		Ticked: false,
	}
	
	mockDB.On("AddShoppingListItem", "New Item").Return(newItem, nil)
	
	handler := &Handler{
		DB:           mockDB,
		TemplatePath: "../../cmd/server/testdata/test_template.html", // Use test template
	}
	
	formData := "item=New+Item"
	req := httptest.NewRequest("POST", "/shopping-list", strings.NewReader(formData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	
	handler.ShoppingListHandler(w, req)
	
	mockDB.AssertExpectations(t)
	
	// Note: We can't easily test the template execution without implementing a custom template
}

func TestShoppingListHandler_Delete(t *testing.T) {
	mockDB := new(tests.MockDB)
	
	shoppingList := []models.ShoppingListItem{
		{
			ID:     primitive.NewObjectID(),
			IDHex:  "456",
			Item:   "Remaining Item",
			Ticked: false,
		},
	}
	
	mockDB.On("DeleteShoppingListItem", "123").Return(nil)
	mockDB.On("GetShoppingList").Return(shoppingList, nil)
	
	handler := &Handler{
		DB:           mockDB,
		TemplatePath: "../../cmd/server/testdata/test_template.html", // Use test template
	}
	
	req := httptest.NewRequest("DELETE", "/shopping-list?item=123", nil)
	w := httptest.NewRecorder()
	
	handler.ShoppingListHandler(w, req)
	
	mockDB.AssertExpectations(t)
	
	// Note: We can't easily test the template execution without implementing a custom template
}

func TestShoppingListTickHandler(t *testing.T) {
	mockDB := new(tests.MockDB)
	
	updatedItem := models.ShoppingListItem{
		ID:     primitive.NewObjectID(),
		IDHex:  "123",
		Item:   "Test Item",
		Ticked: true,
	}
	
	mockDB.On("TickShoppingListItem", "123", true).Return(updatedItem, nil)
	
	handler := &Handler{
		DB:           mockDB,
		TemplatePath: "../../cmd/server/testdata/test_template.html", // Use test template
	}
	
	formData := "123=on"
	req := httptest.NewRequest("POST", "/shopping-list/tick", strings.NewReader(formData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	
	handler.ShoppingListTickHandler(w, req)
	
	mockDB.AssertExpectations(t)
	
	// Note: We can't easily test the template execution without implementing a custom template
}

func TestShoppingListSortHandler(t *testing.T) {
	mockDB := new(tests.MockDB)
	
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
	
	handler := &Handler{
		DB:           mockDB,
		TemplatePath: "../../cmd/server/testdata/test_template.html", // Use test template
	}
	
	// Create JSON payload
	orderUpdate := models.OrderUpdate{Order: orders}
	jsonData, _ := json.Marshal(orderUpdate)
	
	req := httptest.NewRequest("POST", "/shopping-list/sort", strings.NewReader(string(jsonData)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	handler.ShoppingListSortHandler(w, req)
	
	mockDB.AssertExpectations(t)
	
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestShoppingListEditHandler(t *testing.T) {
	mockDB := new(tests.MockDB)
	
	updatedItem := models.ShoppingListItem{
		ID:     primitive.NewObjectID(),
		IDHex:  "123",
		Item:   "Updated Item",
		Ticked: false,
	}
	
	mockDB.On("UpdateShoppingListItem", "123", "Updated Item").Return(updatedItem, nil)
	
	handler := &Handler{
		DB:           mockDB,
		TemplatePath: "../../cmd/server/testdata/test_template.html", // Use test template
	}
	
	formData := "123=Updated+Item"
	req := httptest.NewRequest("POST", "/shopping-list/edit", strings.NewReader(formData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	
	handler.ShoppingListEditHandler(w, req)
	
	mockDB.AssertExpectations(t)
	
	// Note: We can't easily test the template execution without implementing a custom template
}