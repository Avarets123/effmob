BEGIN;

    INSERT INTO songs(id, "group", link, song, release_date) VALUES
    ('984f5c01-6f31-445d-bc72-d6ef2503f321', 'best', 'first', 'https://www.youtube.com/watch?v=Xsp3_a-PMTw', '2024-11-24T22:55:01.385Z'),
    ('3001f2b5-5bbf-49da-a2b8-cc23671a0c1d', 'best', 'second', 'https://www.youtube.com/watch?v=Xsp3_a-PMTw', '2024-10-24T22:55:01.385Z'),
    ('5f7fa130-6ffb-462f-a481-061430d8a619', 'second best', 'last', 'https://www.youtube.com/watch?v=Xsp3_a-PMTw', '2024-01-24T22:55:01.385Z');


    INSERT INTO couplets(song_id, couplet, couplet_num) VALUES
    ('984f5c01-6f31-445d-bc72-d6ef2503f321', 'Ooh baby, dont you know I suffer?', 1),
    ('984f5c01-6f31-445d-bc72-d6ef2503f321', 'Ooh baby, can you hear me moan?', 2),
    ('984f5c01-6f31-445d-bc72-d6ef2503f321', 'You caught me under false pretenses', 3),
    ('984f5c01-6f31-445d-bc72-d6ef2503f321', 'How long before you let me go?', 4),
    ('984f5c01-6f31-445d-bc72-d6ef2503f321', 'You set my soul alight', 5),
    ('984f5c01-6f31-445d-bc72-d6ef2503f321', 'You set my soul alight', 6),
    ('3001f2b5-5bbf-49da-a2b8-cc23671a0c1d', 'one', 1),
    ('3001f2b5-5bbf-49da-a2b8-cc23671a0c1d', 'two', 2),
    ('3001f2b5-5bbf-49da-a2b8-cc23671a0c1d', 'three', 3);

COMMIT;