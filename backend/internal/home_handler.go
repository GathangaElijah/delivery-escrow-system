package internal

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
			{ID: "1", Name: "Product 1", Description: "Description for product 1", Price: 99.99, ImageURL: "/static/img/product1.jpg"},
			{ID: "2", Name: "Product 2", Description: "Description for product 2", Price: 149.99, ImageURL: "/static/img/product2.jpg"},
			{ID: "3", Name: "Product 3", Description: "Description for product 3", Price: 199.99, ImageURL: "/static/img/product3.jpg"},
		}

		data := IndexData{
			Title:    "DES - Delivery Escrow System",
			Products: products,
		}

		// Parse template
		tmpl, err := template.ParseFiles("backend/templates/home.html")
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
