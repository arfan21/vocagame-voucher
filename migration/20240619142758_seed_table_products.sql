-- +goose Up
-- +goose StatementBegin
INSERT INTO
    products (name, price, stock)
VALUES
    ('Genshin Impact', 1000, 20),
    ('Valorant', 500, 15),
    ('PUBG', 300, 10),
    ('Mobile Legend', 200, 5),
    ('Free Fire', 100, 1);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DELETE FROM products
WHERE
    name IN (
        'Genshin Impact',
        'Valorant',
        'PUBG',
        'Mobile Legend',
        'Free Fire'
    );

-- +goose StatementEnd