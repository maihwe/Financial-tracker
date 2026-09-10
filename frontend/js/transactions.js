// ================================
// TRANSACTION MANAGEMENT
// ================================


// Find the transaction form.
const form = document.querySelector("#transaction-form");

// Find the transaction list.
const transactionList = document.querySelector("#transaction-list");

// Find the form heading.
const formTitle = document.querySelector("#form-title");

// Find the submit button.
const submitButton = document.querySelector("#submit-button");

// Find the cancel button.
const cancelButton = document.querySelector("#cancel-button");

// Store the ID of the transaction currently being edited.
let editingTransactionID = null;


// Load all transactions from the Go API.
async function loadTransactions() {

    try {

        // Ask the Go server for all transactions.
        const response = await fetch(
            "/transactions"
        );

        // Check whether the request succeeded.
        if (!response.ok) {

            console.error(
                "Failed to load transactions."
            );

            return;
        }

        // Convert the response into JavaScript data.
        const transactions =
            await response.json();

        // Clear the existing transaction list.
        transactionList.innerHTML = "";


        // Start financial totals at zero.
        let totalIncome = 0;
        let totalExpenses = 0;


        // Go through every transaction.
        transactions.forEach(function(transaction) {

            // Calculate income.
            if (transaction.type === "income") {

                totalIncome +=
                    transaction.amount;
            }

            // Calculate expenses.
            if (transaction.type === "expense") {

                totalExpenses +=
                    transaction.amount;
            }


            // Create a transaction card.
            const item =
                document.createElement("div");

            item.className =
                "transaction-item";


            // Format transaction date.
            const formattedDate =
                new Date(
                    transaction.transaction_at
                ).toLocaleString("en-NG");


            // Get the category name from its ID.
            const categoryName =
                getCategoryName(
                    transaction.category_id
                );


            // Display the transaction.
            item.innerHTML = `
                <h3>${transaction.title}</h3>

                <p>
                    Amount:
                    ${formatCurrency(transaction.amount)}
                </p>

                <p>
                    Category:
                    ${categoryName}
                </p>

                <p>
                    Type:
                    ${transaction.type}
                </p>

                <p>
                    Date:
                    ${formattedDate}
                </p>

                <div class="transaction-buttons">

                    <button
                        class="edit-button"
                        data-id="${transaction.id}"
                    >
                        Edit
                    </button>

                    <button
                        class="delete-button"
                        data-id="${transaction.id}"
                    >
                        Delete
                    </button>

                </div>
            `;


            // Add the card to the page.
            transactionList.appendChild(item);

        });


        // Calculate balance.
        const balance =
            totalIncome - totalExpenses;


        // Display financial summary.
        document.querySelector(
            "#total-income"
        ).textContent =
            formatCurrency(totalIncome);

        document.querySelector(
            "#total-expenses"
        ).textContent =
            formatCurrency(totalExpenses);

        document.querySelector(
            "#balance"
        ).textContent =
            formatCurrency(balance);

    } catch (error) {

        console.error(
            "Error loading transactions:",
            error
        );
    }
}


// Handle clicks inside the transaction list.
transactionList.addEventListener(
    "click",
    async function(event) {

        // Check if Edit was clicked.
        if (
            event.target.classList.contains(
                "edit-button"
            )
        ) {

            // Get the transaction ID.
            const transactionID =
                event.target.dataset.id;


            // Get the transaction from the API.
            const response =
                await fetch(
                    `/transactions/${transactionID}`
                );


            // Check whether the request succeeded.
            if (!response.ok) {

                console.error(
                    "Failed to get transaction."
                );

                return;
            }


            // Convert response to JSON.
            const transaction =
                await response.json();


            // Remember which transaction we are editing.
            editingTransactionID =
                transaction.id;


            // Put the transaction values into the form.
            document.querySelector(
                "#title"
            ).value =
                transaction.title;

            document.querySelector(
                "#amount"
            ).value =
                transaction.amount;

            document.querySelector(
                "#category"
            ).value =
                getCategoryName(
                    transaction.category_id
                );

            document.querySelector(
                "#type"
            ).value =
                transaction.type;


            // Change the form into Edit mode.
            formTitle.textContent =
                "Edit Transaction";

            submitButton.textContent =
                "Update Transaction";

            cancelButton.hidden =
                false;


            // Move the user back to the form.
            form.scrollIntoView({
                behavior: "smooth"
            });
        }


        // Check if Delete was clicked.
        if (
            event.target.classList.contains(
                "delete-button"
            )
        ) {

            // Get the transaction ID.
            const transactionID =
                event.target.dataset.id;


            // Ask the server to delete the transaction.
            const response =
                await fetch(
                    `/transactions/${transactionID}`,
                    {
                        method: "DELETE"
                    }
                );


            // Check whether deletion succeeded.
            if (!response.ok) {

                console.error(
                    "Failed to delete transaction."
                );

                return;
            }


            // Reload the transaction list.
            loadTransactions();
        }
    }
);


// Handle form submission.
form.addEventListener(
    "submit",
    async function(event) {

        // Prevent page refresh.
        event.preventDefault();


        // Collect form values.
        const transaction = {

            title:
                document.querySelector(
                    "#title"
                ).value,

            amount:
                Number(
                    document.querySelector(
                        "#amount"
                    ).value
                ),

            category_id:
                getCategoryID(
                    document.querySelector(
                        "#category"
                    ).value
                ),

            type:
                document.querySelector(
                    "#type"
                ).value
        };


        let response;


        // If no transaction is being edited,
        // create a new transaction.
        if (
            editingTransactionID === null
        ) {

            response =
                await fetch(
                    "/transactions",
                    {
                        method: "POST",

                        headers: {
                            "Content-Type":
                                "application/json"
                        },

                        body:
                            JSON.stringify(
                                transaction
                            )
                    }
                );
        }


        // Otherwise update an existing transaction.
        else {

            response =
                await fetch(
                    `/transactions/${editingTransactionID}`,
                    {
                        method: "PUT",

                        headers: {
                            "Content-Type":
                                "application/json"
                        },

                        body:
                            JSON.stringify(
                                transaction
                            )
                    }
                );
        }


        // Check whether saving succeeded.
        if (!response.ok) {

            console.error(
                "Failed to save transaction."
            );

            return;
        }


        // Reset the form.
        form.reset();


        // Reset edit mode.
        editingTransactionID =
            null;


        // Return form to Add mode.
        formTitle.textContent =
            "Add Transaction";

        submitButton.textContent =
            "Add Transaction";

        cancelButton.hidden =
            true;


        // Reload transactions and summary.
        loadTransactions();
    }
);


// Cancel Edit button.
cancelButton.addEventListener(
    "click",
    function() {

        // Reset the form.
        form.reset();


        // Leave edit mode.
        editingTransactionID =
            null;


        // Return form to Add mode.
        formTitle.textContent =
            "Add Transaction";

        submitButton.textContent =
            "Add Transaction";

        cancelButton.hidden =
            true;
    }
);
