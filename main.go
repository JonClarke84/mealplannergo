package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./htmx")))
	http.HandleFunc("/generate", func(w http.ResponseWriter, r *http.Request) {
		mealPlanHTML := getMealPlanHTML()
		fmt.Fprint(w, mealPlanHTML)
	})

	http.HandleFunc("/saturday", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			fmt.Printf("Error parsing form: %s\n", err)
			return
		}
		meal := r.Form.Get("meal")
		r.Header.Set("Content-Type", "text/html")
		io.WriteString(w, fmt.Sprintf("<p>Saturday's Meal: %s</p>", meal))
	})

	fmt.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
