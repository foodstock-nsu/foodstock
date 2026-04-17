CREATE TYPE transaction_status AS ENUM('PENDING', 'SUCCESS', 'FAILED'); 

CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    order_id UUID REFERENCES orders(id),           
    sbp_transaction_id TEXT NOT NULL,           
    amount NUMERIC(10, 2) NOT NULL,   
    status transaction_status DEFAULT 'PENDING',         
    webhook_received_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL 
        DEFAULT CURRENT_TIMESTAMP
); 