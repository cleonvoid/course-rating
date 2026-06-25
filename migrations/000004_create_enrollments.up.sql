CREATE TABLE IF NOT EXISTS enrollments (
    id          SERIAL PRIMARY KEY,
    user_id     INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id   INTEGER NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    enrolled_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, course_id)
);

-- alice enrolled in: Go, PostgreSQL, Docker
-- bob enrolled in: Go, HTMX
-- carol enrolled in: PostgreSQL, Docker, System Design
-- dave enrolled in: HTMX, System Design
-- eve enrolled in all courses
INSERT INTO enrollments (user_id, course_id) VALUES
    (1, 1), (1, 2), (1, 3),
    (2, 1), (2, 4),
    (3, 2), (3, 3), (3, 5),
    (4, 4), (4, 5),
    (5, 1), (5, 2), (5, 3), (5, 4), (5, 5);
