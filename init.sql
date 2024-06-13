-- Создание таблицы posts
CREATE TABLE posts (
    id VARCHAR(100) PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    comments_disabled BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы comments
CREATE TABLE comments (
    id VARCHAR(100) PRIMARY KEY,
    post_id VARCHAR(100) REFERENCES posts(id) ON DELETE CASCADE,
    parent_id VARCHAR(100) REFERENCES comments(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы pairs для иерархии комментариев
CREATE TABLE pairs (
    parent_id VARCHAR(100) REFERENCES comments(id) ON DELETE CASCADE,
    child_id VARCHAR(100) REFERENCES comments(id) ON DELETE CASCADE,
    PRIMARY KEY (parent_id, child_id)
);
