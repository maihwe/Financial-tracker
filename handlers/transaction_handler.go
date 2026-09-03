package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"financial-tracker/models"
	"financial-tracker/storage"
	"financial-tracker/utils"
)

// TransactionHandler handles requests related to transactions.
func TransactionHandler(w http.ResponseWriter, r *http.Request) {

	// Tell the client that our responses are JSON.
	w.Header().Set("Content-Type", "application/json")

	// Handle PUT requests.
	if r.Method == http.MethodPut {

		// Split a path such as /transactions/3.
		// The result is:
		// ["", "transactions", "3"]
		parts := strings.Split(r.URL.Path, "/")

		// Make sure the path contains a transaction ID.
		if len(parts) != 3 || parts[1] != "transactions" {
			http.Error(w, "Invalid transaction path", http.StatusBadRequest)
			return
		}

		// Convert the ID from a string into an integer.
		id, err := strconv.Atoi(parts[2])

		// Stop if the ID is not a valid number.
		if err != nil {
			http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
			return
		}

		// Create an empty transaction to receive
		// the client's updated JSON data.
		var updatedTransaction models.Transaction

		// Decode the request body.
		err = json.NewDecoder(r.Body).Decode(&updatedTransaction)

		// Stop if the JSON is invalid.
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Validate the new transaction data.
		validationError := utils.ValidateTransaction(updatedTransaction)

		// Stop if validation fails.
		if validationError != "" {
			http.Error(w, validationError, http.StatusBadRequest)
			return
		}

		// Update the transaction in storage.
		updatedTransaction, found := storage.UpdateTransaction(
			id,
			updatedTransaction,
		)

		// Return 404 if the transaction does not exist.
		if !found {
			http.Error(w, "Transaction not found", http.StatusNotFound)
			return
		}

		// Return the updated transaction.
		json.NewEncoder(w).Encode(updatedTransaction)

		return
	}

	// Handle DELETE requests.
	if r.Method == http.MethodDelete {

		// Split a path such as /transactions/3.
		// The result is:
		// ["", "transactions", "3"]
		parts := strings.Split(r.URL.Path, "/")

		// Make sure the path contains a transaction ID.
		if len(parts) != 3 || parts[1] != "transactions" {
			http.Error(w, "Invalid transaction path", http.StatusBadRequest)
			return
		}

		// Convert the ID from a string into an integer.
		id, err := strconv.Atoi(parts[2])

		// Stop if the ID is not a valid number.
		if err != nil {
			http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
			return
		}

		// Try to delete the transaction.
		deleted := storage.DeleteTransaction(id)

		// Return 404 if the transaction does not exist.
		if !deleted {
			http.Error(w, "Transaction not found", http.StatusNotFound)
			return
		}

		// Tell the client that the transaction was successfully deleted.
		w.WriteHeader(http.StatusNoContent)

		return
	}

	// Handle POST requests.
	if r.Method == http.MethodPost {

		// Create an empty transaction to receive
		// the client's JSON data.
		var newTransaction models.Transaction

		// Decode the JSON request body.
		err := json.NewDecoder(r.Body).Decode(&newTransaction)

		// Stop if the JSON is invalid.
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Validate the transaction before saving it.
		validationError := utils.ValidateTransaction(newTransaction)

		// Stop if validation fails.
		if validationError != "" {
			http.Error(w, validationError, http.StatusBadRequest)
			return
		}

		// Store the valid transaction.
		newTransaction = storage.AddTransaction(newTransaction)

		// Tell the client that a new resource was created.
		w.WriteHeader(http.StatusCreated)

		// Return the created transaction as JSON.
		json.NewEncoder(w).Encode(newTransaction)

		return
	}

	// Handle GET requests.
	if r.Method == http.MethodGet {

		// If the path is exactly /transactions,
		// return every transaction.
		if r.URL.Path == "/transactions" {

			// Get all transactions from storage.
			transactions := storage.GetAllTransactions()

			// Return them as JSON.
			json.NewEncoder(w).Encode(transactions)

			return
		}

		// Split a path such as /transactions/3.
		// The result is:
		// ["", "transactions", "3"]
		parts := strings.Split(r.URL.Path, "/")

		// Make sure the path has the correct structure.
		if len(parts) != 3 || parts[1] != "transactions" {
			http.Error(w, "Invalid transaction path", http.StatusBadRequest)
			return
		}

		// Convert the ID from string to integer.
		id, err := strconv.Atoi(parts[2])

		// Stop if the ID is not a valid number.
		if err != nil {
			http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
			return
		}

		// Search for the transaction.
		transaction, found := storage.GetTransactionByID(id)

		// Return 404 if it does not exist.
		if !found {
			http.Error(w, "Transaction not found", http.StatusNotFound)
			return
		}

		// Return the transaction as JSON.
		json.NewEncoder(w).Encode(transaction)

		return
	}

	// Handle unsupported HTTP methods.
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// SummaryHandler returns our financial summary.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {

	// Only GET is allowed for the summary.
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Tell the client that we are returning JSON.
	w.Header().Set("Content-Type", "application/json")

	// Calculate the financial totals.
	totalIncome, totalExpenses, balance := storage.GetFinancialSummary()

	// Create the response.
	summary := map[string]float64{
		"total_income":   totalIncome,
		"total_expenses": totalExpenses,
		"balance":        balance,
	}

	// Return the summary as JSON.
	json.NewEncoder(w).Encode(summary)
}
