CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT current_timestamp
    );

CREATE TABLE IF NOT EXISTS segments (
    id SERIAL PRIMARY KEY,
    segment_name TEXT NOT NULL UNIQUE ,
    created_at timestamptz DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS segment_user (
    id_user BIGINT NOT NULL,
    id_segment BIGINT NOT NULL,
    delete_time timestamptz DEFAULT NULL,
    add_at timestamptz DEFAULT current_timestamp,
    PRIMARY KEY (id_user,id_segment)
);

-- INSERt into users values (1, current_timestamp);
-- insert into segments values (2, 'test2', current_timestamp);
-- insert into segment_user values (1, 1, default,default);
-- select * from segment_user;
-- select * from segments;


