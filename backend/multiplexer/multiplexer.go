package multiplexer

import (
	"des/backend/handlers"
	"net/http"
	"strings"
)

// Here we are registering all incoming request patterns.
func Multiplexer() {
	mux := http.NewServeMux()

	// Serve static files
	files_handler := http.FileServer(http.Dir("../static"))
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

}
