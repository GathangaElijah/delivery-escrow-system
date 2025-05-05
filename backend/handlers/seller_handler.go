package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
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
	ID       string
	Date     string
	Status   string
	Items    []OrderItem
	Customer User
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

		// Mock orders data
		orders := []Order{
			{
				ID:     "ORD-001",
				Date:   "2023-05-15",
				Status: "pending",
				Items: []OrderItem{
					{
						ID:       "PROD-001",
						Name:     "Smartphone",
						Price:    499.99,
						Quantity: 1,
						ImageURL: "/static/img/product1.jpg",
					},
				},
				Customer: User{
					ID:              "B12345",
					Name:            "John Doe",
					Email:           "john@example.com",
					ShippingAddress: "123 Main St, Anytown, AN 12345",
				},
			},
			{
				ID:     "ORD-002",
				Date:   "2023-05-14",
				Status: "processing",
				Items: []OrderItem{
					{
						ID:       "PROD-002",
						Name:     "Laptop",
						Price:    899.99,
						Quantity: 1,
						ImageURL: "/static/img/product2.jpg",
					},
					{
						ID:       "PROD-003",
						Name:     "Headphones",
						Price:    79.99,
						Quantity: 2,
						ImageURL: "/static/img/product3.jpg",
					},
				},
				Customer: User{
					ID:              "B54321",
					Name:            "Jane Smith",
					Email:           "jane@example.com",
					ShippingAddress: "456 Oak Ave, Othertown, OT 54321",
				},
			},
			{
				ID:     "ORD-003",
				Date:   "2023-05-13",
				Status: "shipped",
				Items: []OrderItem{
					{
						ID:       "PROD-004",
						Name:     "Tablet",
						Price:    349.99,
						Quantity: 1,
						ImageURL: "/static/img/product4.jpg",
					},
				},
				Customer: User{
					ID:              "B67890",
					Name:            "Bob Johnson",
					Email:           "bob@example.com",
					ShippingAddress: "789 Pine Rd, Somewhere, SW 67890",
				},
			},
		}

		// Filter processing orders for the shipping tab
		var processingOrders []Order
		for _, order := range orders {
			if order.Status == "processing" {
				processingOrders = append(processingOrders, order)
			}
		}

		// Mock products data
		products := []Product{
			{
				ID:          "PROD-001",
				Name:        "Smartphone",
				Description: "Latest model smartphone with high-resolution camera",
				Price:       499.99,
				ImageURL:    "/static/img/product1.jpg",
			},
			{
				ID:          "PROD-002",
				Name:        "Laptop",
				Description: "Powerful laptop for work and gaming",
				Price:       899.99,
				ImageURL:    "/static/img/product2.jpg",
			},
			{
				ID:          "PROD-003",
				Name:        "Headphones",
				Description: "Noise-cancelling wireless headphones",
				Price:       79.99,
				ImageURL:    "/static/img/product3.jpg",
			},
			{
				ID:          "PROD-004",
				Name:        "Tablet",
				Description: "Lightweight tablet with long battery life",
				Price:       349.99,
				ImageURL:    "/static/img/product4.jpg",
			},
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
		tmpl, err := template.ParseFiles(filepath.Join("templates", "seller_dashboard.html"))
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

// PrepareShipment handles preparing a shipment and generating a QR code
func PrepareShipment() http.HandlerFunc {
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
		orderID := r.FormValue("order_id")
		trackingNumber := r.FormValue("tracking_number")
		shippingNotes := r.FormValue("shipping_notes")

		// Handle package photo upload
		packagePhotoFile, packagePhotoHeader, err := r.FormFile("package_photo")
		if err != nil {
			http.Error(w, "Package photo is required", http.StatusBadRequest)
			return
		}
		defer packagePhotoFile.Close()

		// Read package photo data
		packagePhotoData, err := io.ReadAll(packagePhotoFile)
		if err != nil {
			http.Error(w, "Failed to read package photo", http.StatusInternalServerError)
			return
		}

		// Initialize manufacturer signature data
		var manufacturerSigData []byte

		// Check if manufacturer signature was uploaded
		manufacturerSigFile, manufacturerSigHeader, err := r.FormFile("manufacturer_signature")
		if err == nil {
			defer manufacturerSigFile.Close()

			// Read manufacturer signature data
			manufacturerSigData, err = io.ReadAll(manufacturerSigFile)
			if err != nil {
				http.Error(w, "Failed to read manufacturer signature", http.StatusInternalServerError)
				return
			}

			log.Printf("Manufacturer signature: %s", manufacturerSigHeader.Filename)
		}

		// In a real app, you would:
		// 1. Save the package photo and manufacturer signature to a storage system
		// 2. Get the URLs of the saved files
		// 3. Update the order status in the database

		// Generate a unique hash for the QR code
		// Combine order ID, timestamp, package photo data, and manufacturer signature data
		timestamp := time.Now().Unix()
		hashData := fmt.Sprintf("%s-%d", orderID, timestamp)

		// Add package photo hash
		packagePhotoHash := sha256.Sum256(packagePhotoData)
		hashData += "-" + hex.EncodeToString(packagePhotoHash[:])

		// Add manufacturer signature hash if available
		if len(manufacturerSigData) > 0 {
			manufacturerSigHash := sha256.Sum256(manufacturerSigData)
			hashData += "-" + hex.EncodeToString(manufacturerSigHash[:])
		}

		// Generate final hash
		hash := sha256.Sum256([]byte(hashData))
		qrCodeHash := hex.EncodeToString(hash[:])

		// In a real app, you would:
		// 1. Generate a QR code image from the hash
		// 2. Save the QR code image to a storage system
		// 3. Get the URL of the saved QR code image

		// For now, just create a mock QR code URL
		qrCodeURL := "/static/img/qr-codes/" + qrCodeHash + ".png"

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

		// In a real app, you would:
		// 1. Check if the user is logged in
		// 2. Check if the user is a seller
		// 3. Check if the order belongs to the seller
		// 4. Update the order status in the database

		// For now, just log the action
		log.Printf("Order %s marked as processing", orderID)

		// Redirect back to the seller dashboard
		http.Redirect(w, r, "/seller/dashboard", http.StatusSeeOther)
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

		// In a real app, you would:
		// 1. Check if the user is logged in
		// 2. Check if the user is a seller
		// 3. Check if the order belongs to the seller
		// 4. Update the order status in the database

		// For now, just log the action
		log.Printf("Order %s marked as shipped", orderID)

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

		// Parse form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Get QR code hash from form
		qrCodeHash := r.FormValue("qr_code_hash")

		// In a real app, you would:
		// 1. Check if the QR code hash exists in the database
		// 2. Check if the QR code has already been scanned the maximum number of times
		// 3. Increment the scan count
		// 4. If scanned by the buyer, update the order status to "delivered"
		// 5. If the order is now "delivered", release the funds from escrow

		// For now, just log the action
		log.Printf("QR code scanned: %s", qrCodeHash)

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success": true, "message": "Delivery confirmed"}`))
	}
}
