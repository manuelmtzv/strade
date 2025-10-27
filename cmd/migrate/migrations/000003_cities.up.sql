CREATE TABLE cities (
    id CHAR(5) PRIMARY KEY,           
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL        
);

COMMENT ON COLUMN cities.id IS 'SEPOMEX city code';
COMMENT ON COLUMN cities.name IS 'City name';

CREATE INDEX idx_cities_name ON cities(name);
CREATE INDEX idx_cities_slug ON cities(slug);