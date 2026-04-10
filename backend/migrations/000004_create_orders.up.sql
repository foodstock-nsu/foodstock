CREATE TYPE order_status AS ENUM('PENDING', 'PAID', 'CANCELED');

CREATE TABLE orders (
    id UUID PRIMARY KEY, 
    location_id UUID REFERENCES locations(id),           
    status order_status NOT NULL,           
    total_price NUMERIC(10, 2) CHECK (total_price >= 0),   
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,         
    paid_at TIMESTAMP WITH TIME ZONE
); 