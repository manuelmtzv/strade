CREATE TABLE settlements (
    id CHAR(9) PRIMARY KEY,
    postal_code CHAR(5) NOT NULL,                   
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    settlement_type_id CHAR(2) NOT NULL REFERENCES settlement_types(id), 
    municipality_id CHAR(5) NOT NULL REFERENCES municipalities(id),      
    city_id VARCHAR(10) NOT NULL REFERENCES cities(id),                     
    state_id CHAR(2) NOT NULL REFERENCES states(id),                     
    office_postal_code CHAR(5),                     
    zone VARCHAR(10),                               
    latitude DOUBLE PRECISION,                      
    longitude DOUBLE PRECISION
);

COMMENT ON COLUMN settlements.id IS 'Composite key: StateCode + MunicipalityCode + SettlementID (SEPOMEX id_asenta_cpcons)';
COMMENT ON COLUMN settlements.postal_code IS 'Postal code of the settlement';
COMMENT ON COLUMN settlements.name IS 'Settlement name';
COMMENT ON COLUMN settlements.slug IS 'Settlement slug';
COMMENT ON COLUMN settlements.settlement_type_id IS 'Type of settlement';
COMMENT ON COLUMN settlements.municipality_id IS 'Reference to municipality';
COMMENT ON COLUMN settlements.city_id IS 'Reference to city';
COMMENT ON COLUMN settlements.state_id IS 'Reference to state';
COMMENT ON COLUMN settlements.office_postal_code IS 'Postal code of post office that serves this settlement';
COMMENT ON COLUMN settlements.zone IS 'Urban or Rural zone';
COMMENT ON COLUMN settlements.latitude IS 'Latitude for geocoding';
COMMENT ON COLUMN settlements.longitude IS 'Longitude for geocoding';

CREATE INDEX idx_settlements_postal_code ON settlements(postal_code);
CREATE INDEX idx_settlements_name ON settlements(name);
CREATE INDEX idx_settlements_slug ON settlements(slug);
CREATE INDEX idx_settlements_municipality_id ON settlements(municipality_id);
CREATE INDEX idx_settlements_city_id ON settlements(city_id);
CREATE INDEX idx_settlements_state_municipality_postal_code ON settlements(state_id, municipality_id, postal_code);