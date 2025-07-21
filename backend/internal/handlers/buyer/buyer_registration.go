package auth

import (
	"des/backend/internal/models/buyer"
	"des/backend/internal/service/buyer_service"
	"encoding/json"
	"net/http"
)

func RegisterBuyer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "MethodNotAllowed", http.StatusMethodNotAllowed)
			return
		}

		// This is imported from the buyer package in the models
		var buyerRegister buyer.BuyerRegister

		if err := json.NewDecoder(r.Body).Decode(&buyerRegister); err != nil {
			http.Error(w, "Invalid JSON data", http.StatusBadRequest)
			return
		}

		// This is imported from the buyer_service package in services
		if err := buyer_service.ValidateBuyerRegistration(&buyerRegister); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Buyer registration received",
			"email":   buyerRegister.EmailAddress,
		})
	}

}
