CREATE TYPE item_category AS ENUM('lunch','breakfast','drinks');

CREATE TABLE items (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,           
    description TEXT,           
    category item_category NOT NULL,   
    photo_url TEXT,         
    calories INT CHECK (calories >= 0), 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
); 