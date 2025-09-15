CREATE TABLE IF NOT EXISTS users (
    id uuid NOT NULL,
    email varchar(255) NOT NULL UNIQUE,
    password varchar(255) NOT NULL,
    role varchar(50) NOT NULL CHECK (role IN ('admin', 'user', 'general')),
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);
