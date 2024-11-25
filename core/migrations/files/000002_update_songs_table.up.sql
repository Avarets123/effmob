BEGIN;
CREATE INDEX songs_name_idx ON songs (song);
CREATE INDEX songs_group_idx ON songs ("group");
COMMIT;