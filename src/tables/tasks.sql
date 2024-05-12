CREATE TABLE IF NOT EXISTS tasks (
    id serial PRIMARY KEY,
    name varchar(50) NOT NULL,
    description varchar(1024),
    data bytea NOT NULL,
    owner_id integer REFERENCES users (id) NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);