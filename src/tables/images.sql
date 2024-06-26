CREATE TABLE IF NOT EXISTS images (
    id uuid PRIMARY KEY,
    name varchar(50) NOT NULL,
    owner_id integer REFERENCES users (id) NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);