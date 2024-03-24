package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
)

func serveTemplate(w http.ResponseWriter, thisWeeksMeals MealList) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// execute the template
	if err := tmpl.Execute(w, thisWeeksMeals); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// load env only in development
	// if os.Getenv("GO_ENV") != "production" {
	// 	if err := godotenv.Load(); err != nil {
	// 		fmt.Printf("Error loading .env file: %s\n", err)
	// 		return
	// 	}
	// }

	// db
	MONGO_URI := os.Getenv("GO_SHOPPING_MONGO_ATLAS_URI")

	fmt.Print("MONGO_URI:", MONGO_URI)

	client, err := db(MONGO_URI)
	if err != nil {
		fmt.Printf("We have experienced an error connecting to MongoDB: %s\n", err)
		return
	}

	defer client.Disconnect(context.TODO())

	// routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		thisWeeksMeals, err := getThisWeeksMeals(client)
		if err != nil {
			fmt.Printf("Error getting this week's meals: %s\n", err)
			return
		}
		serveTemplate(w, thisWeeksMeals)
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

		defer updateMeal(client, key, value)

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

	// start server
	fmt.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
