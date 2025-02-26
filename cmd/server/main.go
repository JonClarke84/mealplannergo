package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/JonClarke84/mealplannergo/pkg/db"
	"github.com/JonClarke84/mealplannergo/pkg/handlers"
)

func main() {
	// Get MongoDB connection URI from environment variable
	mongoURI := os.Getenv("GO_SHOPPING_MONGO_ATLAS_URI")
	if mongoURI == "" {
		fmt.Println("Error: GO_SHOPPING_MONGO_ATLAS_URI environment variable is not set")
		os.Exit(1)
	}

	// Initialize database connection
	mongoDB, err := db.NewMongoDB(mongoURI)
	if err != nil {
		fmt.Printf("Error connecting to MongoDB: %s\n", err)
		os.Exit(1)
	}
	defer mongoDB.Close()

	// Initialize handlers
	h := handlers.New(mongoDB)

	// Define routes
	http.HandleFunc("/", h.HomeHandler)
	http.HandleFunc("/meal", h.MealHandler)
	http.HandleFunc("/shopping-list", h.ShoppingListHandler)
	http.HandleFunc("/shopping-list/tick", h.ShoppingListTickHandler)
	http.HandleFunc("/shopping-list/sort", h.ShoppingListSortHandler)
	http.HandleFunc("/shopping-list/edit", h.ShoppingListEditHandler)

	// Serve static files
	publicFileServer := http.FileServer(http.Dir("./public"))
	http.Handle("/public/", http.StripPrefix("/public/", publicFileServer))

	// Start server
	port := "8080"
	fmt.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}