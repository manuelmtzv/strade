CREATE TABLE cities (
    id VARCHAR(10) PRIMARY KEY,           
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    state_id CHAR(2) REFERENCES states(id)
);

COMMENT ON COLUMN cities.id IS 'Composite ID: state_code + city_code (or special codes like 000)';
COMMENT ON COLUMN cities.name IS 'City name';

CREATE INDEX idx_cities_name ON cities(name);
CREATE INDEX idx_cities_slug ON cities(slug);

INSERT INTO cities (id, name, slug) VALUES ('000', 'Sin ciudad', 'sin-ciudad');