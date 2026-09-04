package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"financial-tracker/models"
	"financial-tracker/storage"
	"financial-tracker/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TransactionHandler handles transaction-related HTTP requests.
func TransactionHandler(pool *pgxpool.Pool) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		// Handle GET requests.
		if r.Method == http.MethodGet {

			// GET /transactions
			if r.URL.Path == "/transactions" {

				transactions, err := storage.GetAllTransactionsFromDB(pool)

				if err != nil {
					http.Error(
						w,
						"Could not retrieve transactions",
						http.StatusInternalServerError,
					)
					return
				}

				json.NewEncoder(w).Encode(transactions)
				return
			}

			// GET /transactions/{id}
			parts := strings.Split(
				strings.Trim(r.URL.Path, "/"),
				"/",
			)

			// Expected:
			// ["transactions", "18"]
			if len(parts) != 2 || parts[0] != "transactions" {
				http.Error(
					w,
					"Invalid transaction path",
					http.StatusBadRequest,
				)
				return
			}

			// Convert the ID from string to integer.
			id, err := strconv.Atoi(parts[1])

			if err != nil {
				http.Error(
					w,
					"Invalid transaction ID",
					http.StatusBadRequest,
				)
				return
			}

			// Get the transaction from PostgreSQL.
			transaction, err := storage.GetTransactionByIDFromDB(
				pool,
				id,
			)

			// Return 404 when the transaction doesn't exist.
			if err != nil {
				http.Error(
					w,
					"Transaction not found",
					http.StatusNotFound,
				)
				return
			}

			// Return the transaction.
			json.NewEncoder(w).Encode(transaction)
			return
		}

		// Handle POST requests.
		if r.Method == http.MethodPost {

			var newTransaction models.Transaction

			err := json.NewDecoder(r.Body).Decode(
				&newTransaction,
			)

			if err != nil {
				http.Error(
					w,
					"Invalid JSON",
					http.StatusBadRequest,
				)
				return
			}

			validationError := utils.ValidateTransaction(
				newTransaction,
			)

			if validationError != "" {
				http.Error(
					w,
					validationError,
					http.StatusBadRequest,
				)
				return
			}

			newTransaction, err =
				storage.AddTransactionToDB(
					pool,
					newTransaction,
				)

			if err != nil {
				http.Error(
					w,
					"Could not save transaction",
					http.StatusInternalServerError,
				)
				return
			}

			w.WriteHeader(http.StatusCreated)

			json.NewEncoder(w).Encode(newTransaction)
			return
		}

		// Handle PUT requests.
		if r.Method == http.MethodPut {

			// Extract the transaction ID from the URL.
			parts := strings.Split(
				strings.Trim(r.URL.Path, "/"),
				"/",
			)

			// Expected:
			// ["transactions", "18"]
			if len(parts) != 2 || parts[0] != "transactions" {
				http.Error(
					w,
					"Invalid transaction path",
					http.StatusBadRequest,
				)
				return
			}

			// Convert the ID from string to integer.
			id, err := strconv.Atoi(parts[1])

			if err != nil {
				http.Error(
					w,
					"Invalid transaction ID",
					http.StatusBadRequest,
				)
				return
			}

			// Decode the new transaction values.
			var updatedTransaction models.Transaction

			err = json.NewDecoder(r.Body).Decode(
				&updatedTransaction,
			)

			if err != nil {
				http.Error(
					w,
					"Invalid JSON",
					http.StatusBadRequest,
				)
				return
			}

			// Validate the new transaction values.
			validationError := utils.ValidateTransaction(
				updatedTransaction,
			)

			if validationError != "" {
				http.Error(
					w,
					validationError,
					http.StatusBadRequest,
				)
				return
			}

			// Update the transaction in PostgreSQL.
			result, err := storage.UpdateTransactionInDB(
				pool,
				id,
				updatedTransaction,
			)

			if err != nil {
				http.Error(
					w,
					"Transaction not found",
					http.StatusNotFound,
				)
				return
			}

			// Return the updated transaction.
			json.NewEncoder(w).Encode(result)
			return
		}
		if r.Method == http.MethodDelete {

			parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

			if len(parts) != 2 {
				http.Error(w, "Invalid transaction path", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(parts[1])

			if err != nil {
				http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
				return
			}

			handleDeleteTransaction(w, r, pool, id)
			return
		}

		// Reject unsupported HTTP methods.
		http.Error(
			w,
			"Method not allowed",
			http.StatusMethodNotAllowed,
		)
	}
}

// handleDeleteTransaction deletes a transaction from PostgreSQL.
func handleDeleteTransaction(
	w http.ResponseWriter,
	r *http.Request,
	pool *pgxpool.Pool,
	id int,
) {
	// Delete the transaction from the database.
	err := storage.DeleteTransactionFromDB(pool, id)

	// If nothing was deleted, the transaction does not exist.
	if err != nil {
		http.Error(
			w,
			"Transaction not found",
			http.StatusNotFound,
		)
		return
	}

	// 204 means the deletion was successful
	// and there is no response body.
	w.Header().Del("Content-Type")
	w.WriteHeader(http.StatusNoContent)
}
