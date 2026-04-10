CREATE TABLE order_items (
    id UUID PRIMARY KEY,
    order_id UUID REFERENCES orders(id),           
    item_id UUID REFERENCES items(id),           
    item_amount INT NOT NULL,   
    price_at_purchase NUMERIC(10, 2) NOT NULL 
        CHECK (price_at_purchase >= 0)        
); 