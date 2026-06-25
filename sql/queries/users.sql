-- name: GetUserByEmail :one
SELECT id, name, email, password_hash, created_at
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT id, name, email
FROM users
WHERE id = $1;
