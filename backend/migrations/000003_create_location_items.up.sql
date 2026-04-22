CREATE TABLE IF NOT EXISTS location_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID REFERENCES items(id),
    location_id UUID REFERENCES locations(id),       
    price NUMERIC(10, 2) NOT NULL CHECK (price >= 0),           
    is_available BOOLEAN NOT NULL DEFAULT TRUE,   
    stock_amount INT NOT NULL CHECK (stock_amount >= 0)
); 
