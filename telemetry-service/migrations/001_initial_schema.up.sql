-- Create devices table
CREATE TABLE IF NOT EXISTS devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mac_address VARCHAR(17) NOT NULL UNIQUE,
    hostname VARCHAR(255) NOT NULL,
    os VARCHAR(100) NOT NULL,
    platform VARCHAR(100) NOT NULL,
    version VARCHAR(100) NOT NULL,
    current_username VARCHAR(100) NOT NULL,
    user_id UUID,
    org_unit VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_seen_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create processes table
CREATE TABLE IF NOT EXISTS processes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    pid INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    cmdline TEXT[],
    username VARCHAR(100),
    exe_path TEXT,
    start_time BIGINT,
    status VARCHAR(50),
    sha256 VARCHAR(64),
    version VARCHAR(100),
    file_size BIGINT,
    collected_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create containers table
CREATE TABLE IF NOT EXISTS containers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    container_id VARCHAR(255) NOT NULL,
    image VARCHAR(500) NOT NULL,
    names TEXT[],
    status VARCHAR(100),
    ports TEXT[],
    labels JSONB,
    container_created BIGINT,
    collected_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create browser_sessions table
CREATE TABLE IF NOT EXISTS browser_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    browser_fingerprint VARCHAR(255) NOT NULL,
    user_agent TEXT,
    tabs TEXT[],
    user_id UUID,
    collected_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create threat_findings table
CREATE TABLE IF NOT EXISTS threat_findings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    process_id UUID REFERENCES processes(id) ON DELETE SET NULL,
    container_id UUID REFERENCES containers(id) ON DELETE SET NULL,
    description TEXT NOT NULL,
    severity VARCHAR(20) NOT NULL CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    rule_id VARCHAR(100) NOT NULL,
    rule_name VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_devices_mac_address ON devices(mac_address);
CREATE INDEX IF NOT EXISTS idx_devices_last_seen_at ON devices(last_seen_at);
CREATE INDEX IF NOT EXISTS idx_processes_device_id ON processes(device_id);
CREATE INDEX IF NOT EXISTS idx_processes_collected_at ON processes(collected_at);
CREATE INDEX IF NOT EXISTS idx_processes_pid_device_id ON processes(pid, device_id);
CREATE INDEX IF NOT EXISTS idx_containers_device_id ON containers(device_id);
CREATE INDEX IF NOT EXISTS idx_containers_collected_at ON containers(collected_at);
CREATE INDEX IF NOT EXISTS idx_containers_container_id ON containers(container_id);
CREATE INDEX IF NOT EXISTS idx_browser_sessions_device_id ON browser_sessions(device_id);
CREATE INDEX IF NOT EXISTS idx_browser_sessions_collected_at ON browser_sessions(collected_at);
CREATE INDEX IF NOT EXISTS idx_threat_findings_device_id ON threat_findings(device_id);
CREATE INDEX IF NOT EXISTS idx_threat_findings_severity ON threat_findings(severity);
CREATE INDEX IF NOT EXISTS idx_threat_findings_timestamp ON threat_findings(timestamp);
CREATE INDEX IF NOT EXISTS idx_threat_findings_rule_id ON threat_findings(rule_id);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger for devices table
CREATE TRIGGER update_devices_updated_at BEFORE UPDATE ON devices
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
