package main

import (
	"fmt"
	"net/http"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "*") // Allow all headers
}

func main() {
	http.HandleFunc("/generate", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		mealPlanHTML := `
<!DOCTYPE html>
<html>
<head>
<title>Weekly Meal Plan</title>
</head>
<body>
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
</body>
</html>
`
		fmt.Fprint(w, mealPlanHTML)
	})

	fmt.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
