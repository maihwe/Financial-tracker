package main

import (
	"fmt"
	"net/http"

	"financial-tracker/handlers"
)

// main starts the HTTP server
// and connects our application routes.
func main() {

	// Register the transactions collection endpoint.
	// Example: GET /transactions
	//          POST /transactions
	http.HandleFunc("/transactions", handlers.TransactionHandler)

	// Register individual transaction endpoints.
	// Example: GET /transactions/1
	http.HandleFunc("/transactions/", handlers.TransactionHandler)

	// Register the financial summary endpoint.
	// Example: GET /summary
	http.HandleFunc("/summary", handlers.SummaryHandler)

	// Display the server address.
	fmt.Println("Server running on http://localhost:8080")

	// Start the HTTP server.
	err := http.ListenAndServe(":8080", nil)

	// Display an error if the server fails to start.
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
