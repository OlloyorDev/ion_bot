package main

import (
	"bytes"
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
	fmt.Printf("Sum calculated: %d\n", sum)

	// Prepare the response
	resp := Response{Sum: sum}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

	// Send POST request to FCM
	sendNotification(sum)
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
