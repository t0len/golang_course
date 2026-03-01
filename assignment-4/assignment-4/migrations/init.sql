-- Create movies table on startup
CREATE TABLE IF NOT EXISTS movies (
    id       SERIAL PRIMARY KEY,
    title    VARCHAR(255) NOT NULL,
    genre    VARCHAR(100) NOT NULL,
    budget   INTEGER      NOT NULL,
    hero     VARCHAR(255) NOT NULL,
    heroine  VARCHAR(255) NOT NULL
);

-- Seed some initial data
INSERT INTO movies (title, genre, budget, hero, heroine)
VALUES 
    ('SAW',  'Horror',  500000,  'JONNY DEPP', 'Scarlet'),
    ('TEST', 'Romance', 1000000, 'BALE',        'ARMAS')
ON CONFLICT DO NOTHING;
