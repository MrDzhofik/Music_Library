CREATE TABLE groups (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    group_name TEXT NOT NULL UNIQUE,
);

CREATE TABLE songs (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    group_id INTEGER,
    song_name TEXT NOT NULL UNIQUE,
    release_date date,
    text TEXT,
    link TEXT,
    FOREIGN KEY (group_id) REFERENCES groups(id)
);
