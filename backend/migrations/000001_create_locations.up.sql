CREATE TABLE locations (
    id UUID PRIMARY KEY,
    slug TEXT UNIQUE,           
    name TEXT NOT NULL,           
    address TEXT NOT NULL,   
    is_active BOOLEAN NOT NULL DEFAULT TRUE,         
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
); 