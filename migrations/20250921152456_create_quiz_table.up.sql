CREATE TABLE IF NOT EXISTS quiz (
    id INT PRIMARY KEY ,
    title TEXT NOT NULL ,
    content TEXT NOT NULL ,
    options TEXT[] NOT NUll ,
    correct_answer TEXT NOT NULL 
);