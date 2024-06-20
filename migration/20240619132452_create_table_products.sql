-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    IF NOT EXISTS products (
        id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid (),
        name TEXT NOT NULL,
        price DECIMAL NOT NULL,
        stock INTEGER NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

-- add trigger update_at
CREATE TRIGGER update_products_updated_at BEFORE
UPDATE ON products FOR EACH ROW EXECUTE PROCEDURE trigger_set_updated ();

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE products;

-- +goose StatementEnd