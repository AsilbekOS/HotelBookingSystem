CREATE TABLE IF NOT EXISTS hotels (
    hotel_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    location VARCHAR(255) NOT NULL,
    rating DOUBLE PRECISION,
    address TEXT
);

CREATE TABLE IF NOT EXISTS rooms (
    room_id SERIAL PRIMARY KEY,
    hotel_id INT REFERENCES hotels(hotel_id),
    room_type VARCHAR(50),
    price_per_night DECIMAL(10, 2),
    availability BOOLEAN
);
