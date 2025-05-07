package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strings"

	"des/backend/handlers"
)

// openBrowser tries to open the URL in the default browser
func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		log.Printf("Error opening browser: %v", err)
	}
}

func main() {
	mux := http.NewServeMux()

	// Serve static files
	files_handler := http.FileServer(http.Dir("./backend/static"))
	mux.Handle("/static/", http.StripPrefix("/static", files_handler))

	// Home page
	mux.Handle("/", handlers.Index())

	// Product detail page
	mux.Handle("/product/", handlers.GetProduct())

	// Authentication routes
	mux.Handle("/register", handlers.Register())
	mux.Handle("/register/buyer", handlers.RegisterBuyer())
	mux.Handle("/register/seller", handlers.RegisterSeller())
	mux.Handle("/register/transporter", handlers.RegisterTransporter())
	mux.Handle("/login", handlers.Login())

	// Seller dashboard routes
	mux.Handle("/seller/dashboard", handlers.SellerDashboard())
	mux.Handle("/seller/products/add", handlers.AddProduct())
	mux.Handle("/seller/shipping/prepare", handlers.PrepareShipment())
	mux.Handle("/seller/orders/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/process") {
			handlers.ProcessOrder().ServeHTTP(w, r)
		} else if strings.HasSuffix(r.URL.Path, "/ship") {
			handlers.ShipOrder().ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}
	}))
	mux.Handle("/seller/account/update", handlers.UpdateAccount())
	
	// Escrow-related routes
	mux.Handle("/escrow/status/", handlers.GetEscrowStatus())
	mux.Handle("/scan-qr", handlers.ScanQRCode())
	
	// Checkout routes
	mux.Handle("/checkout/", handlers.Checkout())
	mux.Handle("/checkout/process", handlers.ProcessCheckout())
	mux.Handle("/checkout/confirmation", handlers.CheckoutConfirmation())

	// Cart routes
	mux.Handle("/cart", handlers.ViewCart())
	mux.Handle("/cart/add", handlers.AddToCart())
	mux.Handle("/cart/update", handlers.UpdateCart())
	mux.Handle("/cart/clear", handlers.ClearCart())

	// Add debug routes
	mux.Handle("/debug/info", handlers.DebugInfo())
	mux.Handle("/debug/cart", handlers.DebugCart())
	mux.Handle("/debug/create-test-order", handlers.CreateTestOrder())
	mux.Handle("/debug/add-product2", handlers.AddProduct2ToCart())
	mux.Handle("/test/cart", handlers.TestCart())

	// Define server address
	addr := ":8080"
	url := fmt.Sprintf("http://localhost%s", addr)

	// Print clickable link in terminal
	fmt.Printf("Server starting at \033[34m\033[4m%s\033[0m\n", url)

	// Start the server
	log.Printf("Starting server on %s", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
