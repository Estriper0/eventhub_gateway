CREATE TABLE event (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    about TEXT,
    start_date TIMESTAMP NOT NULL,
    location VARCHAR(255),
    status VARCHAR(50) NOT NULL,
    max_attendees SMALLINT
);