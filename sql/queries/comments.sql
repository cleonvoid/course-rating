-- name: ListCommentsByLesson :many
SELECT c.id, c.lesson_id, c.user_id, c.body, c.created_at, u.name AS user_name
FROM comments c
JOIN users u ON u.id = c.user_id
WHERE c.lesson_id = $1
ORDER BY c.created_at ASC;

-- name: CreateComment :exec
INSERT INTO comments (lesson_id, user_id, body)
VALUES ($1, $2, $3);

-- name: DeleteComment :exec
DELETE FROM comments WHERE id = $1 AND user_id = $2;
