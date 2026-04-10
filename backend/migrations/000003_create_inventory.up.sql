CREATE TABLE inventory (
    id UUID PRIMARY KEY,
    product_id UUID REFERENCES items(id),     
    location_id UUID REFERENCES locations(id),       
    price NUMERIC(10, 2) NOT NULL CHECK (price >= 0),           
    is_available BOOLEAN NOT NULL DEFAULT TRUE,   
    stock_amount INT NOT NULL CHECK (stock_amount >= 0)
); 
