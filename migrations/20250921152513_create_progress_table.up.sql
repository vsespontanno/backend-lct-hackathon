CREATE TABLE IF NOT EXISTS progress
(
    user_id  INT PRIMARY KEY,
    points   INT NOT NULL DEFAULT 0,
    progress INT NOT NULL DEFAULT 0
);