CREATE TABLE reservation(
    id SERIAL PRIMARY KEY,
    provider_id  BIGINT,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_phone TEXT NOT NULL,
    room_id BIGINT NOT NULL,
    hotel_id BIGINT,
    date_from DATE NOT NULL,
    date_to DATE NOT NULL,
    reservation_status TEXT NOT NULL
);