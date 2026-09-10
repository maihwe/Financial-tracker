ALTER TABLE users
DROP CONSTRAINT users_role_check;

ALTER TABLE users
ADD CONSTRAINT users_role_check
CHECK (role IN ('user', 'admin', 'super_admin'));

CREATE UNIQUE INDEX one_super_admin_only
ON users (role)
WHERE role = 'super_admin';
