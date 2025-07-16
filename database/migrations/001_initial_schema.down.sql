-- Rollback Migration 001: Initial Schema
-- This will drop all tables created in the initial schema

DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS threats;
DROP TABLE IF EXISTS containers;
DROP TABLE IF EXISTS processes;
DROP TABLE IF EXISTS device_metrics;
DROP TABLE IF EXISTS devices;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS organizations;

-- Drop the function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop UUID extension
DROP EXTENSION IF EXISTS "uuid-ossp";
