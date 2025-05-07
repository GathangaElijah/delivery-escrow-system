package handlers

import (
	"html/template"
	"net/http"
)

// Product represents a product that can be purchased
type Product struct {
	ID          string
	Name        string
	Description string
	Price       float64
	ImageURL    string
}

// IndexData holds data for the index template
type IndexData struct {
	Title    string
	Products []Product
}

// Index handles the home page request
func Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Sample products - in a real app, these would come from a database
		products := []Product{
			{ID: "1", Name: "Eton Pauls Shirt", Description: "The most expensive shirt", Price: 4000.00, ImageURL: "/static/images/Eton-Pauls-Shirt.jpg"},
			{ID: "2", Name: "Rolex Watch", Description: "Expensive watch", Price: 10000.99, ImageURL: "/static/images/expensive-watch.jpeg"},
			{ID: "3", Name: "Nike Air Max", Description: "Expensive shoe", Price: 1999.99, ImageURL: "/static/images/original-nike-air.jpg"},
		}

		data := IndexData{
			Title:    "DES - Delivery Escrow System",
			Products: products,
		}

		// Parse template
		tmpl, err := template.ParseFiles("templates/home.html")
		if err != nil {
			http.Error(w, "Failed to load template", http.StatusInternalServerError)
			return
		}

		// Execute template
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
		}
	}
}
