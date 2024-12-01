CREATE TABLE IF NOT EXISTS songs(
    song_id bigserial PRIMARY KEY,
    group text NOT NULL,
    song text NOT NULL,
    release_date text NOT NULL,
    song_text text[],
    link text
);