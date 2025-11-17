-- Seed Admin User
-- Default admin user for initial login
-- Email: admin@pln.co.id
-- Password: Admin@12345
-- IMPORTANT: Change this password after first login!

INSERT INTO users (id, username, email, password_hash, role_id)
SELECT UUID(), 'admin', 'admin@pln.co.id',
'$2a$10$PW50IILE7Of1Uxsmx..qgu9eM3r8kEkCxKc9nB29gnw7rh9GENAK6',
id FROM roles WHERE name = 'SYSADMIN' LIMIT 1;
