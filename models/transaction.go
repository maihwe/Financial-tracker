package models

// Transaction represents one movement of money.
// A transaction can either be income or an expense.
type Transaction struct {
	ID       int     `json:"id"`
	Title    string  `json:"title"`
	Amount   float64 `json:"amount"`
	Category string  `json:"category"`
	Type     string  `json:"type"`
}
