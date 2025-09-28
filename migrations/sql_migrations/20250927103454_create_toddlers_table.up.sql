-- +migrate Up
CREATE TABLE toddlers(
    id SERIAL PRIMARY KEY,
    parent_id INT NOT NULL,
    location_id INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    birthdate DATE NOT NULL,
    sex VARCHAR(20) NOT NULL CHECK (sex IN ('M', 'F')),
    height DECIMAL(4,1) NOT NULL,
    profile_picture VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_toddlers_parent FOREIGN KEY (parent_id) REFERENCES parents(id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_toddlers_location FOREIGN KEY (location_id) REFERENCES locations(id) ON UPDATE CASCADE ON DELETE CASCADE
);