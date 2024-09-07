CREATE TABLE IF NOT EXISTS bookings (
    booking_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    hotel_id INT NOT NULL,
    room_id INT NOT NULL,
    room_type VARCHAR(100) NOT NULL,
    check_in_date DATE NOT NULL,        
    check_out_date DATE NOT NULL,        
    total_amount NUMERIC(10, 2) NOT NULL,
    status VARCHAR(50) NOT NULL
);
