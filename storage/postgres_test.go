package storage

import (
	"os"
	"testing"

	"financial-tracker/database"
	"financial-tracker/models"
)

// TestGetTransactionByIDFromDB checks that we can
// retrieve one transaction from PostgreSQL.
func TestGetTransactionByIDFromDB(t *testing.T) {

	// Skip the test if DATABASE_URL is not available.
	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("DATABASE_URL is not set")
	}

	// Connect to PostgreSQL.
	conn, err := database.Connect()
	if err != nil {
		t.Fatal(err)
	}

	defer conn.Close()

	// Create a transaction specifically for this test.
	transaction := models.Transaction{
		Title:    "Get Test",
		Amount:   10000,
		Category: "Testing",
		Type:     "income",
	}

	// Insert the test transaction.
	savedTransaction, err := AddTransactionToDB(
		conn,
		transaction,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Clean up the test transaction when the test finishes.
	defer DeleteTransactionFromDB(
		conn,
		savedTransaction.ID,
	)

	// Retrieve the transaction we just created.
	foundTransaction, err := GetTransactionByIDFromDB(
		conn,
		savedTransaction.ID,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Check that we retrieved the correct transaction.
	if foundTransaction.ID != savedTransaction.ID {
		t.Errorf(
			"expected ID %d, got %d",
			savedTransaction.ID,
			foundTransaction.ID,
		)
	}

	// Check that PostgreSQL created the timestamps.
	if foundTransaction.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}

	if foundTransaction.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

// TestUpdateTransactionInDB checks that we can
// update an existing transaction in PostgreSQL.
func TestUpdateTransactionInDB(t *testing.T) {

	// Skip the test if DATABASE_URL is not available.
	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("DATABASE_URL is not set")
	}

	// Connect to PostgreSQL.
	conn, err := database.Connect()
	if err != nil {
		t.Fatal(err)
	}

	defer conn.Close()

	// Create a transaction specifically for this test.
	transaction := models.Transaction{
		Title:    "Original Salary",
		Amount:   50000,
		Category: "Salary",
		Type:     "income",
	}

	// Insert the test transaction.
	savedTransaction, err := AddTransactionToDB(
		conn,
		transaction,
	)
	if err != nil {
		t.Fatal(err)
	}
	originalCreatedAt := savedTransaction.CreatedAt
	originalUpdatedAt := savedTransaction.UpdatedAt
	// Clean up the test transaction when the test finishes.
	defer DeleteTransactionFromDB(
		conn,
		savedTransaction.ID,
	)

	// Create the new values.
	updatedTransaction := models.Transaction{
		Title:    "Updated Salary",
		Amount:   75000,
		Category: "Salary",
		Type:     "income",
	}

	// Update the transaction.
	result, err := UpdateTransactionInDB(
		conn,
		savedTransaction.ID,
		updatedTransaction,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the ID did not change.
	if result.ID != savedTransaction.ID {
		t.Errorf(
			"expected ID %d, got %d",
			savedTransaction.ID,
			result.ID,
		)
	}

	// Check that the title changed.
	if result.Title != "Updated Salary" {
		t.Errorf(
			"expected title Updated Salary, got %s",
			result.Title,
		)
	}

	// Check that the amount changed.
	if result.Amount != 75000 {
		t.Errorf(
			"expected amount 75000, got %.2f",
			result.Amount,
		)
	}
	// created_at should never change when a transaction is updated.
	if !result.CreatedAt.Equal(originalCreatedAt) {
		t.Error("expected CreatedAt to remain unchanged")
	}

	if result.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}

	if !result.UpdatedAt.After(originalUpdatedAt) {
		t.Error("expected UpdatedAt to change after update")
	}

	// After an update, updated_at should be at or after created_at.
	if result.UpdatedAt.Before(result.CreatedAt) {
		t.Error("expected UpdatedAt to be at or after CreatedAt")
	}
}

// TestDeleteTransactionFromDB checks that we can
// delete a transaction from PostgreSQL.
func TestDeleteTransactionFromDB(t *testing.T) {

	// Skip the test if DATABASE_URL is not available.
	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("DATABASE_URL is not set")
	}

	// Connect to PostgreSQL.
	conn, err := database.Connect()
	if err != nil {
		t.Fatal(err)
	}

	defer conn.Close()

	// Create a transaction specifically for this test.
	transaction := models.Transaction{
		Title:    "Delete Test",
		Amount:   15000,
		Category: "Testing",
		Type:     "expense",
	}

	// Insert the test transaction.
	savedTransaction, err := AddTransactionToDB(
		conn,
		transaction,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Delete the transaction.
	err = DeleteTransactionFromDB(
		conn,
		savedTransaction.ID,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Try to retrieve the deleted transaction.
	_, err = GetTransactionByIDFromDB(
		conn,
		savedTransaction.ID,
	)

	// We expect PostgreSQL to report that the row doesn't exist.
	if err == nil {
		t.Fatal("expected transaction to be deleted")
	}
}

// TestGetFinancialSummaryFromDB checks that PostgreSQL
// correctly calculates income, expenses, and balance.
func TestGetFinancialSummaryFromDB(t *testing.T) {

	// Skip the test if DATABASE_URL is not available.
	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("DATABASE_URL is not set")
	}

	// Connect to PostgreSQL.
	conn, err := database.Connect()
	if err != nil {
		t.Fatal(err)
	}

	defer conn.Close()

	// Create test income.
	income := models.Transaction{
		Title:    "Summary Income Test",
		Amount:   100000,
		Category: "Testing",
		Type:     "income",
	}

	// Save the income transaction.
	savedIncome, err := AddTransactionToDB(conn, income)
	if err != nil {
		t.Fatal(err)
	}

	// Create test expense.
	expense := models.Transaction{
		Title:    "Summary Expense Test",
		Amount:   25000,
		Category: "Testing",
		Type:     "expense",
	}

	// Save the expense transaction.
	savedExpense, err := AddTransactionToDB(conn, expense)
	if err != nil {
		t.Fatal(err)
	}

	// Clean up both test transactions.
	defer DeleteTransactionFromDB(conn, savedIncome.ID)
	defer DeleteTransactionFromDB(conn, savedExpense.ID)

	// Get the financial summary.
	totalIncome, totalExpenses, balance, err :=
		GetFinancialSummaryFromDB(conn)

	if err != nil {
		t.Fatal(err)
	}

	// The database may already contain other transactions,
	// so we check that our expected amounts are reflected
	// in the totals.
	if totalIncome < 100000 {
		t.Errorf(
			"expected total income to be at least 100000, got %.2f",
			totalIncome,
		)
	}

	if totalExpenses < 25000 {
		t.Errorf(
			"expected total expenses to be at least 25000, got %.2f",
			totalExpenses,
		)
	}

	// Verify the balance formula.
	expectedBalance := totalIncome - totalExpenses

	if balance != expectedBalance {
		t.Errorf(
			"expected balance %.2f, got %.2f",
			expectedBalance,
			balance,
		)
	}
}
