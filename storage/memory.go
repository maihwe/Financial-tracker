package storage

import "financial-tracker/models"

// transactions is our temporary in-memory database.
// It stores both income and expenses.
var transactions = []models.Transaction{
	{
		ID:         1,
		Title:      "Starting Salary",
		Amount:     200000,
		CategoryID: 3,
		Type:       "income",
	},
	{
		ID:         2,
		Title:      "Groceries",
		Amount:     25000,
		CategoryID: 1,
		Type:       "expense",
	},
}

// nextID keeps track of the next available transaction ID.
//
// We do not calculate IDs using len(transactions)
// because deleting a transaction changes the length
// of the slice and could create duplicate IDs.
var nextID = 3

// GetAllTransactions returns all transactions currently stored.
func GetAllTransactions() []models.Transaction {
	return transactions
}

// AddTransaction adds a new transaction to storage.
func AddTransaction(transaction models.Transaction) models.Transaction {

	// Assign the next unique ID.
	transaction.ID = nextID

	// Increase the counter so the next transaction
	// receives a different ID.
	nextID++

	// Add the transaction to storage.
	transactions = append(transactions, transaction)

	// Return the newly created transaction.
	return transaction
}

// GetTransactionByID searches for a transaction using its ID.
// It returns the transaction and true when found.
// It returns an empty transaction and false when not found.
func GetTransactionByID(id int) (models.Transaction, bool) {

	// Loop through all stored transactions.
	for _, transaction := range transactions {

		// Check whether the current transaction
		// has the requested ID.
		if transaction.ID == id {
			return transaction, true
		}
	}

	// The transaction was not found.
	return models.Transaction{}, false
}

// UpdateTransaction replaces an existing transaction.
func UpdateTransaction(
	id int,
	updated models.Transaction,
) (models.Transaction, bool) {

	// Search through all transactions.
	for i, transaction := range transactions {

		// Find the transaction that matches the requested ID.
		if transaction.ID == id {

			// Keep the original ID.
			updated.ID = id

			// Replace the old transaction.
			transactions[i] = updated

			// Return the updated transaction.
			return updated, true
		}
	}

	// The transaction was not found.
	return models.Transaction{}, false
}

// DeleteTransaction removes a transaction from storage.
func DeleteTransaction(id int) bool {

	// Search through all transactions.
	for i, transaction := range transactions {

		// Find the transaction we want to delete.
		if transaction.ID == id {

			// Remove the transaction from the slice.
			transactions = append(
				transactions[:i],
				transactions[i+1:]...,
			)

			// Confirm successful deletion.
			return true
		}
	}

	// The transaction was not found.
	return false
}

// GetFinancialSummary calculates total income,
// total expenses, and the current balance.
func GetFinancialSummary() (float64, float64, float64) {

	var totalIncome float64
	var totalExpenses float64

	// Examine every transaction.
	for _, transaction := range transactions {

		// Add income to total income.
		if transaction.Type == "income" {
			totalIncome += transaction.Amount
		}

		// Add expenses to total expenses.
		if transaction.Type == "expense" {
			totalExpenses += transaction.Amount
		}
	}

	// Calculate the current balance.
	balance := totalIncome - totalExpenses

	return totalIncome, totalExpenses, balance
}
