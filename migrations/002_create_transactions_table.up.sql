-- Create transactions table
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    transaction_code VARCHAR(50) UNIQUE NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    
    -- Amounts
    amount DECIMAL(12,2) NOT NULL,
    fee DECIMAL(12,2) DEFAULT 0,
    total_amount DECIMAL(12,2) NOT NULL,
    
    -- Status flow: waiting_payment -> pending -> validated -> processing -> shipped -> completed
    status VARCHAR(30) DEFAULT 'waiting_payment' CHECK (
        status IN ('waiting_payment', 'pending', 'validated', 'processing', 'shipped', 'completed', 'cancelled', 'disputed')
    ),
    
    -- Foreign keys
    buyer_id INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    seller_id INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    admin_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    
    -- Payment info
    payment_proof_url VARCHAR(255),
    payment_verified_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    payment_verified_at TIMESTAMP,
    
    -- Shipping info (data yang dikirim penjual)
    shipping_data TEXT,
    shipping_sent_at TIMESTAMP,
    
    -- Timestamps for status changes
    paid_at TIMESTAMP,
    validated_at TIMESTAMP,
    processing_at TIMESTAMP,
    shipped_at TIMESTAMP,
    completed_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    
    -- Auto complete deadline (2x24 jam setelah shipped)
    auto_complete_at TIMESTAMP,
    
    -- General timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Create indexes
CREATE INDEX idx_transactions_code ON transactions(transaction_code);
CREATE INDEX idx_transactions_buyer_id ON transactions(buyer_id);
CREATE INDEX idx_transactions_seller_id ON transactions(seller_id);
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transactions_created_at ON transactions(created_at);
CREATE INDEX idx_transactions_auto_complete_at ON transactions(auto_complete_at);

-- Create trigger for updated_at
CREATE TRIGGER update_transactions_updated_at
    BEFORE UPDATE ON transactions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();