-- +goose Up
CREATE TABLE users (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name text UNIQUE NOT NULL
);
  
-- +goose Down
DROP TABLE users;
