package models

import "time"

type TransactionType string

const (
	Deposit     TransactionType = "deposit"
	Withdraw    TransactionType = "withdraw"
	TransferIn  TransactionType = "transfer_in"
	TransferOut TransactionType = "transfer_out"
)

type Transaction struct {
	ID          int             `json:"id"`
	AccountID   int             `json:"account_id"`
	Type        TransactionType `json:"type"`
	Amount      float64         `json:"amount"`
	TargetID    *int            `json:"target_id,omitempty"`
	Description string          `json:"description,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
}
