-- +goose Up
CREATE TABLE posts (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  title text,
  url text UNIQUE NOT NULL,
  description text,
  published_at text,
  feed_id uuid NOT NULL REFERENCES feeds ON DELETE CASCADE,
  FOREIGN KEY (feed_id) REFERENCES feeds(id)
);

-- +goose Down
DROP TABLE posts;
