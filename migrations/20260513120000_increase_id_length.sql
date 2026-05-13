-- +goose Up
CREATE TABLE urls_new (
    id VARCHAR(20),
    url TEXT,
    alias TEXT NOT NULL UNIQUE
);

INSERT INTO urls_new (id, url, alias)
SELECT id, url, alias FROM urls;

DROP TABLE urls;

ALTER TABLE urls_new RENAME TO urls;

-- +goose Down
CREATE TABLE urls_new (
    id VARCHAR(10),
    url TEXT NOT NULL,
    alias TEXT NOT NULL UNIQUE
);

INSERT INTO urls_new (id, url, alias)
SELECT id, url, alias FROM urls;

DROP TABLE urls;

ALTER TABLE urls_new RENAME TO urls;