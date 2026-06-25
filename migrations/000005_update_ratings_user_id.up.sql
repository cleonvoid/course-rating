-- Replace text student_id with integer user_id FK to users
DROP TABLE IF EXISTS ratings;

CREATE TABLE ratings (
    id         SERIAL PRIMARY KEY,
    course_id  INTEGER  NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    user_id    INTEGER  NOT NULL REFERENCES users(id)   ON DELETE CASCADE,
    stars      SMALLINT NOT NULL CHECK (stars >= 1 AND stars <= 5),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (course_id, user_id)
);
