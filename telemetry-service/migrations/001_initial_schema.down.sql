-- Drop trigger first
DROP TRIGGER IF EXISTS update_devices_updated_at ON devices;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_threat_findings_rule_id;
DROP INDEX IF EXISTS idx_threat_findings_timestamp;
DROP INDEX IF EXISTS idx_threat_findings_severity;
DROP INDEX IF EXISTS idx_threat_findings_device_id;
DROP INDEX IF EXISTS idx_browser_sessions_collected_at;
DROP INDEX IF EXISTS idx_browser_sessions_device_id;
DROP INDEX IF EXISTS idx_containers_container_id;
DROP INDEX IF EXISTS idx_containers_collected_at;
DROP INDEX IF EXISTS idx_containers_device_id;
DROP INDEX IF EXISTS idx_processes_pid_device_id;
DROP INDEX IF EXISTS idx_processes_collected_at;
DROP INDEX IF EXISTS idx_processes_device_id;
DROP INDEX IF EXISTS idx_devices_last_seen_at;
DROP INDEX IF EXISTS idx_devices_mac_address;

-- Drop tables (in reverse order due to foreign key constraints)
DROP TABLE IF EXISTS threat_findings;
DROP TABLE IF EXISTS browser_sessions;
DROP TABLE IF EXISTS containers;
DROP TABLE IF EXISTS processes;
DROP TABLE IF EXISTS devices;
