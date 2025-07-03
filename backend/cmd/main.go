package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	multiplexer "des/backend/api/v1"
	middleware "des/backend/middleware"
)

func main() {
	mux := multiplexer.Multiplexer()
	// DEBUG
	handler := middleware.EnableCORS(mux)

	// Define server address
	port := os.Getenv("PORT")
	if port == "" {
		port = ":5001"
	}

	url := fmt.Sprintf("http://localhost%s", port)

	// Print clickable link in terminal
	fmt.Printf("Server starting at \033[34m\033[4m%s\033[0m\n", url)

	// Start the server
	log.Printf("Starting server on %s", port)
	err := http.ListenAndServe(port, handler)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
