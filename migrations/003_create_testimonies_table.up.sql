-- Create testimonies table
CREATE TABLE testimonies (
    id SERIAL PRIMARY KEY,
    transaction_id INTEGER UNIQUE NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    admin_notes TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Ensure one testimony per transaction
    UNIQUE(transaction_id, user_id)
);

-- Create indexes
CREATE INDEX idx_testimonies_transaction_id ON testimonies(transaction_id);
CREATE INDEX idx_testimonies_user_id ON testimonies(user_id);
CREATE INDEX idx_testimonies_rating ON testimonies(rating);

-- Create trigger for updated_at
CREATE TRIGGER update_testimonies_updated_at
    BEFORE UPDATE ON testimonies
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();