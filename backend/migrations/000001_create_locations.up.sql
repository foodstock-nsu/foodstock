CREATE TABLE IF NOT EXISTS locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug TEXT NOT NULL,
    name TEXT NOT NULL,           
    address TEXT NOT NULL,   
    is_active BOOLEAN NOT NULL DEFAULT TRUE,         
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX IF NOT EXISTS locations_slug_unique_idx ON locations (slug)
    WHERE deleted_at IS NULL;