CREATE TABLE locations (
    id UUID PRIMARY KEY,
    slug TEXT UNIQUE,           
    name TEXT NOT NULL,           
    address TEXT NOT NULL,   
    is_active BOOLEAN NOT NULL DEFAULT TRUE,         
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
); 

CREATE TYPE item_category AS ENUM('lunch','breakfast','drinks');

CREATE TABLE items (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,           
    description TEXT NOT NULL,           
    category item_category NOT NULL,   
    photo_url TEXT,         
    calories INT CHECK (calories >= 0), 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
); 

CREATE TABLE inventory (
    id UUID PRIMARY KEY,
    product_id UUID REFERENCES items(id),     
    location_id UUID REFERENCES locations(id),       
    price NUMERIC(10, 2) NOT NULL CHECK (price >= 0),           
    is_available BOOLEAN NOT NULL DEFAULT TRUE,   
    stock_amount INT NOT NULL CHECK (stock_amount >= 0)
); 

CREATE TYPE order_status AS ENUM('PENDING', 'PAID', 'CANCELED');

CREATE TABLE orders (
    id UUID PRIMARY KEY, 
    location_id UUID REFERENCES locations(id),           
    status order_status NOT NULL,           
    total_price NUMERIC(10, 2) CHECK (total_price >= 0),   
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,         
    paid_at TIMESTAMP WITH TIME ZONE
); 

CREATE TABLE order_items (
    id UUID PRIMARY KEY,
    order_id UUID REFERENCES orders(id),           
    item_id UUID REFERENCES items(id),           
    item_amount INT NOT NULL,   
    price_at_purchase NUMERIC(10, 2) NOT NULL 
        CHECK (price_at_purchase >= 0)        
); 

CREATE TYPE transaction_status AS ENUM('PENDING', 'SUCCESS', 'FAILED'); 

CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    order_id UUID REFERENCES orders(id),           
    sbp_transaction_id TEXT,           
    amount NUMERIC(10, 2),   
    status transaction_status,         
    webhook_recieved_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL 
        DEFAULT CURRENT_TIMESTAMP
); 