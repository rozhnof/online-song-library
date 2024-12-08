CREATE TABLE groups (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_groups_deleted_at ON groups(deleted_at);

CREATE TABLE songs (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    group_id UUID REFERENCES groups(id) NOT NULL,
    release_date DATE NOT NULL,
    text TEXT NOT NULL,
    link TEXT NOT NULL,
    deleted_at TIMESTAMP,
    UNIQUE(name, group_id)
);

CREATE INDEX idx_songs_release_date ON songs(release_date);
CREATE INDEX idx_songs_text ON songs(text);
CREATE INDEX idx_songs_link ON songs(link);
CREATE INDEX idx_songs_deleted_at ON songs(deleted_at);