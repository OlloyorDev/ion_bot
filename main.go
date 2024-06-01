package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Request structure
type Request struct {
	A int `json:"a"`
	B int `json:"b"`
}

// Response structure
type Response struct {
	Sum int `json:"sum"`
}

// handler function for POST requests
func sumHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	// Decode the incoming JSON payload
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Calculate the sum
	sum := req.A + req.B
	fmt.Printf("<><><><><><><><><><><<> %v\n", sum)

	// Prepare the response
	resp := Response{Sum: sum}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/sum", sumHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if PORT is not set
	}

	fmt.Printf("Server is running on port %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
