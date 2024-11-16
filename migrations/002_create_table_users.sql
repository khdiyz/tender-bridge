-- +goose Up
CREATE TABLE IF NOT EXISTS "users"(
    "id" UUID PRIMARY KEY,
    "role" VARCHAR(64) NOT NULL,
    "username" VARCHAR(64) NOT NULL,
    "email" VARCHAR(64) NOT NULL,
    "password" TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS "users";