-- +goose Up
CREATE TABLE urls {
    id VARCHAR(10) PRIMARY KEY,
    url TEXT NOT NULL,
    alias TEXT NOT NULL UNIQUE
};


-- +goose Down
DROP TABLE urls;
