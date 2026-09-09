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

		// Every transaction request requires authentication.
		userID, err := GetAuthenticatedUserID(r)

		if err != nil {
			http.Error(
				w,
				"Authentication required",
				http.StatusUnauthorized,
			)
			return
		}

		// Handle GET requests.
		if r.Method == http.MethodGet {

			// GET /transactions
			if r.URL.Path == "/transactions" {

				transactions, err :=
					storage.GetAllTransactionsFromDB(
						pool,
						userID,
					)

				if err != nil {
    				http.Error(
        				w,
        				"Could not retrieve transactions: "+err.Error(),
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
			if len(parts) != 2 ||
				parts[0] != "transactions" {

				http.Error(
					w,
					"Invalid transaction path",
					http.StatusBadRequest,
				)
				return
			}

			id, err := strconv.Atoi(parts[1])

			if err != nil {
				http.Error(
					w,
					"Invalid transaction ID",
					http.StatusBadRequest,
				)
				return
			}

			transaction, err :=
				storage.GetTransactionByIDFromDB(
					pool,
					id,
					userID,
				)

			if err != nil {
				http.Error(
					w,
					"Transaction not found",
					http.StatusNotFound,
				)
				return
			}

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

			// The authenticated user's ID comes from
			// the session, not from the frontend.
			newTransaction.UserID = userID

			validationError :=
				utils.ValidateTransaction(
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

			parts := strings.Split(
				strings.Trim(r.URL.Path, "/"),
				"/",
			)

			if len(parts) != 2 ||
				parts[0] != "transactions" {

				http.Error(
					w,
					"Invalid transaction path",
					http.StatusBadRequest,
				)
				return
			}

			id, err := strconv.Atoi(parts[1])

			if err != nil {
				http.Error(
					w,
					"Invalid transaction ID",
					http.StatusBadRequest,
				)
				return
			}

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

			// Never trust a UserID supplied by the client.
			// The authenticated user's ID is authoritative.
			updatedTransaction.UserID = userID

			validationError :=
				utils.ValidateTransaction(
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

			result, err :=
				storage.UpdateTransactionInDB(
					pool,
					id,
					userID,
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

			json.NewEncoder(w).Encode(result)
			return
		}

		// Handle DELETE requests.
		if r.Method == http.MethodDelete {

			parts := strings.Split(
				strings.Trim(r.URL.Path, "/"),
				"/",
			)

			if len(parts) != 2 ||
				parts[0] != "transactions" {

				http.Error(
					w,
					"Invalid transaction path",
					http.StatusBadRequest,
				)
				return
			}

			id, err := strconv.Atoi(parts[1])

			if err != nil {
				http.Error(
					w,
					"Invalid transaction ID",
					http.StatusBadRequest,
				)
				return
			}

			err =
				storage.DeleteTransactionFromDB(
					pool,
					id,
					userID,
				)

			if err != nil {
				http.Error(
					w,
					"Transaction not found",
					http.StatusNotFound,
				)
				return
			}

			w.Header().Del("Content-Type")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		http.Error(
			w,
			"Method not allowed",
			http.StatusMethodNotAllowed,
		)
	}
}
