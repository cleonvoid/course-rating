CREATE TABLE IF NOT EXISTS users (
    id            SERIAL PRIMARY KEY,
    name          TEXT NOT NULL,
    email         TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO users (name, email, password_hash) VALUES
    ('Alice Smith',   'alice@example.com', '$2a$10$Qq8JUCnWqtDwCRQd6UTl5OWRYMSCa6y07CPN3WCpOdw4XTdLaN3wK'),
    ('Bob Johnson',   'bob@example.com',   '$2a$10$0ceX0Klg2a89etIcrqqfYuYMNiesXuvETo.uyzColJJWcPXCgYama'),
    ('Carol Williams','carol@example.com', '$2a$10$Jt/EZBa3/SUZeu/mA2YnIOmuREG1odv0LxQdn.0WWRXJlbzRuTKCa'),
    ('Dave Brown',    'dave@example.com',  '$2a$10$Y4lVZwaOkWp0AHmaBg604.5MygclRCZyAPencGmdmozxPV6aW3icS'),
    ('Eve Davis',     'eve@example.com',   '$2a$10$8qg16aTgjEbL65qQkkC9n.AAhBOWj59IwD3eKV.1JVtffx.g1X.8W');

-- passwords: alice123, bob123, carol123, dave123, eve123
