CREATE TABLE IF NOT EXISTS lessons (
    id         SERIAL PRIMARY KEY,
    course_id  INTEGER NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    title      TEXT NOT NULL,
    body       TEXT NOT NULL DEFAULT '',
    "order"    INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO lessons (course_id, title, body, "order") VALUES
    (1, 'Getting Started with Go', 'Install Go and write your first Hello World program. We cover GOPATH, modules, and the basic project structure every Go developer needs to understand.', 1),
    (1, 'Types and Variables', 'Explore Go''s type system: basic types, composite types, and short variable declaration syntax. Learn why Go''s type inference keeps code concise without sacrificing safety.', 2),
    (1, 'Functions and Error Handling', 'Go''s approach to functions, multiple return values, and idiomatic error handling. Understand why explicit errors beat exceptions in most production scenarios.', 3),
    (2, 'Understanding the Query Planner', 'How PostgreSQL parses, plans, and executes queries. Learn to read EXPLAIN and EXPLAIN ANALYZE output so you can diagnose slow queries with confidence.', 1),
    (2, 'Index Strategies', 'B-tree, hash, GIN, and BRIN indexes compared. When to use each, how to measure index impact, and common indexing mistakes that hurt write performance.', 2),
    (3, 'Docker Fundamentals', 'Images, containers, and the Docker daemon explained. Write your first Dockerfile and understand the layer cache that makes builds fast.', 1),
    (3, 'Kubernetes Architecture', 'Nodes, pods, deployments, and services from the ground up. See how the control plane and kubelet cooperate to keep your workloads running.', 2),
    (4, 'Hypermedia as the Engine', 'The HTMX philosophy: why HTML over the wire beats JSON APIs for many use cases. Understand HATEOAS and how it simplifies frontend architecture.', 1),
    (4, 'hx-get, hx-post, and Targets', 'Core HTMX attributes and how they replace full-page reloads with surgical DOM swaps. Walk through real examples that replace hundreds of lines of JavaScript.', 2),
    (5, 'Scalability Basics', 'Horizontal vs vertical scaling trade-offs. How load balancers, caches, and read replicas each address different bottlenecks at different traffic levels.', 1),
    (5, 'Designing a URL Shortener', 'Step-by-step system design walkthrough: capacity estimation, choosing a data store, handling redirects at scale, and avoiding the pitfalls most candidates miss.', 2);
