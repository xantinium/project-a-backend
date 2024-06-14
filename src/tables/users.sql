CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY,
    first_name varchar(50),
    last_name varchar(50),
    avatar_id uuid,
    oauth_service smallint NOT NULL,
    yandex_profile_id varchar(20),
    google_profile_id varchar(30),
    theme smallint NOT NULL,
    lang smallint NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);