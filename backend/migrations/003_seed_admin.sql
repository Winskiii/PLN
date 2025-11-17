-- Seed Admin User
-- Default admin user for initial login
-- Email: admin@pln.co.id
-- Password: Admin@12345
-- IMPORTANT: Change this password after first login!

INSERT INTO users (id, username, email, password_hash, role_id)
SELECT UUID(), 'admin', 'admin@pln.co.id',
'$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYCn.OC4x/S',
id FROM roles WHERE name = 'SYSADMIN' LIMIT 1;
