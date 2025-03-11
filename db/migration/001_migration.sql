-- +goose Up
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    profile_name VARCHAR(100),
    email VARCHAR(100) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    bio TEXT DEFAULT '',
    phone_number VARCHAR(15),
    profile_pic TEXT,
    online_status VARCHAR(20) DEFAULT 'offline',
    date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE posts (
    post_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
    content TEXT,
    -- image_url TEXT,
    date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    view_count INTEGER DEFAULT 0, -- Tracks unique views
    repost_count INTEGER DEFAULT 0,
    comment_count INTEGER DEFAULT 0, -- Tracks the number of reposts
    like_count INTEGER DEFAULT 0, -- Tracks total likes
    bookmark_count INTEGER DEFAULT 0 -- Tracks total bookmarks
);

CREATE TABLE comments (
    comment_id SERIAL PRIMARY KEY,
    post_id INTEGER REFERENCES posts(post_id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
    parent_id INTEGER REFERENCES comments(comment_id) ON DELETE CASCADE,
    content TEXT,
    date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE likes (
    like_id SERIAL PRIMARY KEY,
    post_id INTEGER REFERENCES posts(post_id) ON DELETE CASCADE,
    comment_id INTEGER REFERENCES comments(comment_id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
    date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (post_id, user_id) -- Ensures a user can like a post only once
);

CREATE TABLE bookmarks (
    bookmark_id SERIAL PRIMARY KEY,
    post_id INTEGER REFERENCES posts(post_id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
    date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (post_id, user_id) -- Ensures a user can bookmark a post only once
);

-- Table to track unique post views
CREATE TABLE post_views (
    view_id SERIAL PRIMARY KEY,
    post_id INTEGER REFERENCES posts(post_id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE, -- Nullable for guest views
    ip_address VARCHAR(45), -- Track guest views via IP
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (post_id, user_id, ip_address) -- Ensures unique views
);

-- Table to track reposts
CREATE TABLE reposts (
    repost_id SERIAL PRIMARY KEY,
    original_post_id INTEGER REFERENCES posts(post_id) ON DELETE CASCADE, -- Original post
    user_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE, -- Who reposted
    reposted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (original_post_id, user_id) -- Ensures a user can only repost once
);

-- +goose Down
DROP TABLE IF EXISTS reposts;
DROP TABLE IF EXISTS post_views;
DROP TABLE IF EXISTS bookmarks;
DROP TABLE IF EXISTS likes;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS users;

