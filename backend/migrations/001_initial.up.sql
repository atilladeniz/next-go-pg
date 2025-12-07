-- Initial migration: Create core tables
-- This migration sets up the initial database schema

-- Enable UUID extension (PostgreSQL)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- user_stats table (matches domain.UserStats)
CREATE TABLE IF NOT EXISTS user_stats (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL UNIQUE,
    project_count INTEGER DEFAULT 0,
    activity_today INTEGER DEFAULT 0,
    notifications INTEGER DEFAULT 0,
    last_login TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    member_since TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index for faster user lookups
CREATE INDEX IF NOT EXISTS idx_user_stats_user_id ON user_stats(user_id);
