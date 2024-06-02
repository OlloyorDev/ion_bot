package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
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

// Global variable to be reset every 12 minutes
var (
	globalVar int
	mu        sync.Mutex
)

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
	fmt.Printf("Sum calculated: %d\n", sum)

	// Prepare the response
	resp := Response{Sum: sum}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

	// Send POST request to FCM
	sendNotification(sum)

	// Update the global variable
	mu.Lock()
	globalVar += sum
	mu.Unlock()
}

// sendNotification sends a POST request to the FCM endpoint
func sendNotification(sum int) {
	url := "https://fcm.googleapis.com/fcm/send"

	// Request body
	data := map[string]interface{}{
		"to": "en9ec3fm30l3gKakrC0axl:APA91bHGR-hIzXVoiQ-5JCCew4-qtKSJiv5GvjlxQNWPxNJv37yf9CIR97zSXYgNsgcU2W09PLR1OJ7XoZJGXEhziJ7ILlRzT5DQ-uzqqfsRSBuj6YQ1bnc_aNr52JkR3hUeQcP3Spjx",
		"notification": map[string]string{
			"title":        "Ion Bot",
			"body":         fmt.Sprintf("Hisob %d", sum),
			"click_action": "FLUTTER_NOTIFICATION_CLICK",
		},
		"data": map[string]string{
			"user_message": fmt.Sprintf("notification kelyaptoi buyurtmadan, sum: %d", sum),
		},
	}

	// Marshal JSON data
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Add your FCM server key as the Authorization header
	req.Header.Set("Authorization", "key=AAAAuI2nL_w:APA91bEfwqmmsd0TSjKuXnKETb133lmjKPIGCXILjIxwNlv8rf9ECPRAagohUZVuMvsJWXg1wpGeviJm6JMTRK_HmKWUysoTLmi_01RCdlyIRG2RVucU57AKpizhIpSRd_KGJ_HyYqXU")

	// Send HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	// Print response status
	fmt.Println("FCM server response:", resp.Status)
}

// handler function for GET requests to /refresh
func refreshHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the current date and time
	currentTime := time.Now()
	formattedTime := currentTime.Format("02.01.2006 15:04:05")

	// Prepare the response
	response := map[string]string{
		"currentDateTime": formattedTime,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// periodicRefresh sends a GET request to the /refresh endpoint every 12 minutes
func periodicRefresh() {
	for {
		startTime := time.Now()

		resp, err := http.Get("https://ion-bot.onrender.com/refresh")
		if err != nil {
			fmt.Println("Error sending periodic refresh request:", err)
			time.Sleep(3 * time.Minute) // Wait before retrying
			continue
		}

		// Ensure the response body is closed properly
		if resp.StatusCode != http.StatusOK {
			fmt.Println("Unexpected status code from refresh endpoint:", resp.StatusCode)
		}
		resp.Body.Close()

		// Reset the global variable
		mu.Lock()
		globalVar = 0
		mu.Unlock()

		fmt.Println("Global variable reset to 0")

		// Wait for the remainder of the 12 minutes
		elapsed := time.Since(startTime)
		time.Sleep(3*time.Minute - elapsed)
	}
}

func main() {
	http.HandleFunc("/sum", sumHandler)
	http.HandleFunc("/refresh", refreshHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if PORT is not set
	}

	go periodicRefresh() // Start the periodic refresh in a new goroutine

	fmt.Printf("Server is running on port %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
