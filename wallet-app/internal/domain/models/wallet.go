package models

import (
    "time"

    "github.com/google/uuid"
)

type TransactionType string

const (
    Deposit  TransactionType = "deposit"
    Withdraw TransactionType = "withdraw"
)

type Wallet struct {
    ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
    UserID    uuid.UUID `json:"user_id" gorm:"type:uuid"`
    Balance   float64   `json:"balance"`
    Currency  string    `json:"currency"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type Transaction struct {
    ID          uuid.UUID       `json:"id" gorm:"type:uuid;primary_key"`
    WalletID    uuid.UUID       `json:"wallet_id" gorm:"type:uuid"`
    Type        TransactionType `json:"type"`
    Amount      float64         `json:"amount"`
    Description string         `json:"description"`
    CreatedAt   time.Time      `json:"created_at"`
}