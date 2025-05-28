package multiplexer

import (
	"des/backend/internal"
	"net/http"
	"strings"
)

// Here we are registering all incoming request patterns.
func Multiplexer() *http.ServeMux{
	mux := http.NewServeMux()

	// Serve static files
	files_handler := http.FileServer(http.Dir("./backend/static"))
	mux.Handle("/static/", http.StripPrefix("/static", files_handler))

	// Home page
	mux.HandleFunc("/",internal.Index())

	// Product detail page
	mux.Handle("/product/",internal.GetProduct())

	// Authentication routes
	mux.Handle("/register",internal.Register())
	mux.Handle("/register/buyer",internal.RegisterBuyer())
	mux.Handle("/register/seller",internal.RegisterSeller())
	mux.Handle("/register/transporter",internal.RegisterTransporter())
	mux.Handle("/login",internal.Login())

	// Seller dashboard routes
	mux.Handle("/seller/dashboard",internal.SellerDashboard())
	mux.Handle("/seller/products/add",internal.AddProduct())
	mux.Handle("/seller/shipping/prepare",internal.PrepareShipment())
	mux.Handle("/seller/orders/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/process") {
		internal.ProcessOrder().ServeHTTP(w, r)
		} else if strings.HasSuffix(r.URL.Path, "/ship") {
		internal.ShipOrder().ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}
	}))
	mux.Handle("/seller/account/update",internal.UpdateAccount())

	// Escrow-related routes
	mux.Handle("/escrow/status/",internal.GetEscrowStatus())
	mux.Handle("/scan-qr",internal.ScanQRCode())

	// Checkout routes
	mux.Handle("/checkout/",internal.Checkout())
	mux.Handle("/checkout/process",internal.ProcessCheckout())
	mux.Handle("/checkout/confirmation",internal.CheckoutConfirmation())

	// Cart routes
	mux.Handle("/cart",internal.ViewCart())
	mux.Handle("/cart/add",internal.AddToCart())
	mux.Handle("/cart/update",internal.UpdateCart())
	mux.Handle("/cart/clear",internal.ClearCart())

	return mux
}
