package main

import (
	"fmt"
	"net/http"
)

func getMealPlanHTML() string {
	return `
		<h1>Your Meal Plan</h1>
		<ul>
			<li>Saturday: Pasta</li>
			<li>Sunday: Grilled Chicken</li>
			<li>Monday: Tacos</li>
			<li>Tuesday: Salmon</li>
			<li>Wednesday: Curry</li>
			<li>Thursday: Pizza</li>
			<li>Friday: Burgers</li>
			<li>Saturday: Soup</li>
		</ul>
	`
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./htmx")))
	http.HandleFunc("/generate", func(w http.ResponseWriter, r *http.Request) {
		// getShoppingListFromAI()
		mealPlanHTML := getMealPlanHTML()
		fmt.Fprint(w, mealPlanHTML)
	})

	fmt.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
