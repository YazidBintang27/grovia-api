-- +migrate Up
CREATE TABLE parents(
    id SERIAL PRIMARY KEY,
    location_id INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    phone_number VARCHAR(100) NOT NULL UNIQUE,
    address VARCHAR(100) NOT NULL,
    nik VARCHAR(100) NOT NULL UNIQUE,
    job VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_parents_location FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE CASCADE ON UPDATE CASCADE,

    CONSTRAINT unique_parents_record UNIQUE (location_id, name, phone_number, nik)
);