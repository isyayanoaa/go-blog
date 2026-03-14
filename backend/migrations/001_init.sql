-- Users
CREATE TABLE users (
    id         SERIAL PRIMARY KEY,
    github_id  BIGINT UNIQUE NOT NULL,
    username   VARCHAR(100) NOT NULL,
    avatar     TEXT,
    is_admin   BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Posts
CREATE TABLE posts (
    id         SERIAL PRIMARY KEY,
    title      VARCHAR(200) NOT NULL,
    slug       VARCHAR(200) UNIQUE NOT NULL,
    content    TEXT NOT NULL,
    summary    TEXT,
    cover      TEXT,
    category   VARCHAR(100),
    tags       TEXT[],
    status     VARCHAR(20) DEFAULT 'draft',  -- published | draft
    like_count INT DEFAULT 0,
    view_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Moments
CREATE TABLE moments (
    id         SERIAL PRIMARY KEY,
    content    TEXT,
    images     TEXT[],
    video      TEXT,
    like_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Comments
CREATE TABLE comments (
    id          SERIAL PRIMARY KEY,
    user_id     INT REFERENCES users(id) ON DELETE CASCADE,
    target_type VARCHAR(20) NOT NULL,  -- post | moment
    target_id   INT NOT NULL,
    parent_id   INT REFERENCES comments(id) ON DELETE CASCADE,
    content     TEXT NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW()
);

-- Likes
CREATE TABLE likes (
    id          SERIAL PRIMARY KEY,
    user_id     INT REFERENCES users(id) ON DELETE CASCADE,
    target_type VARCHAR(20) NOT NULL,  -- post | moment
    target_id   INT NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW(),
    UNIQUE (user_id, target_type, target_id)
);

-- Indexes
CREATE INDEX idx_posts_status ON posts(status);
CREATE INDEX idx_posts_category ON posts(category);
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);
CREATE INDEX idx_moments_created_at ON moments(created_at DESC);
CREATE INDEX idx_comments_target ON comments(target_type, target_id);
CREATE INDEX idx_likes_target ON likes(target_type, target_id);
