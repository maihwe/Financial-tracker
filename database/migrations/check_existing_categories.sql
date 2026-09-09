SELECT
    category,
    COUNT(*) AS transaction_count
FROM transactions
GROUP BY category
ORDER BY category;

