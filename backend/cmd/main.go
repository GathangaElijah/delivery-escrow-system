package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"

	"des/backend/multiplexer"
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
	mux := multiplexer.Multiplexer()

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
