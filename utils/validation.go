package utils

import (
	"strings"

	"financial-tracker/models"
)

// ValidateTransaction checks whether a transaction
// contains acceptable information.
func ValidateTransaction(transaction models.Transaction) string {

	// Remove unnecessary spaces from the title.
	title := strings.TrimSpace(transaction.Title)

	// Make sure a title was provided.
	if title == "" {
		return "Title is required"
	}

	// Make sure the amount is greater than zero.
	if transaction.Amount <= 0 {
		return "Amount must be greater than zero"
	}

	// Remove unnecessary spaces from the category.
	category := strings.TrimSpace(transaction.Category)

	// Make sure a category was provided.
	if category == "" {
		return "Category is required"
	}

	// Convert the transaction type to lowercase
	// so INCOME, Income, and income can be treated consistently.
	transactionType := strings.ToLower(strings.TrimSpace(transaction.Type))

	// Only income and expense are valid transaction types.
	if transactionType != "income" && transactionType != "expense" {
		return "Type must be income or expense"
	}

	// Return an empty string when validation succeeds.
	return ""
}
