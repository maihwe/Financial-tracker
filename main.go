package main

import (
	"fmt"
	"net/http"

	"financial-tracker/database"
	"financial-tracker/handlers"
)

func main() {

	// Connect to the PostgreSQL database.
	pool, err := database.Connect()

	// Stop the application if the database connection fails.
	if err != nil {
		fmt.Println("Database connection failed:", err)
		return
	}

	// Close the connection pool when the server stops.
	defer pool.Close()

	// Create our transaction handler using the database pool.
	transactionHandler := handlers.TransactionHandler(pool)

	// Register the transaction routes.
	http.HandleFunc("/transactions", transactionHandler)
	http.HandleFunc("/transactions/", transactionHandler)

	fmt.Println("Server running on http://localhost:8080")

	// Start the HTTP server.
	err = http.ListenAndServe(":8080", nil)

	// Report a server error if one occurs.
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
