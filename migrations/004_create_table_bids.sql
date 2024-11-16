-- +goose Up

CREATE TABLE IF NOT EXISTS "bids"(
    "id" UUID PRIMARY KEY,
    "contractor_id" UUID NOT NULL,
    "tender_id" UUID NOT NULL,
    "price" BIGINT NOT NULL,
    "delivery_time" INTEGER NOT NULL,
    "comment" TEXT,
    "status" VARCHAR(64),
    FOREIGN KEY (contractor_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (tender_id) REFERENCES tenders(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS "bids";