CREATE TABLE groups (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    group_name TEXT NOT NULL UNIQUE
);

CREATE TABLE songs (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    group_id INTEGER,
    song_name TEXT NOT NULL UNIQUE,
    release_date TEXT DEFAULT NOW()::date,
    text TEXT,
    link TEXT,
    FOREIGN KEY (group_id) REFERENCES groups(id)
);

INSERT INTO groups(group_name) VALUES 
('Imagine Dragons'), 
('Linkin Park');

INSERT INTO songs(group_id, song_name, text, link) VALUES 
(2, 'Believer', 'I am believer', 'https://youtu.be/7wtfhZwyrcc?si=AhOKtDFQw19Cmfmy'),
(2, 'Thunder', 'Before the thunder', 'https://youtu.be/fKopy74weus?si=PPbCVQS28w3Fp0Ga'),
(1, 'In the end', 'It doens''t even matter', 'https://youtu.be/eVTXPUF4Oz4?si=XRrAbzJJqOO4jAJx');
