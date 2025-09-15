CREATE TABLE IF NOT EXISTS schedules (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    movies_id INT NOT NULL REFERENCES movies(id),
    cinemas_id INT NOT NULL REFERENCES cinemas(id),
    times_id INT NOT NULL REFERENCES times(id),
    locations_id INT NOT NULL REFERENCES locations(id),
    date DATE NOT NULL
);
