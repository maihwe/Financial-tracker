package storage

import (
	"os"
	"testing"
	"time"

	"financial-tracker/database"
	"financial-tracker/models"
)

// TestGetTransactionByIDFromDB checks that we can
// retrieve one transaction belonging to the correct user.
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

	// Create a test user.
	user := models.User{
		Email:        "get-test-" + time.Now().Format("20060102150405.000000000") + "@example.com",
		PasswordHash: "test-password-hash",
	}

	createdUser, err := CreateUserInDB(
		conn,
		user,
	)
	if err != nil {
		t.Fatal(err)
	}

	userID := createdUser.ID

	// Create a transaction specifically for this test.
	transaction := models.Transaction{
		UserID:      userID,
		Title:       "Get Test",
		Amount:      10000,
		CategoryID:  14, // Other
		Type:        "income",
		Description: "Transaction retrieval test",
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
		userID,
	)

	// Retrieve the transaction using the correct user ID.
	foundTransaction, err := GetTransactionByIDFromDB(
		conn,
		savedTransaction.ID,
		userID,
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

	// Check that the correct user owns the transaction.
	if foundTransaction.UserID != userID {
		t.Errorf(
			"expected UserID %d, got %d",
			userID,
			foundTransaction.UserID,
		)
	}

	// Check that the category ID was retrieved correctly.
	if foundTransaction.CategoryID != 14 {
		t.Errorf(
			"expected CategoryID 14, got %d",
			foundTransaction.CategoryID,
		)
	}

	// Check that the description was retrieved correctly.
	if foundTransaction.Description != "Transaction retrieval test" {
		t.Errorf(
			"expected description Transaction retrieval test, got %s",
			foundTransaction.Description,
		)
	}

	// Check that PostgreSQL created the timestamps.
	if foundTransaction.TransactionAt.IsZero() {
		t.Error("expected TransactionAt to be set")
	}

	if foundTransaction.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}

	if foundTransaction.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

// TestGetTransactionByIDFromDBRejectsWrongUser verifies that
// one user cannot retrieve another user's transaction.
func TestGetTransactionByIDFromDBRejectsWrongUser(t *testing.T) {

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

	// Create the transaction owner.
	owner := models.User{
		Email:        "owner-test-" + time.Now().Format("20060102150405.000000000") + "@example.com",
		PasswordHash: "test-password-hash",
	}

	createdOwner, err := CreateUserInDB(
		conn,
		owner,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Create another user.
	otherUser := models.User{
		Email:        "other-test-" + time.Now().Format("20060102150405.000000000") + "@example.com",
		PasswordHash: "test-password-hash",
	}

	createdOtherUser, err := CreateUserInDB(
		conn,
		otherUser,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Create a transaction owned by the first user.
	transaction := models.Transaction{
		UserID:      createdOwner.ID,
		Title:       "Private Transaction",
		Amount:      5000,
		CategoryID:  14,
		Type:        "expense",
		Description: "Ownership test",
	}

	savedTransaction, err := AddTransactionToDB(
		conn,
		transaction,
	)
	if err != nil {
		t.Fatal(err)
	}

	defer DeleteTransactionFromDB(
		conn,
		savedTransaction.ID,
		createdOwner.ID,
	)

	// Try to retrieve the transaction using another user's ID.
	_, err = GetTransactionByIDFromDB(
		conn,
		savedTransaction.ID,
		createdOtherUser.ID,
	)

	// The query should not reveal the transaction.
	if err == nil {
		t.Fatal("expected wrong user to be unable to retrieve transaction")
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

	// Create a test user.
	user := models.User{
		Email:        "update-test-" + time.Now().Format("20060102150405.000000000") + "@example.com",
		PasswordHash: "test-password-hash",
	}

	createdUser, err := CreateUserInDB(
		conn,
		user,
	)
	if err != nil {
		t.Fatal(err)
	}

	userID := createdUser.ID

	// Create a transaction specifically for this test.
	transaction := models.Transaction{
		UserID:      userID,
		Title:       "Original Salary",
		Amount:      50000,
		CategoryID:  3, // Salary
		Type:        "income",
		Description: "Original salary transaction",
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
		userID,
	)

	updatedTransaction := models.Transaction{
		UserID:        userID,
		Title:         "Updated Salary",
		Amount:        75000,
		CategoryID:    3, // Salary
		Type:          "income",
		Description:   "Updated salary transaction",
		TransactionAt: savedTransaction.TransactionAt.Add(2 * time.Hour),
	}

	// Update the transaction.
	result, err := UpdateTransactionInDB(
		conn,
		savedTransaction.ID,
		userID,
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

	// Check that ownership did not change.
	if result.UserID != userID {
		t.Errorf(
			"expected UserID %d, got %d",
			userID,
			result.UserID,
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

	// Check that the description changed.
	if result.Description != "Updated salary transaction" {
		t.Errorf(
			"expected description Updated salary transaction, got %s",
			result.Description,
		)
	}

	// Check that the category ID is correct.
	if result.CategoryID != 3 {
		t.Errorf(
			"expected CategoryID 3, got %d",
			result.CategoryID,
		)
	}

	// TransactionAt should change because the update
	// supplied a new transaction time.
	if result.TransactionAt.Equal(savedTransaction.TransactionAt) {
		t.Error("expected TransactionAt to change")
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

	// Create a test user.
	user := models.User{
		Email:        "delete-test-" + time.Now().Format("20060102150405.000000000") + "@example.com",
		PasswordHash: "test-password-hash",
	}

	createdUser, err := CreateUserInDB(
		conn,
		user,
	)
	if err != nil {
		t.Fatal(err)
	}

	userID := createdUser.ID

	// Create a transaction specifically for this test.
	transaction := models.Transaction{
		UserID:      userID,
		Title:       "Delete Test",
		Amount:      15000,
		CategoryID:  14, // Other
		Type:        "expense",
		Description: "Transaction deletion test",
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
		userID,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Try to retrieve the deleted transaction.
	_, err = GetTransactionByIDFromDB(
		conn,
		savedTransaction.ID,
		userID,
	)

	// We expect PostgreSQL to report that the row doesn't exist.
	if err == nil {
		t.Fatal("expected transaction to be deleted")
	}
}

// TestDeleteTransactionFromDBRejectsWrongUser verifies that
// one user cannot delete another user's transaction.
func TestDeleteTransactionFromDBRejectsWrongUser(t *testing.T) {

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

	// Create transaction owner.
	owner := models.User{
		Email:        "delete-owner-" + time.Now().Format("20060102150405.000000000") + "@example.com",
		PasswordHash: "test-password-hash",
	}

	createdOwner, err := CreateUserInDB(
		conn,
		owner,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Create another user.
	otherUser := models.User{
		Email:        "delete-other-" + time.Now().Format("20060102150405.000000000") + "@example.com",
		PasswordHash: "test-password-hash",
	}

	createdOtherUser, err := CreateUserInDB(
		conn,
		otherUser,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Create a transaction owned by the first user.
	transaction := models.Transaction{
		UserID:     createdOwner.ID,
		Title:      "Protected Transaction",
		Amount:     1000,
		CategoryID: 14,
		Type:       "expense",
	}

	savedTransaction, err := AddTransactionToDB(
		conn,
		transaction,
	)
	if err != nil {
		t.Fatal(err)
	}

	defer DeleteTransactionFromDB(
		conn,
		savedTransaction.ID,
		createdOwner.ID,
	)

	// Try to delete it using another user's ID.
	err = DeleteTransactionFromDB(
		conn,
		savedTransaction.ID,
		createdOtherUser.ID,
	)

	if err != nil {
		t.Fatal(err)
	}

	// Verify the transaction still exists.
	_, err = GetTransactionByIDFromDB(
		conn,
		savedTransaction.ID,
		createdOwner.ID,
	)

	if err != nil {
		t.Fatal("expected transaction to still exist after unauthorized delete")
	}
}

// TestGetFinancialSummaryFromDB checks that PostgreSQL
// correctly calculates income, expenses, and balance
// for one specific user.
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

	// Create a test user.
	user := models.User{
		Email:        "summary-test-" + time.Now().Format("20060102150405.000000000") + "@example.com",
		PasswordHash: "test-password-hash",
	}

	createdUser, err := CreateUserInDB(
		conn,
		user,
	)
	if err != nil {
		t.Fatal(err)
	}

	userID := createdUser.ID

	// Create test income.
	income := models.Transaction{
		UserID:      userID,
		Title:       "Summary Income Test",
		Amount:      100000,
		CategoryID:  14, // Other
		Type:        "income",
		Description: "Summary income test",
	}

	// Save the income transaction.
	savedIncome, err := AddTransactionToDB(
		conn,
		income,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Create test expense.
	expense := models.Transaction{
		UserID:      userID,
		Title:       "Summary Expense Test",
		Amount:      25000,
		CategoryID:  14, // Other
		Type:        "expense",
		Description: "Summary expense test",
	}

	// Save the expense transaction.
	savedExpense, err := AddTransactionToDB(
		conn,
		expense,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Clean up both test transactions.
	defer DeleteTransactionFromDB(
		conn,
		savedIncome.ID,
		userID,
	)

	defer DeleteTransactionFromDB(
		conn,
		savedExpense.ID,
		userID,
	)

	// Get the financial summary for this user only.
	totalIncome, totalExpenses, balance, err :=
		GetFinancialSummaryFromDB(
			conn,
			userID,
		)

	if err != nil {
		t.Fatal(err)
	}

	// Because this test user is newly created,
	// these totals should contain exactly our test transactions.
	if totalIncome != 100000 {
		t.Errorf(
			"expected total income 100000, got %.2f",
			totalIncome,
		)
	}

	if totalExpenses != 25000 {
		t.Errorf(
			"expected total expenses 25000, got %.2f",
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
