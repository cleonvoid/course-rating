-- name: ListLessonsByCourse :many
SELECT id, course_id, title, body, "order", created_at
FROM lessons
WHERE course_id = $1
ORDER BY "order", id;

-- name: GetLesson :one
SELECT id, course_id, title, body, "order", created_at
FROM lessons
WHERE id = $1;
