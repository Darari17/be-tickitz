CREATE TABLE IF NOT EXISTS movies (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    backdrop_path TEXT NOT NULL,
    overview TEXT NOT NULL,
    popularity DOUBLE PRECISION NOT NULL,
    poster_path TEXT NOT NULL,
    release_date DATE NOT NULL,
    duration INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    director_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP
);
