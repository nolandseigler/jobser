INSERT into auth.user_account (
    username,
    password
)
VALUES 
('admin', crypt('admin', gen_salt('bf', 8)));
