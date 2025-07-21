package auth

import (
	"encoding/json"
	"log"
	"net/http"
)

// Define the structure to receive
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Define the structure to respond with
type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token,omitempty"` // optional
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler hit")
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON request
	var loginReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Login attempt: email=%s, password=%s", loginReq.Email, loginReq.Password)

	// Simulate user validation (you'll replace with DB logic)
	if loginReq.Email == "user" && loginReq.Password == "user" {
		resp := LoginResponse{
			Message: "Login successful",
			Token:   "mock-jwt-token-123", // Replace with real JWT later
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	} else {
		resp := LoginResponse{
			Message: "Invalid email or password",
		}
		// w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
	}
}
