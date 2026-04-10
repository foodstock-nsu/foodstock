CREATE TYPE transaction_status AS ENUM('PENDING', 'SUCCESS', 'FAILED'); 

CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    order_id UUID REFERENCES orders(id),           
    sbp_transaction_id TEXT,           
    amount NUMERIC(10, 2),   
    status transaction_status,         
    webhook_received_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL 
        DEFAULT CURRENT_TIMESTAMP
); 