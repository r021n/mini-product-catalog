DELETE FROM products;
DELETE FROM categories;
DELETE FROM users WHERE email IN ('admin@example.com', 'user@example.com');