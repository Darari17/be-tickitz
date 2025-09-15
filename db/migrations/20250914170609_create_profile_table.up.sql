CREATE TABLE IF NOT EXISTS profile (
    user_id uuid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    firstname varchar(100),
    lastname varchar(100),
    phone_number varchar(20),
    avatar text,
    point int DEFAULT 0,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp
);
