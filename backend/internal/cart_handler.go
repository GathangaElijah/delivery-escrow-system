package internal

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

// CartItem represents an item in the shopping cart
type CartItem struct {
	ProductID   string
	Name        string
	Description string
	Price       float64
	Quantity    int
	ImageURL    string
}

// CartData holds data for the cart template
type CartData struct {
	Title     string
	CartItems []CartItem
	Total     float64
}

// ViewCart handles the cart page
func ViewCart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get cart from session/cookie
		var cartItems []CartItem
		
		// Try to get cart cookie
		cookie, err := r.Cookie("cart")
		if err == nil {
			// Cookie exists, decode it
			encodedJSON := cookie.Value
			
			// Debug the cookie value
			fmt.Printf("Raw cookie value: %s\n", encodedJSON)
			
			// Decode from base64
			jsonData, err := base64.StdEncoding.DecodeString(encodedJSON)
			if err != nil {
				fmt.Printf("Error decoding base64: %v\n", err)
				// Try to decode directly as JSON as fallback
				err = json.Unmarshal([]byte(encodedJSON), &cartItems)
				if err != nil {
					fmt.Printf("Error unmarshaling cart cookie: %v\n", err)
					// If there's an error, start with an empty cart
					cartItems = []CartItem{}
				}
			} else {
				// Decode JSON
				err = json.Unmarshal(jsonData, &cartItems)
				if err != nil {
					fmt.Printf("Error unmarshaling cart cookie: %v\n", err)
					// If there's an error, start with an empty cart
					cartItems = []CartItem{}
				}
			}
		} else {
			fmt.Printf("No cart cookie found: %v\n", err)
		}
		
		// Debug log
		fmt.Printf("Cart has %d items\n", len(cartItems))
		if len(cartItems) > 0 {
			fmt.Printf("First item: %+v\n", cartItems[0])
		}
		
		// Calculate total
		var total float64
		for _, item := range cartItems {
			total += item.Price * float64(item.Quantity)
		}
		
		data := CartData{
			Title:     "Your Shopping Cart",
			CartItems: cartItems,
			Total:     total,
		}
		
		// Define template functions
		funcMap := template.FuncMap{
			"multiply": func(a, b float64) float64 {
				return a * b
			},
			"mul": func(a, b float64) float64 {
				return a * b
			},
			"float64": func(i int) float64 {
				return float64(i)
			},
		}
		
		// Try to find the template
		tmpl, err := template.New("cart.html").Funcs(funcMap).ParseFiles("backend/templates/cart.html")
		if err != nil {
			// Try another path
			tmpl, err = template.New("cart.html").Funcs(funcMap).ParseFiles("templates/cart.html")
			if err != nil {
				http.Error(w, "Failed to load template: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
		
		// Execute template
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to render template: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

// AddToCart handles adding a product to the cart
func AddToCart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse form data
		err := r.ParseForm()
		if err != nil {
			fmt.Printf("Error parsing form: %v\n", err)
			http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Get product details from form
		productID := r.FormValue("product_id")
		name := r.FormValue("name")
		description := r.FormValue("description")
		priceStr := r.FormValue("price")
		quantityStr := r.FormValue("quantity")
		imageURL := r.FormValue("image_url")

		// Validate required fields
		if productID == "" || name == "" || priceStr == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		// Parse price
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			http.Error(w, "Invalid price format", http.StatusBadRequest)
			return
		}

		// Parse quantity
		quantity := 1
		if quantityStr != "" {
			quantity, err = strconv.Atoi(quantityStr)
			if err != nil || quantity < 1 {
				quantity = 1
			}
		}

		// Create new cart item
		newItem := CartItem{
			ProductID:   productID,
			Name:        name,
			Description: description,
			Price:       price,
			Quantity:    quantity,
			ImageURL:    imageURL,
		}

		// Get existing cart
		var cartItems []CartItem

		// Try to get cart cookie
		cookie, err := r.Cookie("cart")
		if err == nil {
			// Cookie exists, decode it
			encodedJSON := cookie.Value
			
			// Decode from base64
			jsonData, err := base64.StdEncoding.DecodeString(encodedJSON)
			if err != nil {
				fmt.Printf("Error decoding base64: %v\n", err)
				// Try to decode directly as JSON as fallback
				err = json.Unmarshal([]byte(encodedJSON), &cartItems)
				if err != nil {
					fmt.Printf("Error unmarshaling cart cookie: %v\n", err)
					// If there's an error, start with an empty cart
					cartItems = []CartItem{}
				}
			} else {
				// Decode JSON
				err = json.Unmarshal(jsonData, &cartItems)
				if err != nil {
					fmt.Printf("Error unmarshaling cart cookie: %v\n", err)
					// If there's an error, start with an empty cart
					cartItems = []CartItem{}
				}
			}
		}

		// Check if product already in cart
		found := false
		for i, item := range cartItems {
			if item.ProductID == productID {
				// Update quantity
				cartItems[i].Quantity += quantity
				found = true
				break
			}
		}

		// If not found, add to cart
		if !found {
			cartItems = append(cartItems, newItem)
		}

		// Encode cart to JSON
		cartJSON, err := json.Marshal(cartItems)
		if err != nil {
			fmt.Printf("Error marshaling cart: %v\n", err)
			http.Error(w, "Failed to encode cart: "+err.Error(), http.StatusInternalServerError)
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

		// Redirect back to product page or cart
		redirectURL := r.FormValue("redirect")
		if redirectURL == "" {
			redirectURL = "/cart"
		}
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}

// UpdateCart handles updating cart quantities
func UpdateCart() http.HandlerFunc {
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

		// Get existing cart
		var cartItems []CartItem

		// Try to get cart cookie
		cookie, err := r.Cookie("cart")
		if err == nil {
			// Cookie exists, decode it
			encodedJSON := cookie.Value
			
			// Decode from base64
			jsonData, err := base64.StdEncoding.DecodeString(encodedJSON)
			if err != nil {
				fmt.Printf("Error decoding base64: %v\n", err)
				// Try to decode directly as JSON as fallback
				err = json.Unmarshal([]byte(encodedJSON), &cartItems)
				if err != nil {
					fmt.Printf("Error unmarshaling cart cookie: %v\n", err)
					cartItems = []CartItem{}
				}
			} else {
				// Decode JSON
				err = json.Unmarshal(jsonData, &cartItems)
				if err != nil {
					fmt.Printf("Error unmarshaling cart cookie: %v\n", err)
					cartItems = []CartItem{}
				}
			}
		}

		// Update quantities
		for i, item := range cartItems {
			quantityStr := r.FormValue(fmt.Sprintf("quantity_%s", item.ProductID))
			if quantityStr != "" {
				quantity, err := strconv.Atoi(quantityStr)
				if err == nil && quantity > 0 {
					cartItems[i].Quantity = quantity
				} else if quantity <= 0 {
					// Remove item if quantity is 0 or negative
					cartItems = append(cartItems[:i], cartItems[i+1:]...)
					i-- // Adjust index after removal
				}
			}
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

		// Redirect to cart page
		http.Redirect(w, r, "/cart", http.StatusSeeOther)
	}
}

// ClearCart handles removing all items from the cart
func ClearCart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Clear cart cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "cart",
			Value:    "",
			Path:     "/",
			MaxAge:   -1, // Delete cookie
			HttpOnly: true,
		})

		// Redirect to cart page
		http.Redirect(w, r, "/cart", http.StatusSeeOther)
	}
}
