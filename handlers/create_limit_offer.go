package handlers

import (
	"credit-limit-service/db"
	"credit-limit-service/models"
	"encoding/json"
	"net/http"
)

func CreateLimitOffer(w http.ResponseWriter, r *http.Request) {
	var offer models.LimitOffer
	err := json.NewDecoder(r.Body).Decode(&offer)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate the offer data
	if offer.NewLimit <= 0 || offer.OfferActivationTime.After(offer.OfferExpiryTime) || (offer.LimitType != models.AccountLimitType && offer.LimitType != models.PerTransactionLimitType) {
		http.Error(w, "Invalid offer data", http.StatusBadRequest)
		return
	}

	account, err := db.GetAccountByID(offer.AccountID)
	if err != nil {
		http.Error(w, "Error fetching account", http.StatusInternalServerError)
		return
	}

	// check the new limit is greater than the current limit
	if offer.LimitType == models.AccountLimitType && offer.NewLimit <= account.AccountLimit {
		http.Error(w, "New account limit should be greater than the current limit", http.StatusBadRequest)
		return
	}
	if offer.LimitType == models.PerTransactionLimitType && offer.NewLimit <= account.PerTransactionLimit {
		http.Error(w, "New per transaction limit should be greater than the current limit", http.StatusBadRequest)
		return
	}

	offer.Status = models.OfferStatusPending

	err = db.InsertLimitOffer(offer)
	if err != nil {
		http.Error(w, "Error creating limit offer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Limit offer created successfully"))
}
