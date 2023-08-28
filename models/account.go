package models

import (
	"time"
)

type Account struct {
	AccountID                     int       `json:"account_id"` // primary_key
	CustomerID                    int       `json:"customer_id"`
	AccountLimit                  int       `json:"account_limit"`
	PerTransactionLimit           int       `json:"per_transaction_limit"`
	LastAccountLimit              int       `json:"last_account_limit"`
	LastPerTransactionLimit       int       `json:"last_per_transaction_limit"`
	AccountLimitUpdateTime        time.Time `json:"account_limit_update_time"`
	PerTransactionLimitUpdateTime time.Time `json:"per_transaction_limit_update_time"`
}
