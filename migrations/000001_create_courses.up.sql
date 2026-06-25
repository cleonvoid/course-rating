CREATE TABLE IF NOT EXISTS courses (
    id          SERIAL PRIMARY KEY,
    title       TEXT NOT NULL,
    description TEXT NOT NULL,
    instructor  TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO courses (title, description, instructor) VALUES
    ('Introduction to Go', 'Learn the fundamentals of Go programming language.', 'Rob Pike'),
    ('Advanced PostgreSQL', 'Deep dive into PostgreSQL performance and internals.', 'Bruce Momjian'),
    ('Docker & Kubernetes', 'Containerize and orchestrate modern applications.', 'Kelsey Hightower'),
    ('Web Development with HTMX', 'Build hypermedia-driven web apps with minimal JavaScript.', 'Carson Gross'),
    ('System Design Fundamentals', 'Learn to design scalable distributed systems.', 'Alex Xu');
