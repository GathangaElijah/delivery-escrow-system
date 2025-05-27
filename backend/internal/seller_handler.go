package handlers

import (
	"crypto/sha256"
	"des/backend/blockchain"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// SellerDashboardData holds data for the seller dashboard template
type SellerDashboardData struct {
	Title            string
	User             User
	Orders           []Order
	ProcessingOrders []Order
	Products         []Product
}

// Order represents a customer order
type Order struct {
	ID           string
	Date         string
	Status       string
	Items        []OrderItem
	Customer     User
	EscrowStatus string
	EscrowAmount float64
}

// OrderItem represents a product in an order
type OrderItem struct {
	ID       string
	Name     string
	Price    float64
	Quantity int
	ImageURL string
}

// ShipmentData represents data for a shipment
type ShipmentData struct {
	OrderID             string
	PackagePhotoPath    string
	ManufacturerSigPath string
	TrackingNumber      string
	Notes               string
	QRCodeHash          string
	QRCodeURL           string
	CreatedAt           time.Time
	ScanCount           int
	MaxScans            int
	SellerID            string
	BuyerID             string
}

// SellerDashboard handles the seller dashboard page
func SellerDashboard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// In a real app, you would:
		// 1. Check if the user is logged in
		// 2. Check if the user is a seller
		// 3. Get the seller's information from the database

		// Mock seller data for demonstration
		seller := User{
			ID:           "S12345",
			BusinessName: "Example Store",
			Email:        "seller@example.com",
			Wallet:       "5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY",
			Role:         "seller",
		}

		// Get cart items from cookies - this would normally come from a database
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

		// Initialize orders slice
		var orders []Order

		// Convert cart items to orders if there are any
		if len(cartItems) > 0 {
			// Create an order from cart items
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
			orders = append(orders, Order{
				ID:     fmt.Sprintf("ORD-%d", time.Now().Unix()),
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
			})
		}

		// Check for orders in session
		session, err := r.Cookie("orders")
		if err == nil {
			// Cookie exists, decode it
			encodedJSON := session.Value

			var sessionOrders []Order
			// Decode from base64
			jsonData, err := base64.StdEncoding.DecodeString(encodedJSON)
			if err != nil {
				log.Printf("Error decoding base64: %v", err)
			} else {
				// Decode JSON
				err = json.Unmarshal(jsonData, &sessionOrders)
				if err != nil {
					log.Printf("Error unmarshaling orders cookie: %v", err)
				} else {
					// Add session orders to orders slice
					orders = append(orders, sessionOrders...)
				}
			}
		}

		// If still no orders, create a test order with Product 2
		if len(orders) == 0 {
			orders = []Order{
				{
					ID:     fmt.Sprintf("ORD-TEST-%d", time.Now().Unix()),
					Date:   time.Now().Format("2006-01-02"),
					Status: "pending",
					Items: []OrderItem{
						{
							ID:       "P002",
							Name:     "Product 2",
							Price:    149.99,
							Quantity: 1,
							ImageURL: "/static/img/product2.jpg",
						},
					},
					Customer: User{
						ID:              "B54321",
						Name:            "John Doe",
						Email:           "john@example.com",
						ShippingAddress: "123 Main St, Anytown, USA",
					},
					EscrowStatus: "Funds Staked",
					EscrowAmount: 149.99,
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
		}

		// Filter orders for processing
		var processingOrders []Order
		for _, order := range orders {
			if order.Status == "processing" {
				processingOrders = append(processingOrders, order)
			}
		}

		// Get products from orders
		var products []Product
		productMap := make(map[string]Product)

		// First add products from cart items
		for _, item := range cartItems {
			productMap[item.ProductID] = Product{
				ID:          item.ProductID,
				Name:        item.Name,
				Description: item.Description,
				Price:       item.Price,
				ImageURL:    item.ImageURL,
			}
		}

		// Then add products from orders
		for _, order := range orders {
			for _, item := range order.Items {
				if _, exists := productMap[item.ID]; !exists {
					productMap[item.ID] = Product{
						ID:          item.ID,
						Name:        item.Name,
						Description: "Product from order",
						Price:       item.Price,
						ImageURL:    item.ImageURL,
					}
				}
			}
		}

		// Convert map to slice
		for _, product := range productMap {
			products = append(products, product)
		}

		// Prepare data for the template
		data := SellerDashboardData{
			Title:            "DES - Seller Dashboard",
			User:             seller,
			Orders:           orders,
			ProcessingOrders: processingOrders,
			Products:         products,
		}

		// Parse template
		tmpl, err := template.ParseFiles("backend/templates/seller_dashboard.html")
		if err != nil {
			// Try alternative path
			tmpl, err = template.ParseFiles("backend/templates/seller_dashboard.html")
			if err != nil {
				http.Error(w, "Failed to load template: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Execute template
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Failed to render template: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// AddProduct handles adding a new product
func AddProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse multipart form data (for file uploads)
		if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max memory
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Get form values
		productName := r.FormValue("product_name")
		productDescription := r.FormValue("product_description")
		productPrice := r.FormValue("product_price")
		productCategory := r.FormValue("product_category")
		productQuantity := r.FormValue("product_quantity")
		productSKU := r.FormValue("product_sku")

		// Convert price to float
		price, err := strconv.ParseFloat(productPrice, 64)
		if err != nil {
			http.Error(w, "Invalid price format", http.StatusBadRequest)
			return
		}

		// Convert quantity to int
		quantity, err := strconv.Atoi(productQuantity)
		if err != nil {
			http.Error(w, "Invalid quantity format", http.StatusBadRequest)
			return
		}

		// Handle product image upload
		productImageFile, productImageHeader, err := r.FormFile("product_image")
		if err != nil {
			http.Error(w, "Product image is required", http.StatusBadRequest)
			return
		}
		defer productImageFile.Close()

		// In a real app, you would:
		// 1. Save the image to a storage system
		// 2. Get the URL of the saved image
		// 3. Store the product in a database

		// For now, just log the product information
		log.Printf("New product added: %s", productName)
		log.Printf("Description: %s", productDescription)
		log.Printf("Price: $%.2f", price)
		log.Printf("Category: %s", productCategory)
		log.Printf("Quantity: %d", quantity)
		log.Printf("SKU: %s", productSKU)
		log.Printf("Image: %s", productImageHeader.Filename)

		// Check if manufacturer signature was uploaded
		manufacturerSigFile, manufacturerSigHeader, err := r.FormFile("manufacturer_signature")
		if err == nil {
			defer manufacturerSigFile.Close()
			log.Printf("Manufacturer signature: %s", manufacturerSigHeader.Filename)
		}

		// Redirect back to the seller dashboard
		http.Redirect(w, r, "/seller/dashboard", http.StatusSeeOther)
	}
}

// PrepareShipment handles preparing a shipment
func PrepareShipment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse multipart form
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Get form values
		orderID := r.FormValue("order_id")
		trackingNumber := r.FormValue("tracking_number")
		shippingNotes := r.FormValue("shipping_notes")

		// Validate required fields
		if orderID == "" {
			http.Error(w, "Order ID is required", http.StatusBadRequest)
			return
		}

		// Handle package photo upload
		packagePhotoFile, packagePhotoHeader, err := r.FormFile("package_photo")
		if err != nil {
			http.Error(w, "Package photo is required", http.StatusBadRequest)
			return
		}
		defer packagePhotoFile.Close()

		// Handle manufacturer signature upload (optional)
		var manufacturerSigHeader *multipart.FileHeader
		manufacturerSigFile, manufacturerSigHeader, err := r.FormFile("manufacturer_signature")
		if err == nil {
			defer manufacturerSigFile.Close()
		}

		// In a real app, you would:
		// 1. Save the uploaded files to a storage system
		// 2. Generate a QR code for the shipment
		// 3. Update the order status in the database

		// For now, just mock the QR code generation
		// Generate a unique hash for the QR code
		h := sha256.New()
		io.WriteString(h, fmt.Sprintf("%s-%d", orderID, time.Now().UnixNano()))
		qrCodeHash := hex.EncodeToString(h.Sum(nil))

		// Mock QR code URL
		qrCodeURL := fmt.Sprintf("/static/img/qr/%s.png", qrCodeHash)

		// Create a shipment record
		shipment := ShipmentData{
			OrderID:             orderID,
			PackagePhotoPath:    "/static/img/packages/" + packagePhotoHeader.Filename,
			ManufacturerSigPath: "/static/img/signatures/" + manufacturerSigHeader.Filename,
			TrackingNumber:      trackingNumber,
			Notes:               shippingNotes,
			QRCodeHash:          qrCodeHash,
			QRCodeURL:           qrCodeURL,
			CreatedAt:           time.Now(),
			ScanCount:           0,
			MaxScans:            2,        // Can only be scanned twice
			SellerID:            "S12345", // In a real app, this would be the logged-in seller's ID
			BuyerID:             "B54321", // In a real app, this would be the buyer's ID from the order
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

		// Find the order and update its status
		for i, order := range orders {
			if order.ID == orderID {
				orders[i].Status = "shipped"
				log.Printf("Order %s status updated to shipped", orderID)
				break
			}
		}

		// Save updated orders back to cookie
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

		// For now, just log the shipment information
		log.Printf("New shipment prepared for order: %s", shipment.OrderID)
		log.Printf("QR Code Hash: %s", shipment.QRCodeHash)
		log.Printf("Created at: %s", shipment.CreatedAt.Format(time.RFC3339))

		// In a real app, you would store the shipment in a database

		// Return the QR code URL as JSON
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{"qrCodeUrl": "%s"}`, qrCodeURL)))
	}
}

// ProcessOrder handles processing an order
func ProcessOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract order ID from URL
		// Assuming URL pattern like /seller/orders/123/process
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 5 {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		orderID := pathParts[3]
		log.Printf("Processing order: %s", orderID)

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

		// Find the order and update its status
		orderFound := false
		for i, order := range orders {
			if order.ID == orderID {
				orders[i].Status = "processing"
				orderFound = true
				log.Printf("Order %s status updated to processing", orderID)
				break
			}
		}

		if !orderFound {
			log.Printf("Order %s not found", orderID)
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}

		// Save updated orders back to cookie
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

		// In a real app, you would:
		// 1. Update the order status in the database
		// 2. Notify the buyer that their order is being processed
		// 3. Add the order to the seller's processing queue

		// Redirect back to the seller dashboard
		http.Redirect(w, r, "/seller/dashboard#orders", http.StatusSeeOther)
	}
}

// ShipOrder handles shipping an order
func ShipOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract order ID from URL
		// Assuming URL pattern like /seller/orders/123/ship
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 5 {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		orderID := pathParts[3]
		log.Printf("Shipping order: %s", orderID)

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

		// Find the order and update its status
		orderFound := false
		for i, order := range orders {
			if order.ID == orderID {
				orders[i].Status = "shipped"
				orderFound = true
				log.Printf("Order %s status updated to shipped", orderID)
				break
			}
		}

		if !orderFound {
			log.Printf("Order %s not found", orderID)
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}

		// Save updated orders back to cookie
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

		// In a real app, you would:
		// 1. Update the order status in the database
		// 2. Notify the buyer that their order has been shipped
		// 3. Generate tracking information

		// Redirect to the shipping preparation page
		http.Redirect(w, r, "/seller/dashboard#shipping", http.StatusSeeOther)
	}
}

// UpdateAccount handles updating seller account information
func UpdateAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Get form values
		businessName := r.FormValue("business_name")
		email := r.FormValue("email")
		wallet := r.FormValue("wallet")
		// currentPassword := r.FormValue("current_password")
		newPassword := r.FormValue("new_password")
		confirmPassword := r.FormValue("confirm_password")

		// Validate form data
		if newPassword != "" && newPassword != confirmPassword {
			http.Error(w, "New passwords do not match", http.StatusBadRequest)
			return
		}

		// In a real app, you would:
		// 1. Check if the user is logged in
		// 2. Check if the current password is correct
		// 3. Update the user information in the database
		// 4. If a new password was provided, hash it and update it in the database

		// For now, just log the action
		log.Printf("Account updated for business: %s", businessName)
		log.Printf("New email: %s", email)
		log.Printf("New wallet: %s", wallet)
		if newPassword != "" {
			log.Printf("Password changed")
		}

		// Redirect back to the seller dashboard
		http.Redirect(w, r, "/seller/dashboard#account", http.StatusSeeOther)
	}
}

// ScanQRCode handles scanning a QR code for delivery confirmation
func ScanQRCode() http.HandlerFunc {
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

		// Get QR code hash from form
		qrCodeHash := r.FormValue("qr_code_hash")
		userRole := r.FormValue("user_role") // "buyer" or "seller" or "transporter"
		userWallet := r.FormValue("user_wallet")

		log.Printf("QR code scan request received for hash: %s", qrCodeHash)
		log.Printf("User role: %s", userRole)
		log.Printf("User wallet: %s", userWallet)

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

		// Create contract client
		contractClient := blockchain.NewContractClient()

		// Find the order associated with this QR code and update its status
		orderFound := false
		var orderID string

		for i, order := range orders {
			if order.Status == "shipped" {
				// In a real implementation, we would check if the QR code hash matches the order
				orderID = order.ID

				if userRole == "buyer" {
					// When the buyer scans the QR code, mark as delivered and trigger smart contract
					orders[i].Status = "delivered"
					orders[i].EscrowStatus = "Delivery Confirmed"

					// Call the smart contract's confirm_delivery function
					err := contractClient.ConfirmDelivery(orderID, userWallet)
					if err != nil {
						log.Printf("Error confirming delivery: %v", err)
					} else {
						log.Printf("Delivery confirmed for order %s", orderID)
						orders[i].EscrowStatus = "Funds Released by Smart Contract"
					}
				} else if userRole == "transporter" {
					// When the transporter scans, they're submitting proof of delivery
					// Call the submit_proof function on the smart contract
					err := contractClient.SubmitProof(orderID, qrCodeHash, userWallet)
					if err != nil {
						log.Printf("Error submitting proof: %v", err)
					} else {
						log.Printf("Proof submitted for order %s", orderID)
						orders[i].EscrowStatus = "Proof of Delivery Submitted"
					}
				}

				orderFound = true
				break
			}
		}

		if !orderFound {
			log.Printf("No shipped orders found to update")
		}

		// Save updated orders back to cookie
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

		// For now, just log the action
		log.Printf("QR code scanned: %s", qrCodeHash)

		var responseMessage string
		if userRole == "buyer" {
			responseMessage = "Delivery confirmed. Funds will be automatically released by the smart contract."
		} else if userRole == "seller" {
			responseMessage = "Scan recorded. Waiting for buyer confirmation."
		} else if userRole == "transporter" {
			responseMessage = "Proof of delivery submitted. Waiting for buyer confirmation."
		} else {
			responseMessage = "Scan recorded."
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{"success": true, "message": "%s"}`, responseMessage)))
	}
}

// ReleaseFunds handles the buyer releasing funds from escrow
func ReleaseFunds() http.HandlerFunc {
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

		// Get form values
		orderID := r.FormValue("order_id")
		buyerWallet := r.FormValue("buyer_wallet")

		log.Printf("Fund release request for order: %s", orderID)
		log.Printf("Buyer wallet: %s", buyerWallet)

		// In a real app, you would:
		// 1. Check if the user is logged in and is the buyer for this order
		// 2. Get the escrow contract address from the database
		// 3. Call the release function on the escrow contract
		// 4. Update the order status in the database

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

		// Find the order and update its status
		for i, order := range orders {
			if order.ID == orderID {
				orders[i].EscrowStatus = "Funds Released"
				log.Printf("Order %s escrow status updated to Funds Released", orderID)
				break
			}
		}

		// Save updated orders back to cookie
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

		// For now, just log the action
		log.Printf("Funds released for order: %s", orderID)

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success": true, "message": "Funds successfully released to seller"}`))
	}
}

// RaiseDispute handles raising a dispute on the escrow contract
func RaiseDispute() http.HandlerFunc {
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

		// Get form values
		orderID := r.FormValue("order_id")
		buyerWallet := r.FormValue("buyer_wallet")
		disputeReason := r.FormValue("dispute_reason")

		log.Printf("Dispute raised for order: %s", orderID)
		log.Printf("Buyer wallet: %s", buyerWallet)
		log.Printf("Dispute reason: %s", disputeReason)

		// Create contract client
		contractClient := blockchain.NewContractClient()

		// Raise dispute on contract
		err := contractClient.RaiseDispute(orderID, buyerWallet, disputeReason)
		if err != nil {
			log.Printf("Error raising dispute: %v", err)
			http.Error(w, "Failed to raise dispute", http.StatusInternalServerError)
			return
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

		// Find the order and update its status
		for i, order := range orders {
			if order.ID == orderID {
				orders[i].EscrowStatus = "Dispute Raised"
				log.Printf("Order %s escrow status updated to Dispute Raised", orderID)
				break
			}
		}

		// Save updated orders back to cookie
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

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success": true, "message": "Dispute recorded. Our team will review your case."}`))
	}
}

// GetEscrowStatus handles checking the status of the escrow contract
func GetEscrowStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract order ID from URL
		// Assuming URL pattern like /escrow/status/123
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		orderID := pathParts[3]
		log.Printf("Checking escrow status for order: %s", orderID)

		// Create contract client
		contractClient := blockchain.NewContractClient()

		// Get contract status
		contractStatus, err := contractClient.GetContractStatus(orderID)
		if err != nil {
			log.Printf("Error getting contract status: %v", err)
			http.Error(w, "Failed to get contract status", http.StatusInternalServerError)
			return
		}

		// Get existing orders from cookie for additional info
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

		// Find the order for additional info
		var escrowStatus string
		var escrowAmount float64

		for _, order := range orders {
			if order.ID == orderID {
				escrowStatus = order.EscrowStatus
				escrowAmount = order.EscrowAmount
				break
			}
		}

		// If we couldn't find the order in the cookie, use contract status
		if escrowStatus == "" {
			if contractStatus.IsDelivered {
				escrowStatus = "Funds Released"
			} else if contractStatus.DisputeRaised {
				escrowStatus = "Dispute Raised"
			} else if contractStatus.ProofOfDelivery != "" {
				escrowStatus = "Proof of Delivery Submitted"
			} else {
				escrowStatus = "Funds in Escrow"
			}

			escrowAmount = contractStatus.Balance
		}

		// Return the status as JSON
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{
			"orderID": "%s",
			"escrowStatus": "%s",
			"escrowAmount": %.2f,
			"isDelivered": %t,
			"buyer": "%s",
			"seller": "%s",
			"transporter": "%s",
			"disputeRaised": %t,
			"message": "Smart contract automatically handles fund release upon delivery confirmation"
		}`, orderID, escrowStatus, escrowAmount, contractStatus.IsDelivered,
			contractStatus.Buyer, contractStatus.Seller, contractStatus.Transporter,
			contractStatus.DisputeRaised)))
	}
}
