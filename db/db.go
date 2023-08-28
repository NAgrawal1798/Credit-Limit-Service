package db

import (
	"credit-limit-service/models"
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var (
	db *sql.DB
)

func InitDB() error {
	connectionString := "user=postgres dbname=credit_limit_db password=Vegapay sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	// Create the accounts table if it doesn't exist
	err = createAccountsTable()
	if err != nil {
		return err
	}

	// Create the limit_offers table if it doesn't exist
	err = createLimitOffersTable()
	if err != nil {
		return err
	}

	log.Println("Connected to the database successfully")
	return nil
}

func createLimitOffersTable() error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS limit_offers (
            offer_id SERIAL PRIMARY KEY,
            account_id INT NOT NULL,
            limit_type VARCHAR(20) NOT NULL,
            new_limit INT NOT NULL,
            offer_activation_time TIMESTAMP NOT NULL,
            offer_expiry_time TIMESTAMP NOT NULL,
            status VARCHAR(20) NOT NULL,
            created_at TIMESTAMP DEFAULT NOW(),
            updated_at TIMESTAMP
        );
    `)
	if err != nil {
		return err
	}
	return nil
}

func createAccountsTable() error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS accounts (
            account_id SERIAL PRIMARY KEY,
            customer_id INT,
            account_limit INT,
            per_transaction_limit INT,
            last_account_limit INT,
            last_per_transaction_limit INT,
            account_limit_update_time TIMESTAMP,
            per_transaction_limit_update_time TIMESTAMP
        );
    `)
	if err != nil {
		return err
	}
	return nil
}

func GetAccountByID(accountID int) (models.Account, error) {
	var account models.Account
	err := db.QueryRow(
		"SELECT * FROM accounts WHERE account_id = $1",
		accountID,
	).Scan(
		&account.AccountID, &account.CustomerID, &account.AccountLimit,
		&account.PerTransactionLimit, &account.LastAccountLimit,
		&account.LastPerTransactionLimit, &account.AccountLimitUpdateTime,
		&account.PerTransactionLimitUpdateTime,
	)
	if err != nil {
		log.Println("Error fetching account:", err)
		return models.Account{}, err
	}
	return account, nil
}

func InsertAccount(account models.Account) error {

	_, err := db.Exec(
		"INSERT INTO accounts (customer_id, account_limit, per_transaction_limit, last_account_limit, last_per_transaction_limit, account_limit_update_time, per_transaction_limit_update_time) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		account.CustomerID, account.AccountLimit, account.PerTransactionLimit,
		account.LastAccountLimit, account.LastPerTransactionLimit,
		account.AccountLimitUpdateTime, account.PerTransactionLimitUpdateTime,
	)
	if err != nil {
		log.Println("Error inserting account:", err)
		return err
	}
	return nil
}

func InsertLimitOffer(offer models.LimitOffer) error {
	_, err := db.Exec(
		"INSERT INTO limit_offers (account_id, limit_type, new_limit, offer_activation_time, offer_expiry_time, status) VALUES ($1, $2, $3, $4, $5, $6)",
		offer.AccountID, offer.LimitType, offer.NewLimit, offer.OfferActivationTime, offer.OfferExpiryTime, models.OfferStatusPending,
	)
	if err != nil {
		return err
	}
	return nil
}

func GetActiveLimitOffers(accountID int, activeDate *time.Time) ([]models.LimitOffer, error) {
	var offers []models.LimitOffer

	query := `
		SELECT * FROM limit_offers 
		WHERE account_id = $1 AND status = $2 
		AND ($3::timestamp IS NULL OR (offer_activation_time <= $3 AND offer_expiry_time >= $3::timestamp))
	`

	rows, err := db.Query(query, accountID, models.OfferStatusPending, activeDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var offer models.LimitOffer
		var updatedAt sql.NullTime
		err := rows.Scan(
			&offer.OfferID, &offer.AccountID, &offer.LimitType,
			&offer.NewLimit, &offer.OfferActivationTime, &offer.OfferExpiryTime,
			&offer.Status, &offer.CreatedAt, &updatedAt,
		)
		if err != nil {
			return nil, err
		}
		if updatedAt.Valid {
			offer.UpdatedAt = &updatedAt.Time
		} else {
			offer.UpdatedAt = nil
		}
		offers = append(offers, offer)
	}

	if len(offers) == 0 {
		return []models.LimitOffer{}, nil
	}

	return offers, nil
}

func UpdateLimitOfferStatus(limitOfferID int, status string) error {
	activeTime := time.Now()

	result, err := db.Exec(
		"UPDATE limit_offers SET status = $1, updated_at = $2 WHERE offer_id = $3 AND status = $4 AND offer_activation_time <= $5 AND offer_expiry_time >= $6",
		status, activeTime, limitOfferID, models.OfferStatusPending, activeTime, activeTime,
	)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("Offer is not active or not pending")
	}

	return nil
}

func GetLimitOfferByID(offerID int) (models.LimitOffer, error) {
	var offer models.LimitOffer
	err := db.QueryRow(
		"SELECT * FROM limit_offers WHERE offer_id = $1",
		offerID,
	).Scan(
		&offer.OfferID, &offer.AccountID, &offer.LimitType,
		&offer.NewLimit, &offer.OfferActivationTime, &offer.OfferExpiryTime,
		&offer.Status, &offer.CreatedAt, &offer.UpdatedAt,
	)
	if err != nil {
		log.Println("Error fetching limit offer:", err)
		return models.LimitOffer{}, err
	}
	return offer, nil
}

func UpdateAccount(account models.Account, limitOffer models.LimitOffer) error {
	var updateQuery string
	var args []interface{}

	if limitOffer.LimitType == models.AccountLimitType {
		updateQuery = "UPDATE accounts SET account_limit = $1, last_account_limit = $2, account_limit_update_time = $3 WHERE account_id = $4"
		args = append(args, limitOffer.NewLimit, account.AccountLimit, account.AccountLimitUpdateTime, account.AccountID)
	} else if limitOffer.LimitType == models.PerTransactionLimitType {
		updateQuery = "UPDATE accounts SET per_transaction_limit = $1, last_per_transaction_limit = $2, per_transaction_limit_update_time = $3 WHERE account_id = $4"
		args = append(args, limitOffer.NewLimit, account.PerTransactionLimit, account.PerTransactionLimitUpdateTime, account.AccountID)
	} else {
		return errors.New("Invalid LimitType")
	}

	_, err := db.Exec(updateQuery, args...)
	if err != nil {
		log.Println("Error updating account:", err)
		return err
	}
	return nil
}
