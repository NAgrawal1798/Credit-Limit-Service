package main

import (
	"credit-limit-service/db"
	"credit-limit-service/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize the database connection
	err := db.InitDB()
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	// Create limit offer
	r.HandleFunc("/create-limit-offer", handlers.CreateLimitOffer).Methods("POST")

	// List active limit offers
	r.HandleFunc("/list-active-limit-offers", handlers.ListActiveLimitOffers).Methods("GET")

	// Update limit offer status
	r.HandleFunc("/update-limit-offer-status/{limit_offer_id}/{status}", handlers.UpdateLimitOfferStatus).Methods("PUT")

	// Create account
	r.HandleFunc("/create-account", handlers.CreateAccount).Methods("POST")

	// Get account
	r.HandleFunc("/get-account/{account_id}", handlers.GetAccount).Methods("GET")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
