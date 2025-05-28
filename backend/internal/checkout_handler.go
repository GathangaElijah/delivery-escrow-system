package internal

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

// CheckoutData holds data for the checkout template
type CheckoutData struct {
	Title     string
	Product   Product
	CartItems []CartItem
	Total     float64
	User      User
}

// Checkout handles the checkout page
func Checkout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract path parts
		pathParts := strings.Split(r.URL.Path, "/")

		// Mock user data - in a real app, get from session
		user := User{
			ID:     "B54321",
			Name:   "John Doe",
			Email:  "buyer@example.com",
			Wallet: "5FHneW46xGXgs5mUiveU4sbTyGBzmstUspZC92UhjJM694ty",
			Role:   "buyer",
		}

		// Check if we're checking out a single product or the cart
		if len(pathParts) > 2 && pathParts[2] != "cart" && pathParts[2] != "confirmation" {
			// Single product checkout
			productID := pathParts[2]

			// Mock product data - in a real app, get from database
			product := Product{
				ID:          productID,
				Name:        fmt.Sprintf("Product %s", productID),
				Description: fmt.Sprintf("Description for product %s", productID),
				Price:       99.99,
				ImageURL:    fmt.Sprintf("/static/img/product%s.jpg", productID),
			}

			data := CheckoutData{
				Title:   "Checkout",
				Product: product,
				User:    user,
			}

			tmpl, err := template.ParseFiles("backend/templates/checkout.html")
			if err != nil {
				http.Error(w, "Failed to load template", http.StatusInternalServerError)
				return
			}

			tmpl.Execute(w, data)
			return
		}

		// Cart checkout
		var cartItems []CartItem

		// Try to get cart cookie
		cookie, err := r.Cookie("cart")
		if err == nil {
			// Cookie exists, decode it
			json.Unmarshal([]byte(cookie.Value), &cartItems)
		}

		// Calculate total
		var total float64
		for _, item := range cartItems {
			total += item.Price * float64(item.Quantity)
		}

		data := CheckoutData{
			Title:     "Checkout",
			CartItems: cartItems,
			Total:     total,
			User:      user,
		}

		tmpl, err := template.ParseFiles("backend/templates/checkout_cart.html")
		if err != nil {
			http.Error(w, "Failed to load template", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data)
	}
}

// ProcessCheckout handles the checkout process
func ProcessCheckout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse form
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Get cart items from cookies
		var cartItems []CartItem
		cookie, err := r.Cookie("cart")
		if err == nil {
			// Cookie exists, decode it
			encodedJSON := cookie.Value

			// Decode from base64
			jsonData, err := base64.StdEncoding.DecodeString(encodedJSON)
			if err != nil {
				log.Printf("Error decoding base64: %v", err)
				// Try to decode directly as JSON as fallback
				err = json.Unmarshal([]byte(encodedJSON), &cartItems)
				if err != nil {
					log.Printf("Error unmarshaling cart cookie: %v", err)
				}
			} else {
				// Decode JSON
				err = json.Unmarshal(jsonData, &cartItems)
				if err != nil {
					log.Printf("Error unmarshaling cart cookie: %v", err)
				}
			}
		}

		// If no cart items, return error
		if len(cartItems) == 0 {
			http.Error(w, "Cart is empty", http.StatusBadRequest)
			return
		}

		// Get form values
		buyerWallet := r.FormValue("buyer_wallet")
		buyerName := r.FormValue("buyer_name")
		buyerEmail := r.FormValue("buyer_email")
		shippingAddress := r.FormValue("shipping_address")

		// Validate required fields
		if buyerWallet == "" {
			buyerWallet = "5FHneW46xGXgs5mUiveU4sbTyGBzmstUspZC92UhjJM694ty" // Default wallet for testing
		}
		if buyerName == "" {
			buyerName = "John Doe" // Default name for testing
		}
		if buyerEmail == "" {
			buyerEmail = "john@example.com" // Default email for testing
		}
		if shippingAddress == "" {
			shippingAddress = "123 Main St, Anytown, USA" // Default address for testing
		}

		// Create order items from cart items
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
		order := Order{
			ID:     fmt.Sprintf("ORD-CHECKOUT-%d", time.Now().Unix()),
			Date:   time.Now().Format("2006-01-02"),
			Status: "pending",
			Items:  orderItems,
			Customer: User{
				ID:              "B" + fmt.Sprintf("%d", time.Now().Unix())[:5],
				Name:            buyerName,
				Email:           buyerEmail,
				ShippingAddress: shippingAddress,
			},
			EscrowStatus: "Funds Staked",
			EscrowAmount: total,
		}

		// Get existing orders from cookie
		var orders []Order
		ordersCookie, err := r.Cookie("orders")
		if err == nil {
			// Cookie exists, decode it
			encodedJSON := ordersCookie.Value

			// Decode from base64
			jsonData, err := base64.StdEncoding.DecodeString(encodedJSON)
			if err != nil {
				log.Printf("Error decoding base64: %v", err)
			} else {
				// Decode JSON
				err = json.Unmarshal(jsonData, &orders)
				if err != nil {
					log.Printf("Error unmarshaling orders cookie: %v", err)
				}
			}
		}

		// Add new order to orders
		orders = append(orders, order)

		// Save orders to cookie
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

		// Log the action
		log.Printf("Order created: %s", order.ID)
		log.Printf("Total amount: $%.2f", total)
		log.Printf("Buyer wallet: %s", buyerWallet)

		// Clear cart after successful checkout
		http.SetCookie(w, &http.Cookie{
			Name:     "cart",
			Value:    "",
			Path:     "/",
			MaxAge:   -1, // Delete cookie
			HttpOnly: true,
		})

		// Redirect to confirmation page
		http.Redirect(w, r, "/checkout/confirmation", http.StatusSeeOther)
	}
}

// CheckoutConfirmation handles the checkout confirmation page
func CheckoutConfirmation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("backend/templates/checkout_confirmation.html")
		if err != nil {
			http.Error(w, "Failed to load template", http.StatusInternalServerError)
			return
		}

		data := struct {
			Title string
		}{
			Title: "Checkout Confirmation",
		}

		tmpl.Execute(w, data)
	}
}
