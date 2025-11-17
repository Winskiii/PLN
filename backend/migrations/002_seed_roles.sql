-- Seed Roles
-- Insert default roles for the Work Management System

INSERT INTO roles (id, name, description) VALUES
(UUID(), 'SYSADMIN', 'System Administrator - Full access to all features and data'),
(UUID(), 'ADMIN', 'Administrator - Manage users and all projects'),
(UUID(), 'MANAGER', 'Manager - Manage assigned projects and create tasks'),
(UUID(), 'EMPLOYEE', 'Employee - View and update assigned tasks');
