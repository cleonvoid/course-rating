-- name: ListCourses :many
SELECT
    c.id,
    c.title,
    c.description,
    c.instructor,
    c.created_at,
    COALESCE(ROUND(AVG(r.stars)::numeric, 1), 0)::float8 AS avg_rating,
    COUNT(r.id)::int AS rating_count
FROM courses c
LEFT JOIN ratings r ON r.course_id = c.id
GROUP BY c.id
ORDER BY c.id;

-- name: GetCourse :one
SELECT
    c.id,
    c.title,
    c.description,
    c.instructor,
    c.created_at,
    COALESCE(ROUND(AVG(r.stars)::numeric, 1), 0)::float8 AS avg_rating,
    COUNT(r.id)::int AS rating_count
FROM courses c
LEFT JOIN ratings r ON r.course_id = c.id
WHERE c.id = $1
GROUP BY c.id;
