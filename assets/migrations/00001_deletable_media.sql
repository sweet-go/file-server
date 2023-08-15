-- +migrate Up notransaction

CREATE TYPE delete_rule AS ENUM ('MANUAL_DELETE');

CREATE TABLE IF NOT EXISTS "deletable_media" (
    "id" VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    delete_rule delete_rule DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL
);

-- +migrate Down

DROP TABLE IF EXISTS "deletable_media";