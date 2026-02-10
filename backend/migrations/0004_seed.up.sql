INSERT INTO users (name, email, password_hash, role)
VALUES
  ('Admin', 'admin@example.com', '$2a$10$1MtIJ0kaBf4bnVzx.NsRIe1y.jAl2ud9Iz6BOoWTxWQ582Nq9kA1G', 'admin'),
  ('User',  'user@example.com',  '$2a$10$1MtIJ0kaBf4bnVzx.NsRIe1y.jAl2ud9Iz6BOoWTxWQ582Nq9kA1G', 'user')
ON CONFLICT (email) DO NOTHING;

WITH c1 AS (
    INSERT INTO categories (name) VALUES ('Electronics')
    ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
    RETURNING id
),
c2 AS (
    INSERT INTO categories (name) VALUES ('Books')
    ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
    RETURNING id
),
c3 AS (
    INSERT INTO categories (name) VALUES ('Home')
    ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
    RETURNING id
)
INSERT INTO products (category_id, name, description, price)
VALUES
 ((SELECT id FROM c1), 'Mechanical Keyboard', 'Keyboard for typing.', 899000.00),
 ((SELECT id FROM c1), 'Wireless Mouse', 'Simple wireless mouse.', 199000.00),
 ((SELECT id FROM c2), 'Go Programming Book', 'Learn Go from basics.', 250000.00),
 ((SELECT id FROM c3), 'Desk Lamp', 'Warm light desk lamp.', 150000.00);