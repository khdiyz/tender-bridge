-- +goose Up
CREATE TYPE tender_status AS ENUM (
    'open',
    'closed',
    'awarded'
);

CREATE TABLE IF NOT EXISTS "tenders"(
    "id" UUID PRIMARY KEY,
    "client_id" UUID NOT NULL,
    "title" VARCHAR(255) NOT NULL,
    "description" TEXT,
    "deadline" TIMESTAMP NOT NULL,
    "budget" BIGINT NOT NULL,
    "file" VARCHAR(64),
    "status" tender_status NOT NULL,
    FOREIGN KEY (client_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS "tenders";

DROP TYPE "tender_status";