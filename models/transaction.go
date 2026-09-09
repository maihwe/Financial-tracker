package models

import "time"

// Transaction represents one movement of money.
// A transaction can either be income or an expense.
type Transaction struct {
	ID            int       `json:"id"`
	UserID        int       `json:"-"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Amount        float64   `json:"amount"`
	CategoryID    int       `json:"category_id"`
	Type          string    `json:"type"`
	TransactionAt time.Time `json:"transaction_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
