package multiplexer

import (
	"des/backend/auth"
	"net/http"
)

// Here we are registering all incoming request patterns.
func Multiplexer() *http.ServeMux {
	mux := http.NewServeMux()

	// Serve static files
	files_handler := http.FileServer(http.Dir("./backend/static"))
	mux.Handle("/static/", http.StripPrefix("/static", files_handler))
	
	// Authentication routes
	mux.HandleFunc("/login", auth.LoginHandler)

	return mux
}
