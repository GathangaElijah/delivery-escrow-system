package handlers

import (
    "fmt"
    "html/template"
    "net/http"
    "path/filepath"
    "strings"
)

// ProductData holds data for the product template
type ProductData struct {
    Title   string
    Product Product
}

// GetProduct handles requests for individual product pages
func GetProduct() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Extract product ID from URL
        // Assuming URL pattern like /product/123
        pathParts := strings.Split(r.URL.Path, "/")
        if len(pathParts) < 3 {
            http.Error(w, "Invalid product URL", http.StatusBadRequest)
            return
        }
        
        productID := pathParts[2]
        
        // In a real app, I would fetch the product from a database
        // For now, am using a mock product based on the ID
        product := Product{
            ID:          productID,
            Name:        fmt.Sprintf("Product %s", productID),
            Description: fmt.Sprintf("Detailed description for product %s. This product is part of our premium collection.", productID),
            Price:       99.99 + (float64(len(productID)) * 10), // Just a mock price calculation
            ImageURL:    fmt.Sprintf("/static/img/product%s.jpg", productID),
        }
        
        data := ProductData{
            Title:   fmt.Sprintf("DES - %s", product.Name),
            Product: product,
        }
        
        // Parse template
        tmpl, err := template.ParseFiles(filepath.Join("backend/templates", "product.html"))
        if err != nil {
            http.Error(w, "Failed to load template", http.StatusInternalServerError)
            return
        }
        
        // Execute template
        if err := tmpl.Execute(w, data); err != nil {
            http.Error(w, "Failed to render template", http.StatusInternalServerError)
            return
        }
    }
}