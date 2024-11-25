BEGIN;

CREATE TABLE couplets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    song_id UUID NOT NULL,
    couplet VARCHAR(512) NOT NULL,
    couplet_num INTEGER NOT NULL,
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_couplets FOREIGN KEY (song_id) REFERENCES songs(id)
);

COMMIT;