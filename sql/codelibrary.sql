SET client_min_messages TO WARNING;

CREATE TABLE IF NOT EXISTS "user" (
    id uuid PRIMARY KEY NOT NULL,
    username varchar(255) NOT NULL,
    password_hash varchar(255) NOT NULL,
    CONSTRAINT username_unique UNIQUE (username)
);

CREATE TABLE IF NOT EXISTS language (
    id varchar(255) PRIMARY KEY NOT NULL,
    name text NOT NULL
);

CREATE TABLE IF NOT EXISTS codesample (
    id uuid PRIMARY KEY NOT NULL,
    submitted_by_id uuid NOT NULL
        REFERENCES "user"(id)
        ON DELETE RESTRICT,
    language_id varchar(255) NOT NULL
        REFERENCES language(id)
        ON DELETE RESTRICT,
    title text NOT NULL,
    description text NOT NULL,
    body text NOT NULL,
    search_index tsvector NOT NULL,
    created timestamp with time zone NOT NULL,
    modified timestamp with time zone NOT NULL
);

CREATE INDEX IF NOT EXISTS codesample_fulltext_index
    ON codesample USING GIN (search_index);

INSERT INTO language (id, name) VALUES
    ('ada', 'Ada'),
    ('bash', 'Bash'),
    ('c', 'C'),
    ('cobol', 'COBOL'),
    ('cpp', 'C++'),
    ('csharp', 'C#'),
    ('d', 'D'),
    ('dart', 'Dart'),
    ('fortran', 'Fortan'),
    ('go', 'Go'),
    ('java', 'Java'),
    ('javascript', 'JavaScript'),
    ('julia', 'Julia'),
    ('kotlin', 'Kotlin'),
    ('lisp', 'Lisp'),
    ('logo', 'Logo'),
    ('lua', 'Lua'),
    ('matlab', 'Matlab'),
    ('mysql', 'MySQL'),
    ('objc', 'Objective-C'),
    ('perl', 'perl'),
    ('php', 'PHP'),
    ('postgres', 'Postgres'),
    ('powershell', 'PowerShell'),
    ('prolog', 'Prolog'),
    ('python', 'Python'),
    ('r', 'R'),
    ('ruby', 'Ruby'),
    ('rust', 'Rust'),
    ('scala', 'Scala'),
    ('scheme', 'Scheme'),
    ('scratch', 'Scratch'),
    ('sql', 'SQL'),
    ('swift', 'Swift'),
    ('typescript', 'TypeScript'),
    ('visualbasic', 'Visual Basic'),
    ('zsh', 'zsh')
ON CONFLICT DO NOTHING;
