-- Seed: Default admin user for development
-- Username: admin | Password: admin123
-- Password hash generated with bcrypt cost 10

INSERT INTO users (id, username, password_hash, created_at)
VALUES (
    '01JMQZX0000000000000000000',
    'admin',
    '$2a$10$bPPh3kDSLtx1nhRz4dNgTea0BChZ8x.O1/KNI3axovlHdnWO9UGCO',
    NOW()
) ON CONFLICT (username) DO NOTHING;
