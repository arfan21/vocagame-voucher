-- +goose Up
-- +goose StatementBegin
INSERT INTO
    payment_methods (id, name)
VALUES
    (1, 'SNAP Midtrans');

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT
    'down SQL query';

-- +goose StatementEnd