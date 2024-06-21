-- +goose Up
-- +goose StatementBegin
CREATE TYPE payment_status AS ENUM ('waiting_payment', 'completed', 'failed');

CREATE TABLE
    IF NOT EXISTS transactions (
        id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid (),
        product_id UUID NOT NULL,
        payment_method_id INTEGER NOT NULL,
        email TEXT NOT NULL,
        quantity INTEGER NOT NULL,
        total_price DECIMAL NOT NULL,
        status payment_status NOT NULL DEFAULT 'waiting_payment',
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT fk_product_id FOREIGN KEY (product_id) REFERENCES products (id),
        CONSTRAINT fk_payment_method_id FOREIGN KEY (payment_method_id) REFERENCES payment_methods (id)
    );

-- add trigger update_at
CREATE TRIGGER update_transactions_updated_at BEFORE
UPDATE ON transactions FOR EACH ROW EXECUTE PROCEDURE trigger_set_updated ();

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE transactions;

DROP TYPE payment_status;

-- +goose StatementEnd