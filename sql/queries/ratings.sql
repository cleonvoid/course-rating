-- name: UpsertRating :one
INSERT INTO ratings (course_id, user_id, stars)
VALUES ($1, $2, $3)
ON CONFLICT (course_id, user_id)
DO UPDATE SET stars = EXCLUDED.stars
RETURNING id, course_id, user_id, stars, created_at;

-- name: GetRatingByUserAndCourse :one
SELECT id, course_id, user_id, stars, created_at
FROM ratings
WHERE course_id = $1 AND user_id = $2;
