package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"
)

// DebugInfo provides information about the application environment
func DebugInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get current working directory
		cwd, err := os.Getwd()
		if err != nil {
			http.Error(w, "Failed to get current directory: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if template directories exist
		templateDirs := []string{
			"templates",
			"backend/templates",
			"../templates",
		}

		var output strings.Builder
		output.WriteString(fmt.Sprintf("Current working directory: %s\n\n", cwd))

		// Check template directories
		output.WriteString("Template directories:\n")
		for _, dir := range templateDirs {
			if info, err := os.Stat(dir); err == nil && info.IsDir() {
				output.WriteString(fmt.Sprintf("✅ %s (exists)\n", dir))

				// List files in directory
				files, err := os.ReadDir(dir)
				if err == nil {
					output.WriteString("   Files:\n")
					for _, file := range files {
						output.WriteString(fmt.Sprintf("   - %s\n", file.Name()))
					}
				}
			} else {
				output.WriteString(fmt.Sprintf("❌ %s (not found)\n", dir))
			}
		}

		// Check specific template paths
		templatePaths := []string{
			"backend/templates/cart.html",
			"templates/cart.html",
			"../templates/cart.html",
		}

		output.WriteString("\nTemplate files:\n")
		for _, path := range templatePaths {
			if _, err := os.Stat(path); err == nil {
				output.WriteString(fmt.Sprintf("✅ %s (exists)\n", path))
			} else {
				output.WriteString(fmt.Sprintf("❌ %s (not found: %v)\n", path, err))
			}
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(output.String()))
	}
}

// DebugCart displays the current cart contents
func DebugCart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Try to get cart cookie
		cookie, err := r.Cookie("cart")

		var output strings.Builder
		output.WriteString("Cart Debug Information:\n\n")

		if err != nil {
			output.WriteString(fmt.Sprintf("No cart cookie found: %v\n", err))
		} else {
			output.WriteString(fmt.Sprintf("Cart cookie value: %s\n\n", cookie.Value))

			// Try to parse the cart
			var cartItems []CartItem

			// Try to decode from base64
			jsonData, err := base64.StdEncoding.DecodeString(cookie.Value)
			if err != nil {
				output.WriteString(fmt.Sprintf("Error decoding base64: %v\n", err))
				// Try to decode directly as JSON as fallback
				err = json.Unmarshal([]byte(cookie.Value), &cartItems)
				if err != nil {
					output.WriteString(fmt.Sprintf("Error parsing cart JSON: %v\n", err))
				}
			} else {
				// Decode JSON
				err = json.Unmarshal(jsonData, &cartItems)
				if err != nil {
					output.WriteString(fmt.Sprintf("Error parsing cart JSON: %v\n", err))
				}
			}

			if err == nil {
				output.WriteString(fmt.Sprintf("Cart contains %d items:\n", len(cartItems)))

				for i, item := range cartItems {
					output.WriteString(fmt.Sprintf("\nItem %d:\n", i+1))
					output.WriteString(fmt.Sprintf("  ProductID: %s\n", item.ProductID))
					output.WriteString(fmt.Sprintf("  Name: %s\n", item.Name))
					output.WriteString(fmt.Sprintf("  Price: $%.2f\n", item.Price))
					output.WriteString(fmt.Sprintf("  Quantity: %d\n", item.Quantity))
					output.WriteString(fmt.Sprintf("  ImageURL: %s\n", item.ImageURL))
				}
			}
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(output.String()))
	}
}

// TestCart displays a test page for cart functionality
func TestCart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse template
		tmpl, err := template.ParseFiles("backend/templates/test_cart.html")
		if err != nil {
			// Try another path
			tmpl, err = template.ParseFiles("templates/test_cart.html")
			if err != nil {
				http.Error(w, "Failed to load template: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Execute template
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Failed to render template: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

// CreateTestOrder creates a test order for demonstration purposes
func CreateTestOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a test cart with some items
		cartItems := []CartItem{
			{
				ProductID:   "TEST-001",
				Name:        "Test Product 1",
				Description: "This is a test product for demonstration",
				Price:       99.99,
				Quantity:    2,
				ImageURL:    "/static/img/product1.jpg",
			},
			{
				ProductID:   "TEST-002",
				Name:        "Test Product 2",
				Description: "Another test product for demonstration",
				Price:       149.99,
				Quantity:    1,
				ImageURL:    "/static/img/product2.jpg",
			},
		}

		// Encode cart to JSON
		cartJSON, err := json.Marshal(cartItems)
		if err != nil {
			http.Error(w, "Failed to encode cart", http.StatusInternalServerError)
			return
		}

		// Use base64 encoding
		encodedJSON := base64.StdEncoding.EncodeToString(cartJSON)

		// Set cart cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "cart",
			Value:    encodedJSON,
			Path:     "/",
			MaxAge:   86400, // 1 day
			HttpOnly: true,
		})

		// Redirect to seller dashboard
		http.Redirect(w, r, "/seller/dashboard", http.StatusSeeOther)
	}
}

// AddProduct2ToCart adds Product 2 to the cart for testing
func AddProduct2ToCart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a cart with Product 2
		cartItems := []CartItem{
			{
				ProductID:   "P002",
				Name:        "Product 2",
				Description: "This is Product 2 for demonstration",
				Price:       149.99,
				Quantity:    1,
				ImageURL:    "/static/img/product2.jpg",
			},
		}

		// Encode cart to JSON
		cartJSON, err := json.Marshal(cartItems)
		if err != nil {
			http.Error(w, "Failed to encode cart", http.StatusInternalServerError)
			return
		}

		// Use base64 encoding
		encodedJSON := base64.StdEncoding.EncodeToString(cartJSON)

		// Set cart cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "cart",
			Value:    encodedJSON,
			Path:     "/",
			MaxAge:   86400, // 1 day
			HttpOnly: true,
		})

		// Create an order from this cart
		orderItems := make([]OrderItem, 0)
		for _, item := range cartItems {
			orderItems = append(orderItems, OrderItem{
				ID:       item.ProductID,
				Name:     item.Name,
				Price:    item.Price,
				Quantity: item.Quantity,
				ImageURL: item.ImageURL,
			})
		}

		// Calculate total
		var total float64
		for _, item := range orderItems {
			total += item.Price * float64(item.Quantity)
		}

		// Create a new order
		orders := []Order{
			{
				ID:     fmt.Sprintf("ORD-P2-%d", time.Now().Unix()),
				Date:   time.Now().Format("2006-01-02"),
				Status: "pending",
				Items:  orderItems,
				Customer: User{
					ID:              "B54321",
					Name:            "John Doe",
					Email:           "john@example.com",
					ShippingAddress: "123 Main St, Anytown, USA",
				},
				EscrowStatus: "Funds Staked",
				EscrowAmount: total,
			},
		}

		// Save this order to a cookie for persistence
		orderJSON, err := json.Marshal(orders)
		if err == nil {
			encodedOrderJSON := base64.StdEncoding.EncodeToString(orderJSON)
			http.SetCookie(w, &http.Cookie{
				Name:     "orders",
				Value:    encodedOrderJSON,
				Path:     "/",
				MaxAge:   86400, // 1 day
				HttpOnly: true,
			})
		}

		// Redirect to seller dashboard
		http.Redirect(w, r, "/seller/dashboard", http.StatusSeeOther)
	}
}
