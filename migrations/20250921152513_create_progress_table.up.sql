CREATE TABLE progress
(
    user_id  INT PRIMARY KEY,
    points   INT NOT NULL DEFAULT 0,
    progress INT NOT NULL DEFAULT 0
);