CREATE TYPE transaction_status AS ENUM('PENDING', 'SUCCESS', 'FAILED'); 

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID REFERENCES orders(id),           
    sbp_transaction_id TEXT NOT NULL,           
    amount NUMERIC(10, 2) NOT NULL,   
    status transaction_status NOT NULL DEFAULT 'PENDING',
    paid_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL 
        DEFAULT CURRENT_TIMESTAMP
); 