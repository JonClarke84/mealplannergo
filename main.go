package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
)

func serveTemplate(w http.ResponseWriter, newestList MealList) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// execute the template
	if err := tmpl.Execute(w, newestList); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func serveShoppingListTemplate(w http.ResponseWriter, shoppingList []string) {
	tmpl, err := template.ParseFiles("templates/shopping-list.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// execute the template
	if err := tmpl.Execute(w, shoppingList); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// db
	MONGO_URI := os.Getenv("GO_SHOPPING_MONGO_ATLAS_URI")

	client, err := db(MONGO_URI)
	if err != nil {
		fmt.Printf("We have experienced an error connecting to MongoDB, shutting down the server: %s\n", err)
		return
	}

	// routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		newestList, err := getNewestList(client)
		if err != nil {
			fmt.Printf("Error getting this week's meals: %s\n", err)
			return
		}
		serveTemplate(w, newestList)
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
			io.WriteString(w, fmt.Sprintf(
				"<input type='text' name='%s' id='%s-input' style='flex-grow: 1;' />"+
					"<button "+
					"class='failed' "+
					"hx-post='/meal' "+
					"hx-trigger='click' "+
					"hx-include='#%s-input' "+
					"hx-target='#%s-container' "+
					">💾</button>"+
					"<div class='error'>Failed to save %s: %s</div>",
				key, key, value, key, key, value))
			return
		}

		w.Header().Set("Content-Type", "text/html")
		io.WriteString(w, fmt.Sprintf(
			"<input type='text' name='%s' id='%s-input' style='flex-grow: 1;' value='%s' />"+
				"<button "+
				"class='saved' "+
				"hx-post='/meal' "+
				"hx-trigger='click' "+
				"hx-include='#%s-input' "+
				"hx-target='#%s-container' "+
				">💾</button>",
			key, key, value, key, key))
	})

	http.HandleFunc("/shopping-list", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			fmt.Printf("Error parsing shopping list: %s\n", err)
			return
		}

		// CREATE
		if r.Method == "POST" {
			item := r.PostFormValue("item")
			if err := addShoppingListItem(client, item); err != nil {
				http.Error(w, "Failed to add item", http.StatusInternalServerError)
				return
			}
			tmpl := template.Must(template.ParseFiles("templates/index.html"))
			tmpl.ExecuteTemplate(w, "shopping-list-item", item)
		}

		// EDIT
		if r.Method == "PUT" {
			oldItem := r.URL.Query().Get("shopping-list-item")
			newItem := r.PostFormValue("item")

			if err := updateShoppingListItem(client, oldItem, newItem); err != nil {
				http.Error(w, "Failed to update item", http.StatusInternalServerError)
				return
			}

			tmpl := template.Must(template.ParseFiles("templates/index.html"))
			tmpl.ExecuteTemplate(w, "shopping-list-item", newItem)
		}

		// DELETE
		if r.Method == "DELETE" {
			item := r.URL.Query().Get("shopping-list-item")
			if err := deleteShoppingListItem(client, item); err != nil {
				http.Error(w, "Failed to remove item", http.StatusInternalServerError)
				return
			}
			// respond 200 ok
			w.WriteHeader(http.StatusOK)
		}
	})

	// start server
	fmt.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}

	// close the database connection when done
	defer client.Disconnect(context.TODO())
}
