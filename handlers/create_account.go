package handlers

import (
	"credit-limit-service/db"
	"credit-limit-service/models"
	"encoding/json"
	"net/http"
)

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	var account models.Account
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = db.InsertAccount(account)
	if err != nil {
		http.Error(w, "Error inserting account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Account created successfully"))
}
