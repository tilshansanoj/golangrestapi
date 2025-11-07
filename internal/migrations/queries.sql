-- name: CreateUser :one
INSERT INTO users (username, email, password, created, updated)
VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
    RETURNING id, username, email, created, updated;

-- name: GetUserByID :one
SELECT id, username, email, created, updated
FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT id, username, email, created, updated
FROM users
ORDER BY id;

-- name: GetUserByUsername :one
SELECT id, username, email, password, created, updated
FROM users
WHERE username = $1;

-- name: UpdateUser :one
UPDATE users
    SET username = $2,
    email = $3
WHERE id = $1
RETURNING *;

-- name: DeleteUser :one
DELETE FROM users
WHERE id = $1
RETURNING *;

-- name: CreateBlog :one
INSERT INTO blogs (title, content, user_id, created, updated)
VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
    RETURNING id, title, content, user_id, created, updated;

-- name: ListBlogs :many
SELECT id, title, content, user_id, created, updated
FROM blogs
ORDER BY id;

-- name: GetBlogbyId :one
SELECT id, title, content, user_id, created, updated
FROM blogs
WHERE id = $1;

-- name: UpdateBlog :one
UPDATE blogs
    SET title = $2,
    content = $3
WHERE id = $1
RETURNING *;

-- name: DeleteBlog :one
DELETE FROM blogs
WHERE id = $1
RETURNING *;

-- name: CreateChild :one
INSERT INTO children (username, email, password, parent_id, created, updated)
VALUES ($1, $2, $3, $4,CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
    RETURNING id, username, email, parent_id, created, updated;

-- name: GetChildByID :one
SELECT id, username, email, parent_id, created, updated
FROM children
WHERE id = $1;

-- name: ListChildren :many
SELECT id, username, email, parent_id, created, updated
FROM children
ORDER BY id;