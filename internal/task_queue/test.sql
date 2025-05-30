DROP TABLE IF EXISTS tasks;
CREATE TABLE tasks(
    uuid TEXT PRIMARY KEY,
    type TEXT NOT NULL,
    handler TEXT NOT NULL,
    payload TEXT NOT NULL,
    status TEXT NOT NULL,
    retries INTEGER DEFAULT 0,
    max_retries INTEGET DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);