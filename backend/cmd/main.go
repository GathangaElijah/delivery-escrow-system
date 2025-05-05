package main

import (
	"des/handlers"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
)

// openBrowser tries to open the URL in the default browser
func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		log.Printf("Error opening browser: %v", err)
	}
}

func main() {
	mux := http.NewServeMux()

	// Serve static files
	files_handler := http.FileServer(http.Dir("../static"))
	mux.Handle("/static/", http.StripPrefix("/static", files_handler))

	// Home page
	mux.Handle("/", handlers.Index())
	
	// Product detail page
	mux.Handle("/product/", handlers.GetProduct())
	
	// Authentication routes
	mux.Handle("/register", handlers.Register())
	mux.Handle("/register/buyer", handlers.RegisterBuyer())
	mux.Handle("/register/seller", handlers.RegisterSeller())
	mux.Handle("/register/transporter", handlers.RegisterTransporter())
	mux.Handle("/login", handlers.Login())

	// Define server address
	addr := ":8080"
	url := fmt.Sprintf("http://localhost%s", addr)

	// Print clickable link in terminal
	fmt.Printf("Server starting at \033[34m\033[4m%s\033[0m\n", url)
	
	// Start the server
	log.Printf("Starting server on %s", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
