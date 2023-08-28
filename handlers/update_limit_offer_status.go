package handlers

import (
	"credit-limit-service/db"
	"credit-limit-service/models"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func UpdateLimitOfferStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	limitOfferIDStr := vars["limit_offer_id"]
	limitOfferID, err := strconv.Atoi(limitOfferIDStr)
	if err != nil {
		http.Error(w, "Invalid limit offer ID", http.StatusBadRequest)
		return
	}

	status := vars["status"]
	if status != string(models.OfferStatusAccepted) && status != string(models.OfferStatusRejected) {
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	err = db.UpdateLimitOfferStatus(limitOfferID, status)
	if err != nil {
		http.Error(w, "Error updating limit offer status", http.StatusInternalServerError)
		return
	}

	if status == string(models.OfferStatusAccepted) {
		offer, err := db.GetLimitOfferByID(limitOfferID)
		if err != nil {
			http.Error(w, "Error fetching limit offer", http.StatusInternalServerError)
			return
		}

		account, err := db.GetAccountByID(offer.AccountID)
		if err != nil {
			http.Error(w, "Error fetching account", http.StatusInternalServerError)
			return
		}

		err = updateAccountBasedOnLimitOffer(account, offer)
		if err != nil {
			http.Error(w, "Error updating account", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Limit offer status updated successfully"))
}

func updateAccountBasedOnLimitOffer(account models.Account, offer models.LimitOffer) error {
	if offer.LimitType == models.AccountLimitType {
		account.AccountLimit = offer.NewLimit
		account.AccountLimitUpdateTime = time.Now()
	} else if offer.LimitType == models.PerTransactionLimitType {
		account.PerTransactionLimit = offer.NewLimit
		account.PerTransactionLimitUpdateTime = time.Now()
	} else {
		return errors.New("Invalid LimitType")
	}

	err := db.UpdateAccount(account, offer)
	if err != nil {
		log.Println("Error updating account:", err)
		return err
	}
	return nil
}
