package main

import (
	"des/handlers"
	
	"net/http"
)

func main() {

	mux := http.NewServeMux()

	files_handler := http.FileServer(http.Dir("../static"))
	mux.Handle("/static/", http.StripPrefix("/static", files_handler))

	mux.Handle("/", handlers.Index())
}
