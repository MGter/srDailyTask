-- migrations/003_add_user_avatar.sql
-- Add avatar URL for user profile images

ALTER TABLE users ADD COLUMN avatar_url VARCHAR(255);
