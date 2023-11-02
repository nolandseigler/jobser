INSERT into user_account (
    username,
    password
)
VALUES 
('admin', crypt('admin', gen_salt('bf', 8)));
