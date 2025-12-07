-- Rollback initial migration

DROP INDEX IF EXISTS idx_user_stats_user_id;
DROP TABLE IF EXISTS user_stats;
