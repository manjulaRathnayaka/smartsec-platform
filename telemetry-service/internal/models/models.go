package models

import (
	"time"
)

// Device represents a device in the system
type Device struct {
	ID          string    `json:"id" db:"id"`
	MacAddress  string    `json:"mac_address" db:"mac_address"`
	Hostname    string    `json:"hostname" db:"hostname"`
	OS          string    `json:"os" db:"os"`
	Platform    string    `json:"platform" db:"platform"`
	Version     string    `json:"version" db:"version"`
	CurrentUser string    `json:"current_user" db:"current_username"`
	UserID      *string   `json:"user_id" db:"user_id"`
	OrgUnit     *string   `json:"org_unit" db:"org_unit"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	LastSeenAt  time.Time `json:"last_seen_at" db:"last_seen_at"`
}

// Process represents a process record
type Process struct {
	ID          string    `json:"id" db:"id"`
	DeviceID    string    `json:"device_id" db:"device_id"`
	PID         int32     `json:"pid" db:"pid"`
	Name        string    `json:"name" db:"name"`
	Cmdline     []string  `json:"cmdline" db:"cmdline"`
	Username    string    `json:"username" db:"username"`
	ExePath     string    `json:"exe_path" db:"exe_path"`
	StartTime   int64     `json:"start_time" db:"start_time"`
	Status      string    `json:"status" db:"status"`
	SHA256      string    `json:"sha256" db:"sha256"`
	Version     string    `json:"version" db:"version"`
	FileSize    int64     `json:"file_size" db:"file_size"`
	CollectedAt time.Time `json:"collected_at" db:"collected_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// Container represents a container record
type Container struct {
	ID               string            `json:"id" db:"id"`
	DeviceID         string            `json:"device_id" db:"device_id"`
	ContainerID      string            `json:"container_id" db:"container_id"`
	Image            string            `json:"image" db:"image"`
	Names            []string          `json:"names" db:"names"`
	Status           string            `json:"status" db:"status"`
	Ports            []string          `json:"ports" db:"ports"`
	Labels           map[string]string `json:"labels" db:"labels"`
	ContainerCreated int64             `json:"container_created" db:"container_created"`
	CollectedAt      time.Time         `json:"collected_at" db:"collected_at"`
	CreatedAt        time.Time         `json:"created_at" db:"created_at"`
}

// BrowserSession represents a browser session record
type BrowserSession struct {
	ID                 string    `json:"id" db:"id"`
	DeviceID           string    `json:"device_id" db:"device_id"`
	BrowserFingerprint string    `json:"browser_fingerprint" db:"browser_fingerprint"`
	UserAgent          string    `json:"user_agent" db:"user_agent"`
	Tabs               []string  `json:"tabs" db:"tabs"`
	UserID             *string   `json:"user_id" db:"user_id"`
	CollectedAt        time.Time `json:"collected_at" db:"collected_at"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}

// ThreatFinding represents a threat finding record
type ThreatFinding struct {
	ID          string    `json:"id" db:"id"`
	DeviceID    string    `json:"device_id" db:"device_id"`
	ProcessID   *string   `json:"process_id" db:"process_id"`
	ContainerID *string   `json:"container_id" db:"container_id"`
	Description string    `json:"description" db:"description"`
	Severity    string    `json:"severity" db:"severity"`
	RuleID      string    `json:"rule_id" db:"rule_id"`
	RuleName    string    `json:"rule_name" db:"rule_name"`
	Timestamp   time.Time `json:"timestamp" db:"timestamp"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// TelemetryRequest represents the incoming telemetry request
type TelemetryRequest struct {
	Timestamp    time.Time       `json:"timestamp" validate:"required"`
	HostMetadata HostMetadata    `json:"host_metadata" validate:"required"`
	Processes    []ProcessInfo   `json:"processes"`
	Containers   []ContainerInfo `json:"containers"`
	MacAddress   string          `json:"mac_address" validate:"required"`
}

// HostMetadata represents host metadata from the agent
type HostMetadata struct {
	Hostname    string `json:"hostname" validate:"required"`
	OS          string `json:"os" validate:"required"`
	Platform    string `json:"platform" validate:"required"`
	Version     string `json:"version" validate:"required"`
	CurrentUser string `json:"current_user" validate:"required"`
	Uptime      uint64 `json:"uptime"`
}

// ProcessInfo represents process information from the agent
type ProcessInfo struct {
	PID       int32    `json:"pid" validate:"required"`
	Name      string   `json:"name" validate:"required"`
	Cmdline   []string `json:"cmdline"`
	Username  string   `json:"username"`
	ExePath   string   `json:"exe_path"`
	StartTime int64    `json:"start_time"`
	Status    string   `json:"status"`
	SHA256    string   `json:"sha256"`
	Version   string   `json:"version"`
	FileSize  int64    `json:"file_size"`
}

// ContainerInfo represents container information from the agent
type ContainerInfo struct {
	ID      string            `json:"id" validate:"required"`
	Image   string            `json:"image" validate:"required"`
	Names   []string          `json:"names"`
	Status  string            `json:"status"`
	Ports   []string          `json:"ports"`
	Labels  map[string]string `json:"labels"`
	Created int64             `json:"created"`
}
