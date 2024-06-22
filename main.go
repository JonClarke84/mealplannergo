package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShoppingListItem struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"ID,omitempty"`
	IDHex  string             `bson:"IDHex,omitempty" json:"IDHex,omitempty"`
	Item   string             `bson:"Item" json:"Item"`
	Ticked bool               `bson:"Ticked" json:"Ticked"`
}

type ShoppingList struct {
	ShoppingList []ShoppingListItem
	ID           primitive.ObjectID
	SortOrder    []primitive.ObjectID
}

type Meal struct {
	Day  string
	Meal string
}

type MealPlan struct {
	Meals []Meal
}
type PageData struct {
	MealPlan     []Meal
	ShoppingList []ShoppingListItem
}

func main() {
	// client
	MONGO_URI := os.Getenv("GO_SHOPPING_MONGO_ATLAS_URI")

	client, err := db(MONGO_URI)
	if err != nil {
		fmt.Printf("We have experienced an error connecting to MongoDB, shutting down the server: %s\n", err)
		return
	}

	// routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		shoppingList, err := getShoppingList(client)
		if err != nil {
			fmt.Printf("Error getting this week's meals: %s\n", err)
			return
		}

		mealPlan, err := getMealPlan(client)
		if err != nil {
			fmt.Printf("Error getting meal plan: %s\n", err)
			return
		}

		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pageData := PageData{
			MealPlan:     mealPlan.Meals,
			ShoppingList: shoppingList,
		}

		tmpl.Execute(w, pageData)
	})

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))

	http.HandleFunc("/meal", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			fmt.Printf("Error parsing form: %s\n", err)
			return
		}

		var key string
		var value string
		for k, v := range r.PostForm {
			key = k
			value = v[0]
			break
		}

		if err := updateMeal(client, key, value); err != nil {
			http.Error(w, "Failed to update meal", http.StatusInternalServerError)
			return
		}

		updatedMeal := Meal{
			Day:  key,
			Meal: value,
		}

		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.ExecuteTemplate(w, "meal-input", updatedMeal)
	})

	http.HandleFunc("/shopping-list", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			fmt.Printf("Error parsing shopping list: %s\n", err)
			return
		}
		// CREATE
		if r.Method == "POST" {
			item := r.PostFormValue("item")
			newItem, err := addShoppingListItem(client, item)
			if err != nil {
				fmt.Printf("Error adding shopping list item: %s\n", err)
				return
			}
			tmpl := template.Must(template.ParseFiles("templates/index.html"))
			tmpl.ExecuteTemplate(w, "shopping-list-item", newItem)
		}

		// DELETE
		if r.Method == "DELETE" {
			item := r.URL.Query().Get("item")
			if err := deleteShoppingListItem(client, item); err != nil {
				http.Error(w, "Failed to delete item", http.StatusInternalServerError)
				return
			}
			shoppingList, err := getShoppingList(client)
			if err != nil {
				fmt.Printf("Error getting shopping list: %s\n", err)
				return
			}

			tmpl := template.Must(template.ParseFiles("templates/index.html"))
			tmpl.ExecuteTemplate(w, "shopping-list", shoppingList)
		}
	})

	http.HandleFunc("/shopping-list/tick", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			fmt.Printf("Error parsing shopping list: %s\n", err)
			return
		}
		var itemId string
		var ticked bool

		for k, v := range r.PostForm {
			itemId = k
			ticked = v[0] == "on"
			break
		}

		shoppingListItem, err := tickShoppingListItem(client, itemId, ticked)
		if err != nil {
			fmt.Printf("Error ticking shopping list item: %s\n", err)
			return
		}

		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.ExecuteTemplate(w, "shopping-list-item", shoppingListItem)
	})

	http.HandleFunc("/shopping-list/sort", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			fmt.Printf("Error parsing shopping list: %s\n", err)
			return
		}
		// collect all values of any field named "order" and put them in a slice
		var newOrder []int
		for _, v := range r.PostForm["order"] {
			orderInt, _ := strconv.Atoi(v)
			newOrder = append(newOrder, orderInt)
		}

		// if err := updateShoppingListOrder(client, newOrder); err != nil {
		// 	fmt.Printf("Error sorting shopping list: %s\n", err)
		// 	http.Error(w, "Failed to sort shopping list", http.StatusInternalServerError)
		// 	return
		// }

		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/shopping-list/edit", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			fmt.Printf("Error parsing form: %s\n", err)
			return
		}

		var itemId string
		var updatedItem string
		for k, v := range r.PostForm {
			itemId = k
			updatedItem = v[0]
			break
		}

		shoppingListItem, err := updateShoppingListItem(client, itemId, updatedItem)
		if err != nil {
			fmt.Printf("Error updating shopping list item: %s\n", err)
			return
		}

		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.ExecuteTemplate(w, "shopping-list-item", shoppingListItem)
	})

	// start server
	fmt.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}

	// close the database connection when done
	defer client.Disconnect(context.TODO())
}
