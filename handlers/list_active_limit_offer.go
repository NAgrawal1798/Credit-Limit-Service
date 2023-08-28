package handlers

import (
	"credit-limit-service/db"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

func ListActiveLimitOffers(w http.ResponseWriter, r *http.Request) {
	accountIDStr := r.URL.Query().Get("account_id")
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	activeDateStr := r.URL.Query().Get("active_date")
	var activeDate *time.Time
	if activeDateStr != "" {
		parsedActiveDate, err := time.Parse("2006-01-02T15:04:05Z", activeDateStr)
		if err != nil {
			http.Error(w, "Invalid active date", http.StatusBadRequest)
			return
		}
		activeDate = &parsedActiveDate
	}

	offers, err := db.GetActiveLimitOffers(accountID, activeDate)
	if err != nil {
		http.Error(w, "Error fetching active offers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(offers)
}
