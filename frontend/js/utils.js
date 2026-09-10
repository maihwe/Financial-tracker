// ================================
// SHARED UTILITY FUNCTIONS
// ================================


// Global categories.
const categoryIDs = {
    "Food": 1,
    "Transportation": 2,
    "Salary": 3,
    "Business": 4,
    "Education": 5,
    "Health": 6,
    "Bills & Utilities": 7,
    "Housing": 8,
    "Shopping": 9,
    "Entertainment": 10,
    "Savings": 11,
    "Debt": 12,
    "Investment": 13,
    "Other": 14,
    "Trade": 15,
    "Hygiene": 16
};


// Find a category ID from its name.
function getCategoryID(categoryName) {

    const name =
        categoryName.trim();

    return categoryIDs[name];
}


// Find a category name from its ID.
function getCategoryName(categoryID) {

    for (const categoryName in categoryIDs) {

        if (
            categoryIDs[categoryName] ===
            categoryID
        ) {

            return categoryName;
        }
    }

    return "Unknown";
}


// Format numbers as Nigerian Naira.
function formatCurrency(amount) {

    return new Intl.NumberFormat(
        "en-NG",
        {
            style: "currency",
            currency: "NGN"
        }
    ).format(amount);
}
