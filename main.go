package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./htmx")))

	// http.HandleFunc("/generate", func(w http.ResponseWriter, r *http.Request) {
	// 	mealPlanHTML := getMealPlanHTML()
	// 	fmt.Fprint(w, mealPlanHTML)
	// })

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

		r.Header.Set("Content-Type", "text/html")
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

	fmt.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
