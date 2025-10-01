CREATE TABLE IF NOT EXISTS prizes (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    descr TEXT NOT NULL, --зарезервированное слово description
    image_url TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);