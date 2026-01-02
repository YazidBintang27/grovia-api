-- +migrate Up
CREATE TABLE predicts(
    id SERIAL PRIMARY KEY,
    toddler_id INT NOT NULL,
    location_id INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    height DECIMAL(4,1) NOT NULL,
    age INT NOT NULL,
    sex VARCHAR(20) NOT NULL CHECK (sex IN ('M', 'F')),
    zscore DECIMAL(4,1) NOT NULL,
    nutritional_status VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_predicts_toddler FOREIGN KEY (toddler_id) REFERENCES toddlers(id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_predicts_location FOREIGN KEY (location_id) REFERENCES locations(id) ON UPDATE CASCADE ON DELETE CASCADE
);