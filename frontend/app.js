// Get the transaction list from our Go API.
async function loadTransactions() {

    // Ask the Go server for all transactions.
    const response = await fetch("http://localhost:8080/transactions");

    // Convert the response from JSON into JavaScript data.
    const transactions = await response.json();

    // Find the HTML element where transactions should appear.
    const transactionList = document.querySelector("#transaction-list");

    let totalIncome = 0;

    transactions.forEach(function(transaction) {

    if (transaction.type === "income") {
        totalIncome += transaction.amount;
    }

        const item = document.createElement("div");

        item.textContent =
            transaction.title + " - ₦" + transaction.amount;

        transactionList.appendChild(item);

    });
}

// Load transactions when the page opens.
loadTransactions();


// Find the transaction form on the page.
const form = document.querySelector("form");

// Listen for the form being submitted.
form.addEventListener("submit", async function(event) {

    // Stop the browser from refreshing the page.
    event.preventDefault();

    // Collect the values from the form.
    const transaction = {
        title: form.elements.title.value,
        amount: Number(form.elements.amount.value),
        category: form.elements.category.value,
        type: form.elements.type.value
    };

    // Send the transaction to the Go API.
    const response = await fetch("http://localhost:8080/transactions", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(transaction)
    });

    // Check what the Go server returned.
    const savedTransaction = await response.json();

    console.log("Saved transaction:", savedTransaction);

});