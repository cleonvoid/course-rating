DROP TABLE IF EXISTS ratings;

CREATE TABLE ratings (
    id         SERIAL PRIMARY KEY,
    course_id  INTEGER  NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    student_id TEXT     NOT NULL,
    stars      SMALLINT NOT NULL CHECK (stars >= 1 AND stars <= 5),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (course_id, student_id)
);
