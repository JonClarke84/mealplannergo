package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/JonClarke84/mealplannergo/pkg/db"
	"github.com/JonClarke84/mealplannergo/pkg/models"
)

// Handler contains all the dependencies needed for handling HTTP requests
type Handler struct {
	DB db.DBInterface
	// Template path relative to where the app is run from
	TemplatePath string
}

// New creates a new Handler with the given database connection
func New(db db.DBInterface) *Handler {
	return &Handler{
		DB:           db,
		TemplatePath: "./pkg/templates/index.html",
	}
}

// HomeHandler handles the root path request
func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	shoppingList, err := h.DB.GetShoppingList()
	if err != nil {
		fmt.Printf("Error getting this week's meals: %s\n", err)
		http.Error(w, "Failed to get shopping list", http.StatusInternalServerError)
		return
	}

	mealPlan, err := h.DB.GetMealPlan()
	if err != nil {
		fmt.Printf("Error getting meal plan: %s\n", err)
		http.Error(w, "Failed to get meal plan", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(h.TemplatePath)
	if err != nil {
		fmt.Printf("Error parsing template: %s (path: %s)\n", err, h.TemplatePath)
		absPath, _ := filepath.Abs(h.TemplatePath)
		http.Error(w, fmt.Sprintf("Template error: %s (tried: %s)", err, absPath), http.StatusInternalServerError)
		return
	}

	pageData := models.PageData{
		MealPlan:     mealPlan.Meals,
		ShoppingList: shoppingList,
	}

	tmpl.Execute(w, pageData)
}

// MealHandler handles updating a meal
func (h *Handler) MealHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Printf("Error parsing form: %s\n", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	var key string
	var value string
	for k, v := range r.PostForm {
		key = k
		value = v[0]
		break
	}

	if err := h.DB.UpdateMeal(key, value); err != nil {
		http.Error(w, "Failed to update meal", http.StatusInternalServerError)
		return
	}

	updatedMeal := models.Meal{
		Day:  key,
		Meal: value,
	}

	tmpl := template.Must(template.ParseFiles(h.TemplatePath))
	tmpl.ExecuteTemplate(w, "meal-input", updatedMeal)
}

// ShoppingListHandler handles operations on the shopping list
func (h *Handler) ShoppingListHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Printf("Error parsing shopping list: %s\n", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// CREATE
	if r.Method == "POST" {
		item := r.PostFormValue("item")
		newItem, err := h.DB.AddShoppingListItem(item)
		if err != nil {
			fmt.Printf("Error adding shopping list item: %s\n", err)
			http.Error(w, "Failed to add item", http.StatusInternalServerError)
			return
		}
		tmpl := template.Must(template.ParseFiles(h.TemplatePath))
		tmpl.ExecuteTemplate(w, "shopping-list-item", newItem)
	}

	// DELETE
	if r.Method == "DELETE" {
		item := r.URL.Query().Get("item")
		if err := h.DB.DeleteShoppingListItem(item); err != nil {
			http.Error(w, "Failed to delete item", http.StatusInternalServerError)
			return
		}
		shoppingList, err := h.DB.GetShoppingList()
		if err != nil {
			fmt.Printf("Error getting shopping list: %s\n", err)
			http.Error(w, "Failed to get updated shopping list", http.StatusInternalServerError)
			return
		}

		tmpl := template.Must(template.ParseFiles(h.TemplatePath))
		tmpl.ExecuteTemplate(w, "shopping-list", shoppingList)
	}
}

// ShoppingListTickHandler handles toggling a shopping list item's ticked state
func (h *Handler) ShoppingListTickHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Printf("Error parsing shopping list: %s\n", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	var itemId string
	var ticked bool

	for k, v := range r.PostForm {
		itemId = k
		ticked = v[0] == "on"
		break
	}

	shoppingListItem, err := h.DB.TickShoppingListItem(itemId, ticked)
	if err != nil {
		fmt.Printf("Error ticking shopping list item: %s\n", err)
		http.Error(w, "Failed to update item", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(h.TemplatePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "shopping-list-item", shoppingListItem)
}

// ShoppingListSortHandler handles reordering shopping list items
func (h *Handler) ShoppingListSortHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var updates models.OrderUpdate
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the order in the database
	err = h.DB.SortShoppingList(updates.Order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Respond with OK status
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// ShoppingListEditHandler handles editing shopping list items
func (h *Handler) ShoppingListEditHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Printf("Error parsing form: %s\n", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	var itemId string
	var updatedItem string
	for k, v := range r.PostForm {
		itemId = k
		updatedItem = v[0]
		break
	}

	shoppingListItem, err := h.DB.UpdateShoppingListItem(itemId, updatedItem)
	if err != nil {
		fmt.Printf("Error updating shopping list item: %s\n", err)
		http.Error(w, "Failed to update item", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles(h.TemplatePath))
	tmpl.ExecuteTemplate(w, "shopping-list-item", shoppingListItem)
}