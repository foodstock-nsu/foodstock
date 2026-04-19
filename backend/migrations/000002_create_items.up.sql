CREATE TYPE item_category AS ENUM(
    'lunch','breakfast','drinks',
    'snacks','desserts','other'
);

CREATE TABLE items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,           
    description TEXT,           
    category item_category NOT NULL,   
    photo_url TEXT NOT NULL,
    calories INT CHECK (calories >= 0),
    proteins NUMERIC(6, 2) CHECK (proteins >= 0),
    fats NUMERIC(6, 2) CHECK (fats >= 0),
    carbs NUMERIC(6, 2) CHECK (carbs >= 0),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);