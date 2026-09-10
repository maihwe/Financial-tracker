// ================================
// ADMIN SYSTEM STATISTICS
// ================================

async function loadStatistics() {

    const statisticsContainer =
        document.querySelector("#statistics-container");

    try {

        const response =
            await fetch("/admin/statistics");

        if (!response.ok) {

            statisticsContainer.innerHTML =
                "<p>Failed to load statistics.</p>";

            return;
        }

        const statistics =
            await response.json();

        statisticsContainer.innerHTML = `
            <div class="stat-card">
                <h3>Total Users</h3>
                <p>${statistics.total_users}</p>
            </div>

            <div class="stat-card">
                <h3>Total Transactions</h3>
                <p>${statistics.total_transactions}</p>
            </div>

            <div class="stat-card">
                <h3>Total Income</h3>
                <p>₦${statistics.total_income.toLocaleString()}</p>
            </div>

            <div class="stat-card">
                <h3>Total Expenses</h3>
                <p>₦${statistics.total_expenses.toLocaleString()}</p>
            </div>

            <div class="stat-card">
                <h3>System Balance</h3>
                <p>₦${statistics.balance.toLocaleString()}</p>
            </div>
        `;

    } catch (error) {

        console.error(
            "Error loading statistics:",
            error
        );

        statisticsContainer.innerHTML =
            "<p>Unable to load statistics.</p>";
    }
}

loadStatistics();