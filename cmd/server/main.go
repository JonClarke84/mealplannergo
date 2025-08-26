package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/JonClarke84/mealplannergo/pkg/config"
	"github.com/JonClarke84/mealplannergo/pkg/db"
	"github.com/JonClarke84/mealplannergo/pkg/handlers"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Display environment information
	fmt.Printf("Starting server in %s environment\n", cfg.Environment)
	fmt.Printf("Using database: %s\n", cfg.DatabaseName)

	// Initialize database connection
	mongoDB, err := db.NewMongoDB(cfg.MongoURI, cfg.DatabaseName)
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
	fmt.Printf("Server starting on port %s...\n", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}