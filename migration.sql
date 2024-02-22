DROP TABLE IF EXISTS conn_log;
CREATE TABLE IF NOT EXISTS conn_log (
    user_id BIGINT,
    ip_addr VARCHAR(15),
    ts TIMESTAMP
);

CREATE INDEX idx_user_id ON conn_log (user_id);

INSERT INTO conn_log (user_id, ip_addr, ts)
SELECT
    (random() * 4 + 1)::int AS user_id,
    ('{"127.0.0.1", "127.0.0.2", "127.0.0.3"}'::text[])[floor(random()*3)+1] AS ip_addr,
    timestamp '2021-01-01 00:00:00' + (random() * (365 * 24 * 60 * 60)) * interval '1 second' AS ts
FROM
    generate_series(1, 1000000) s(i);