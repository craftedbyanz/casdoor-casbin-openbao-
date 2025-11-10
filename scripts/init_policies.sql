-- Sample policies for testing
-- Run this after starting the server to add initial policies

-- Policies (p_type, v0=subject, v1=object, v2=action)
INSERT INTO casbin_rule (p_type, v0, v1, v2) VALUES 
('p', 'admin', '/api/users', 'read'),
('p', 'admin', '/api/admin/*', 'write'),
('p', 'user', '/api/users/profile', 'read'),
('p', 'user', '/api/protected', 'read');

-- Role assignments (g_type, v0=user, v1=role)
INSERT INTO casbin_rule (p_type, v0, v1) VALUES 
('g', 'admin', 'admin'),
('g', 'testuser', 'user');