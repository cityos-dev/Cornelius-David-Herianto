-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS files(
    id          VARCHAR,
    size        INTEGER,
    path        VARCHAR,
    created_at  TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT files_pk PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS files;
-- +goose StatementEnd
