-- Add checkin_id column to wallet table for syncing with checkin records
ALTER TABLE wallet ADD COLUMN checkin_id BIGINT UNSIGNED DEFAULT 0 AFTER user_id;

-- Add cancel checkin API support
-- When deleting a checkin-related wallet record, also delete the checkin record