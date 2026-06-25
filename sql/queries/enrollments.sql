-- name: IsEnrolled :one
SELECT EXISTS (
    SELECT 1 FROM enrollments WHERE user_id = $1 AND course_id = $2
) AS enrolled;

-- name: EnrollUser :exec
INSERT INTO enrollments (user_id, course_id)
VALUES ($1, $2)
ON CONFLICT (user_id, course_id) DO NOTHING;
