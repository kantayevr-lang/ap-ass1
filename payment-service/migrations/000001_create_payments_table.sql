CREATE TABLE IF NOT EXISTS payments (
    id VARCHAR(50) PRIMARY KEY,
    order_id VARCHAR(50) NOT NULL,
    transaction_id VARCHAR(50),
    amount BIGINT NOT NULL, 
    status VARCHAR(20) NOT NULL
);