package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

type ShoppingListItem struct {
  Item string `bson:"Item" json:"Item"`
  Ticked bool
  Id string `bson:"Id" json:"Id,omitempty"`
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

  //   fmt.Printf("pageData = %+v\n", pageData)
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
      fmt.Printf("item = %s\n", item)
			if err := addShoppingListItem(client, item); err != nil {
				http.Error(w, "Failed to create item", http.StatusInternalServerError)
				return
			}
			tmpl := template.Must(template.ParseFiles("templates/index.html"))
			tmpl.ExecuteTemplate(w, "shopping-list-item", item)
		}

		// DELETE
		if r.Method == "DELETE" {
			item := r.URL.Query().Get("item")
			if err := deleteShoppingListItem(client, item); err != nil {
				http.Error(w, "Failed to delete item", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
		}
	})

	http.HandleFunc("/shopping-list/sort", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			fmt.Printf("Error parsing shopping list: %s\n", err)
			return
		}

		var items []string

		for _, v := range r.PostForm {
			items = append(items, v...)
		}

		if err := sortShoppingList(client, items); err != nil {
			http.Error(w, "Failed to sort shopping list", http.StatusInternalServerError)
			return
		}

		newestList, err := getShoppingList(client)
		if err != nil {
			fmt.Printf("Error getting this week's shopping list: %s\n", err)
			return
		}

		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.ExecuteTemplate(w, "shopping-list", newestList)
	})

	http.HandleFunc("/shopping-list/edit", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			fmt.Printf("Error parsing form: %s\n", err)
			return
		}

		var oldItem string
		var updatedItem string
		for k, v := range r.PostForm {
			oldItem = k
			updatedItem = v[0]
			break
		}

		if err := updateShoppingListItem(client, oldItem, updatedItem); err != nil {
			http.Error(w, "Failed to update meal", http.StatusInternalServerError)
			return
		}

		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.ExecuteTemplate(w, "shopping-list-item", updatedItem)
	})

	// start server
	fmt.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}

	// close the database connection when done
	defer client.Disconnect(context.TODO())
}
