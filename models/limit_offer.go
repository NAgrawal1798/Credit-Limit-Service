package models

import (
	"time"
)

type OfferStatus string

const (
	OfferStatusPending  OfferStatus = "PENDING"
	OfferStatusAccepted OfferStatus = "ACCEPTED"
	OfferStatusRejected OfferStatus = "REJECTED"
)

type LimitType string

const (
	AccountLimitType        LimitType = "ACCOUNT_LIMIT"
	PerTransactionLimitType LimitType = "PER_TRANSACTION_LIMIT"
)

type LimitOffer struct {
	OfferID             int         `json:"offer_id"`
	AccountID           int         `json:"account_id"`
	LimitType           LimitType   `json:"limit_type"`
	NewLimit            int         `json:"new_limit"`
	OfferActivationTime time.Time   `json:"offer_activation_time"`
	OfferExpiryTime     time.Time   `json:"offer_expiry_time"`
	Status              OfferStatus `json:"status"`
	CreatedAt           time.Time   `json:"created_at"`
	UpdatedAt           *time.Time  `json:"updated_at"`
}
