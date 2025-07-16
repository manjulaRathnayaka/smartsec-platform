package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"telemetry-service/internal/models"
)

type DeviceRepository struct {
	db *sql.DB
}

func NewDeviceRepository(db *sql.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) CreateOrUpdate(device *models.Device) error {
	query := `
		INSERT INTO devices (id, mac_address, hostname, os, platform, version, current_user, user_id, org_unit, last_seen_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (mac_address) DO UPDATE SET
			hostname = EXCLUDED.hostname,
			os = EXCLUDED.os,
			platform = EXCLUDED.platform,
			version = EXCLUDED.version,
			current_user = EXCLUDED.current_user,
			user_id = EXCLUDED.user_id,
			org_unit = EXCLUDED.org_unit,
			last_seen_at = EXCLUDED.last_seen_at,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, created_at, updated_at`

	if device.ID == "" {
		device.ID = uuid.New().String()
	}

	err := r.db.QueryRow(
		query,
		device.ID,
		device.MacAddress,
		device.Hostname,
		device.OS,
		device.Platform,
		device.Version,
		device.CurrentUser,
		device.UserID,
		device.OrgUnit,
		device.LastSeenAt,
	).Scan(&device.ID, &device.CreatedAt, &device.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create or update device: %w", err)
	}

	return nil
}

func (r *DeviceRepository) GetByMacAddress(macAddress string) (*models.Device, error) {
	query := `
		SELECT id, mac_address, hostname, os, platform, version, current_user, user_id, org_unit,
			   created_at, updated_at, last_seen_at
		FROM devices
		WHERE mac_address = $1`

	device := &models.Device{}
	err := r.db.QueryRow(query, macAddress).Scan(
		&device.ID,
		&device.MacAddress,
		&device.Hostname,
		&device.OS,
		&device.Platform,
		&device.Version,
		&device.CurrentUser,
		&device.UserID,
		&device.OrgUnit,
		&device.CreatedAt,
		&device.UpdatedAt,
		&device.LastSeenAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get device by mac address: %w", err)
	}

	return device, nil
}

func (r *DeviceRepository) GetByID(id string) (*models.Device, error) {
	query := `
		SELECT id, mac_address, hostname, os, platform, version, current_user, user_id, org_unit,
			   created_at, updated_at, last_seen_at
		FROM devices
		WHERE id = $1`

	device := &models.Device{}
	err := r.db.QueryRow(query, id).Scan(
		&device.ID,
		&device.MacAddress,
		&device.Hostname,
		&device.OS,
		&device.Platform,
		&device.Version,
		&device.CurrentUser,
		&device.UserID,
		&device.OrgUnit,
		&device.CreatedAt,
		&device.UpdatedAt,
		&device.LastSeenAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get device by ID: %w", err)
	}

	return device, nil
}

func (r *DeviceRepository) List(limit, offset int) ([]*models.Device, error) {
	query := `
		SELECT id, mac_address, hostname, os, platform, version, current_user, user_id, org_unit,
			   created_at, updated_at, last_seen_at
		FROM devices
		ORDER BY last_seen_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list devices: %w", err)
	}
	defer rows.Close()

	var devices []*models.Device
	for rows.Next() {
		device := &models.Device{}
		err := rows.Scan(
			&device.ID,
			&device.MacAddress,
			&device.Hostname,
			&device.OS,
			&device.Platform,
			&device.Version,
			&device.CurrentUser,
			&device.UserID,
			&device.OrgUnit,
			&device.CreatedAt,
			&device.UpdatedAt,
			&device.LastSeenAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan device row: %w", err)
		}
		devices = append(devices, device)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating device rows: %w", err)
	}

	return devices, nil
}

type ProcessRepository struct {
	db *sql.DB
}

func NewProcessRepository(db *sql.DB) *ProcessRepository {
	return &ProcessRepository{db: db}
}

func (r *ProcessRepository) Create(process *models.Process) error {
	query := `
		INSERT INTO processes (id, device_id, pid, name, cmdline, username, exe_path, start_time, status, sha256, version, file_size, collected_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING created_at`

	if process.ID == "" {
		process.ID = uuid.New().String()
	}

	err := r.db.QueryRow(
		query,
		process.ID,
		process.DeviceID,
		process.PID,
		process.Name,
		pq.Array(process.Cmdline),
		process.Username,
		process.ExePath,
		process.StartTime,
		process.Status,
		process.SHA256,
		process.Version,
		process.FileSize,
		process.CollectedAt,
	).Scan(&process.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create process: %w", err)
	}

	return nil
}

func (r *ProcessRepository) GetByDeviceID(deviceID string, limit, offset int) ([]*models.Process, error) {
	query := `
		SELECT id, device_id, pid, name, cmdline, username, exe_path, start_time, status, sha256, version, file_size, collected_at, created_at
		FROM processes
		WHERE device_id = $1
		ORDER BY collected_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, deviceID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get processes by device ID: %w", err)
	}
	defer rows.Close()

	var processes []*models.Process
	for rows.Next() {
		process := &models.Process{}
		err := rows.Scan(
			&process.ID,
			&process.DeviceID,
			&process.PID,
			&process.Name,
			pq.Array(&process.Cmdline),
			&process.Username,
			&process.ExePath,
			&process.StartTime,
			&process.Status,
			&process.SHA256,
			&process.Version,
			&process.FileSize,
			&process.CollectedAt,
			&process.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan process row: %w", err)
		}
		processes = append(processes, process)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating process rows: %w", err)
	}

	return processes, nil
}

func (r *ProcessRepository) DeleteOldProcesses(deviceID string, before time.Time) error {
	query := `DELETE FROM processes WHERE device_id = $1 AND collected_at < $2`

	_, err := r.db.Exec(query, deviceID, before)
	if err != nil {
		return fmt.Errorf("failed to delete old processes: %w", err)
	}

	return nil
}
