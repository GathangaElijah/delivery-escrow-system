package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
)

// User represents a registered user
type User struct {
	ID              string
	Name            string
	Email           string
	Wallet          string
	Password        string // In a real app, this would be hashed
	Role            string // "buyer", "seller", "transporter", "manufacturer"
	ShippingAddress string // For buyers
	BusinessName    string // For sellers and manufacturers
	VehicleInfo     string // For transporters
}

// RegisterData holds data for the registration template
type RegisterData struct {
	Title string
	Error string
	Role  string
}

// LoginData holds data for the login template
type LoginData struct {
	Title string
	Error string
}

// RegisterBuyer handles buyer registration
func RegisterBuyer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// For GET requests, just show the form
		if r.Method == http.MethodGet {
			data := RegisterData{
				Title: "DES - Register as Buyer",
				Error: "",
				Role:  "buyer",
			}

			tmpl, err := template.ParseFiles(filepath.Join("templates", "buyer_registration_form.html"))
			if err != nil {
				http.Error(w, "Failed to load template", http.StatusInternalServerError)
				return
			}

			if err := tmpl.Execute(w, data); err != nil {
				http.Error(w, "Failed to render template", http.StatusInternalServerError)
				return
			}
			return
		}

		// For POST requests, process the form
		if r.Method == http.MethodPost {
			// Parse form data
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Failed to parse form", http.StatusBadRequest)
				return
			}

			// Get form values
			name := r.FormValue("name")
			email := r.FormValue("email")
			wallet := r.FormValue("wallet")
			shippingAddress := r.FormValue("shipping_address")
			password := r.FormValue("password")
			confirmPassword := r.FormValue("confirm_password")

			// Validate form data
			var errorMsg string

			// Check if passwords match
			if password != confirmPassword {
				errorMsg = "Passwords do not match"
			}

			// Validate email format
			emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
			if !emailRegex.MatchString(email) {
				errorMsg = "Invalid email format"
			}

			// Validate wallet address (basic check for Nova wallet format)
			// Nova wallet addresses typically start with "5" and are 48 characters long
			walletRegex := regexp.MustCompile(`^5[a-zA-Z0-9]{47}$`)
			if !walletRegex.MatchString(wallet) {
				errorMsg = "Invalid Nova wallet address format"
			}

			// If there are validation errors, show the form again with error message
			if errorMsg != "" {
				data := RegisterData{
					Title: "DES - Register as Buyer",
					Error: errorMsg,
					Role:  "buyer",
				}

				tmpl, err := template.ParseFiles(filepath.Join("templates", "buyer_registration_form.html"))
				if err != nil {
					http.Error(w, "Failed to load template", http.StatusInternalServerError)
					return
				}

				if err := tmpl.Execute(w, data); err != nil {
					http.Error(w, "Failed to render template", http.StatusInternalServerError)
					return
				}
				return
			}

			// In a real app, you would:
			// 1. Hash the password
			// 2. Store the user in a database
			// 3. Create a session for the user

			// Create a new buyer user
			buyer := User{
				ID:              "B" + email[:5], // Simple ID generation for demo
				Name:            name,
				Email:           email,
				Wallet:          wallet,
				Password:        password, // Should be hashed in production
				Role:            "buyer",
				ShippingAddress: shippingAddress,
			}

			// For now, just log the registration
			log.Printf("New buyer registered: %s (%s) with wallet %s", buyer.Name, buyer.Email, buyer.Wallet)

			// Redirect to home page or login page
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	}
}

// RegisterSeller handles seller registration
func RegisterSeller() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// For GET requests, just show the form
		if r.Method == http.MethodGet {
			data := RegisterData{
				Title: "DES - Register as Seller",
				Error: "",
				Role:  "seller",
			}

			tmpl, err := template.ParseFiles(filepath.Join("templates", "seller_registration_form.html"))
			if err != nil {
				http.Error(w, "Failed to load template", http.StatusInternalServerError)
				return
			}

			if err := tmpl.Execute(w, data); err != nil {
				http.Error(w, "Failed to render template", http.StatusInternalServerError)
				return
			}
			return
		}

		// For POST requests, process the form
		if r.Method == http.MethodPost {
			// Parse multipart form data (for file uploads)
			if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max memory
				http.Error(w, "Failed to parse form", http.StatusBadRequest)
				return
			}

			// Get form values
			businessName := r.FormValue("business_name")
			email := r.FormValue("email")
			paymentType := r.FormValue("payment_type")
			password := r.FormValue("password")
			confirmPassword := r.FormValue("confirm_password")

			// Get payment details based on type
			var paymentInfo string
			if paymentType == "wallet" {
				wallet := r.FormValue("wallet")
				paymentInfo = "Wallet: " + wallet
			} else if paymentType == "bank" {
				bankName := r.FormValue("bank_name")
				accountNumber := r.FormValue("account_number")
				paymentInfo = "Bank: " + bankName + ", Account: " + accountNumber
			}

			// Validate form data
			var errorMsg string

			// Check if passwords match
			if password != confirmPassword {
				errorMsg = "Passwords do not match"
			}

			// Validate email format
			emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
			if !emailRegex.MatchString(email) {
				errorMsg = "Invalid email format"
			}

			// Check for file uploads
			_, businessLicenseHeader, err := r.FormFile("business_license")
			if err != nil {
				errorMsg = "Business license file is required"
			}

			_, businessPermitHeader, err := r.FormFile("business_permit")
			if err != nil {
				errorMsg = "Business permit file is required"
			}

			// If there are validation errors, show the form again with error message
			if errorMsg != "" {
				data := RegisterData{
					Title: "DES - Register as Seller",
					Error: errorMsg,
					Role:  "seller",
				}

				tmpl, err := template.ParseFiles(filepath.Join("templates", "seller_registration_form.html"))
				if err != nil {
					http.Error(w, "Failed to load template", http.StatusInternalServerError)
					return
				}

				if err := tmpl.Execute(w, data); err != nil {
					http.Error(w, "Failed to render template", http.StatusInternalServerError)
					return
				}
				return
			}

			// In a real app, you would:
			// 1. Save the uploaded files to a storage system
			// 2. Hash the password
			// 3. Store the user in a database
			// 4. Create a session for the user

			// For now, just log the registration
			log.Printf("New seller registered: %s (%s)", businessName, email)
			log.Printf("Payment info: %s", paymentInfo)
			log.Printf("Business license: %s", businessLicenseHeader.Filename)
			log.Printf("Business permit: %s", businessPermitHeader.Filename)

			// Redirect to home page or login page
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	}
}

// RegisterTransporter handles transporter registration
func RegisterTransporter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// For GET requests, just show the form
		if r.Method == http.MethodGet {
			data := RegisterData{
				Title: "DES - Register as Transporter",
				Error: "",
				Role:  "transporter",
			}

			tmpl, err := template.ParseFiles(filepath.Join("templates", "transporter_registration_form.html"))
			if err != nil {
				http.Error(w, "Failed to load template", http.StatusInternalServerError)
				return
			}

			if err := tmpl.Execute(w, data); err != nil {
				http.Error(w, "Failed to render template", http.StatusInternalServerError)
				return
			}
			return
		}

		// For POST requests, process the form
		if r.Method == http.MethodPost {
			// Parse multipart form data (for file uploads)
			if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max memory
				http.Error(w, "Failed to parse form", http.StatusBadRequest)
				return
			}

			// Get form values
			name := r.FormValue("name")
			email := r.FormValue("email")
			wallet := r.FormValue("wallet")
			vehicleInfo := r.FormValue("vehicle_info")
			password := r.FormValue("password")
			confirmPassword := r.FormValue("confirm_password")

			// Validate form data
			var errorMsg string

			// Check if passwords match
			if password != confirmPassword {
				errorMsg = "Passwords do not match"
			}

			// Validate email format
			emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
			if !emailRegex.MatchString(email) {
				errorMsg = "Invalid email format"
			}

			// Validate wallet address (basic check for Nova wallet format)
			walletRegex := regexp.MustCompile(`^5[a-zA-Z0-9]{47}$`)
			if !walletRegex.MatchString(wallet) {
				errorMsg = "Invalid Nova wallet address format"
			}

			// Check for ID card upload
			_, idCardHeader, err := r.FormFile("id_card")
			if err != nil {
				errorMsg = "National ID card is required"
			}

			// If there are validation errors, show the form again with error message
			if errorMsg != "" {
				data := RegisterData{
					Title: "DES - Register as Transporter",
					Error: errorMsg,
					Role:  "transporter",
				}

				tmpl, err := template.ParseFiles(filepath.Join("templates", "transporter_registration_form.html"))
				if err != nil {
					http.Error(w, "Failed to load template", http.StatusInternalServerError)
					return
				}

				if err := tmpl.Execute(w, data); err != nil {
					http.Error(w, "Failed to render template", http.StatusInternalServerError)
					return
				}
				return
			}

			// In a real app, you would:
			// 1. Save the uploaded ID card to a storage system
			// 2. Hash the password
			// 3. Store the user in a database
			// 4. Create a session for the user

			// Create a new transporter user
			transporter := User{
				ID:          "T" + email[:5], // Simple ID generation for demo
				Name:        name,
				Email:       email,
				Wallet:      wallet,
				Password:    password, // Should be hashed in production
				Role:        "transporter",
				VehicleInfo: vehicleInfo,
			}

			// For now, just log the registration
			log.Printf("New transporter registered: %s (%s) with wallet %s", transporter.Name, transporter.Email, transporter.Wallet)
			log.Printf("ID Card: %s", idCardHeader.Filename)
			log.Printf("Vehicle Info: %s", transporter.VehicleInfo)

			// Redirect to home page or login page
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	}
}

// Register handles the main registration page and redirects to specific registration forms
func Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Redirect to buyer registration by default
		http.Redirect(w, r, "/register/buyer", http.StatusSeeOther)
	}
}

// Login handles user login
func Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// For GET requests, just show the form
		if r.Method == http.MethodGet {
			data := LoginData{
				Title: "DES - Login",
				Error: "",
			}

			tmpl, err := template.ParseFiles(filepath.Join("templates", "login.html"))
			if err != nil {
				http.Error(w, "Failed to load template", http.StatusInternalServerError)
				return
			}

			if err := tmpl.Execute(w, data); err != nil {
				http.Error(w, "Failed to render template", http.StatusInternalServerError)
				return
			}
			return
		}

		// For POST requests, process the form
		if r.Method == http.MethodPost {
			// Parse form data
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Failed to parse form", http.StatusBadRequest)
				return
			}

			// Get form values
			email := r.FormValue("email")
			password := r.FormValue("password")

			// In a real app, you would:
			// 1. Look up the user by email in your database
			// 2. Verify the password hash
			// 3. Create a session for the user

			// For now, just log the login attempt
			log.Printf("Login attempt: %s", email)

			// Simulate a successful login
			// In a real app, you would check credentials against a database
			if email != "" && password != "" {
				// Redirect to home page after successful login
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			// If login fails, show the form again with an error message
			data := LoginData{
				Title: "DES - Login",
				Error: "Invalid email or password",
			}

			tmpl, err := template.ParseFiles(filepath.Join("templates", "login.html"))
			if err != nil {
				http.Error(w, "Failed to load template", http.StatusInternalServerError)
				return
			}

			if err := tmpl.Execute(w, data); err != nil {
				http.Error(w, "Failed to render template", http.StatusInternalServerError)
				return
			}
		}
	}
}
