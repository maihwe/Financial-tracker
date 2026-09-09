# Financial Tracker

A full-stack financial management application for tracking income and expenses, calculating total income, total expenses, and current balance.

## Overview

Financial Tracker is a web application built with **Go**, **PostgreSQL**, **HTML**, **CSS**, and **JavaScript**.

The application allows users to record financial transactions and manage them through a simple interface.

The financial balance is calculated as:

```text
Balance = Total Income - Total Expenses
```

The project was built to practice real-world backend development, database integration, REST API design, frontend development, CRUD operations, validation, testing, and Git/GitHub workflows.

---

## Features

* Add income transactions
* Add expense transactions
* View transaction history
* Edit existing transactions
* Delete transactions
* Categorize transactions
* Display transaction type
* Display transaction creation date
* Automatically calculate total income
* Automatically calculate total expenses
* Automatically calculate current balance
* Nigerian Naira (₦) currency formatting
* PostgreSQL database persistence
* Input validation
* REST API
* Automated backend tests

---

## Technology Stack

### Backend

* Go
* `net/http`
* PostgreSQL
* pgx/v5

### Frontend

* HTML5
* CSS3
* JavaScript
* Fetch API

### Database

* PostgreSQL
* Neon PostgreSQL

### Development Tools

* Git
* GitHub
* VS Code
* Linux terminal

---

## Project Structure

```text
Financial-tracker/
│
├── main.go
├── go.mod
│
├── models/
│   └── transaction.go
│
├── storage/
│   ├── memory.go
│   ├── postgres.go
│   └── postgres_test.go
│
├── handlers/
│   └── transaction_handler.go
│
├── utils/
│   └── validation.go
│
├── database/
│   ├── connection.go
│   ├── cmd/
│   │   └── migrate.go
│   │
│   └── migrations/
│       ├── 001_create_transactions.sql
│       └── 002_add_transaction_timestamps.sql
│
└── frontend/
    ├── index.html
    ├── style.css
    └── app.js
```

---

## How It Works

The application follows a simple client-server architecture.

```text
Browser
   │
   │ HTTP Request
   ▼
Go HTTP Server
   │
   ▼
Transaction Handlers
   │
   ▼
Storage Layer
   │
   ▼
PostgreSQL Database
```

When a user creates, updates, retrieves, or deletes a transaction, the frontend communicates with the Go backend through HTTP requests.

The backend processes the request and communicates with PostgreSQL for persistent data storage.

---

## Transaction Model

Each transaction contains:

```text
ID
Title
Amount
Category
Type
CreatedAt
UpdatedAt
```

The `type` field determines whether a transaction is an:

```text
income
```

or:

```text
expense
```

---

## Financial Calculation

The application calculates the financial summary from the transaction data.

```text
Total Income
     -
Total Expenses
     =
Balance
```

For example:

```text
Income:    ₦100,000
Expenses:   ₦25,000
-------------------
Balance:    ₦75,000
```

---

## API Endpoints

### Get all transactions

```http
GET /transactions
```

Returns all transactions.

---

### Get one transaction

```http
GET /transactions/{id}
```

Returns a single transaction by ID.

---

### Create a transaction

```http
POST /transactions
```

Example request body:

```json
{
    "title": "Salary",
    "amount": 100000,
    "category": "Work",
    "type": "income"
}
```

---

### Update a transaction

```http
PUT /transactions/{id}
```

Example request body:

```json
{
    "title": "Updated Salary",
    "amount": 120000,
    "category": "Work",
    "type": "income"
}
```

The `updated_at` timestamp is changed when a transaction is updated, while the original `created_at` timestamp is preserved.

---

### Delete a transaction

```http
DELETE /transactions/{id}
```

Deletes the specified transaction.

A successful deletion returns:

```http
204 No Content
```

---

## Validation

The backend validates transaction data before storing it.

A transaction must have:

* A title
* An amount greater than zero
* A category
* A valid transaction type

Valid transaction types are:

```text
income
expense
```

---

## Database

The project uses PostgreSQL for persistent storage.

Database migrations are located in:

```text
database/migrations/
```

The migrations create the transactions table and add transaction timestamps.

The application reads the PostgreSQL connection string from the environment variable:

```text
DATABASE_URL
```

The database connection string should **never be committed to GitHub**.

---

## Running the Project

### 1. Clone the repository

```bash
git clone https://github.com/maihwe/Financial-tracker.git
```

### 2. Enter the project

```bash
cd Financial-tracker
```

### 3. Install dependencies

```bash
go mod download
```

### 4. Set the database connection

Set your PostgreSQL connection string as the `DATABASE_URL` environment variable.

Do not place your real database credentials directly inside the source code.

### 5. Run database migrations

Run the migration command provided by the project.

### 6. Start the application

```bash
go run .
```

The application runs on:

```text
http://localhost:8080
```

Open the address in a web browser.

---

## Testing

The Go test suite can be run with:

```bash
go test ./...
```

The project currently includes tests for the PostgreSQL storage layer.

A successful test run should report the storage package as passing.

---

## Git Workflow

The project was developed using Git for version control.

Typical workflow:

```bash
git status

git add .

git commit -m "describe your changes"

git push origin main
```

The final frontend CRUD implementation was committed with:

```text
Complete financial tracker frontend CRUD
```

Commit:

```text
c33ea25
```

---

## What I Learned

Building this project provided practical experience with:

* Go project organization
* HTTP servers
* REST APIs
* CRUD operations
* PostgreSQL
* Database migrations
* SQL
* Environment variables
* Backend validation
* JSON
* JavaScript Fetch API
* DOM manipulation
* HTML forms
* Frontend/backend communication
* Automated testing
* Git and GitHub
* Debugging full-stack applications

---

## Future Improvements

Possible future improvements include:

* User authentication
* User-specific financial data
* Transaction search
* Filtering by category
* Filtering by income/expense
* Date-based filtering
* Financial reports
* Charts and visualizations
* Monthly spending summaries
* Export transactions to CSV
* Deployment to a cloud platform

These features are intentionally left as future improvements so the current project remains focused on its core financial tracking functionality.

---

## Project Status

**Completed — Version 1**

The core financial tracking functionality is implemented and tested.

```text
Backend       ✅
Database      ✅
API           ✅
Frontend      ✅
CRUD          ✅
Validation    ✅
Testing       ✅
GitHub        ✅
Documentation ⏳
```

---

## Author

**Maihwe**

GitHub:

https://github.com/maihwe
