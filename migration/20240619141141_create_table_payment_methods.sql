-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    IF NOT EXISTS payment_methods (id SERIAL PRIMARY KEY, name TEXT NOT NULL);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE payment_methods;

-- +goose StatementEnd