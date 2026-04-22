CREATE TABLE IF NOT EXISTS locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug TEXT UNIQUE NOT NULL,           
    name TEXT NOT NULL,           
    address TEXT NOT NULL,   
    is_active BOOLEAN NOT NULL DEFAULT TRUE,         
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
); 