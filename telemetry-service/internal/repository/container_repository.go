package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"telemetry-service/internal/models"
)

type ContainerRepository struct {
	db *sql.DB
}

func NewContainerRepository(db *sql.DB) *ContainerRepository {
	return &ContainerRepository{db: db}
}

func (r *ContainerRepository) Create(container *models.Container) error {
	// Convert labels map to JSON
	labelsJSON, err := json.Marshal(container.Labels)
	if err != nil {
		return fmt.Errorf("failed to marshal labels: %w", err)
	}

	query := `
		INSERT INTO containers (id, device_id, container_id, image, names, status, ports, labels, container_created, collected_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING created_at`

	if container.ID == "" {
		container.ID = uuid.New().String()
	}

	err = r.db.QueryRow(
		query,
		container.ID,
		container.DeviceID,
		container.ContainerID,
		container.Image,
		pq.Array(container.Names),
		container.Status,
		pq.Array(container.Ports),
		labelsJSON,
		container.ContainerCreated,
		container.CollectedAt,
	).Scan(&container.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	return nil
}

func (r *ContainerRepository) GetByDeviceID(deviceID string, limit, offset int) ([]*models.Container, error) {
	query := `
		SELECT id, device_id, container_id, image, names, status, ports, labels, container_created, collected_at, created_at
		FROM containers
		WHERE device_id = $1
		ORDER BY collected_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, deviceID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get containers by device ID: %w", err)
	}
	defer rows.Close()

	var containers []*models.Container
	for rows.Next() {
		container := &models.Container{}
		var labelsJSON []byte

		err := rows.Scan(
			&container.ID,
			&container.DeviceID,
			&container.ContainerID,
			&container.Image,
			pq.Array(&container.Names),
			&container.Status,
			pq.Array(&container.Ports),
			&labelsJSON,
			&container.ContainerCreated,
			&container.CollectedAt,
			&container.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan container row: %w", err)
		}

		// Unmarshal labels JSON
		if len(labelsJSON) > 0 {
			err := json.Unmarshal(labelsJSON, &container.Labels)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal labels: %w", err)
			}
		}

		containers = append(containers, container)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating container rows: %w", err)
	}

	return containers, nil
}

func (r *ContainerRepository) DeleteOldContainers(deviceID string, before time.Time) error {
	query := `DELETE FROM containers WHERE device_id = $1 AND collected_at < $2`

	_, err := r.db.Exec(query, deviceID, before)
	if err != nil {
		return fmt.Errorf("failed to delete old containers: %w", err)
	}

	return nil
}

type ThreatRepository struct {
	db *sql.DB
}

func NewThreatRepository(db *sql.DB) *ThreatRepository {
	return &ThreatRepository{db: db}
}

func (r *ThreatRepository) Create(threat *models.ThreatFinding) error {
	query := `
		INSERT INTO threat_findings (id, device_id, process_id, container_id, description, severity, rule_id, rule_name, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at`

	if threat.ID == "" {
		threat.ID = uuid.New().String()
	}

	err := r.db.QueryRow(
		query,
		threat.ID,
		threat.DeviceID,
		threat.ProcessID,
		threat.ContainerID,
		threat.Description,
		threat.Severity,
		threat.RuleID,
		threat.RuleName,
		threat.Timestamp,
	).Scan(&threat.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create threat finding: %w", err)
	}

	return nil
}

func (r *ThreatRepository) GetByDeviceID(deviceID string, limit, offset int) ([]*models.ThreatFinding, error) {
	query := `
		SELECT id, device_id, process_id, container_id, description, severity, rule_id, rule_name, timestamp, created_at
		FROM threat_findings
		WHERE device_id = $1
		ORDER BY timestamp DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, deviceID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get threat findings by device ID: %w", err)
	}
	defer rows.Close()

	var threats []*models.ThreatFinding
	for rows.Next() {
		threat := &models.ThreatFinding{}
		err := rows.Scan(
			&threat.ID,
			&threat.DeviceID,
			&threat.ProcessID,
			&threat.ContainerID,
			&threat.Description,
			&threat.Severity,
			&threat.RuleID,
			&threat.RuleName,
			&threat.Timestamp,
			&threat.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan threat finding row: %w", err)
		}
		threats = append(threats, threat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating threat finding rows: %w", err)
	}

	return threats, nil
}

func (r *ThreatRepository) GetBySeverity(severity string, limit, offset int) ([]*models.ThreatFinding, error) {
	query := `
		SELECT id, device_id, process_id, container_id, description, severity, rule_id, rule_name, timestamp, created_at
		FROM threat_findings
		WHERE severity = $1
		ORDER BY timestamp DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, severity, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get threat findings by severity: %w", err)
	}
	defer rows.Close()

	var threats []*models.ThreatFinding
	for rows.Next() {
		threat := &models.ThreatFinding{}
		err := rows.Scan(
			&threat.ID,
			&threat.DeviceID,
			&threat.ProcessID,
			&threat.ContainerID,
			&threat.Description,
			&threat.Severity,
			&threat.RuleID,
			&threat.RuleName,
			&threat.Timestamp,
			&threat.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan threat finding row: %w", err)
		}
		threats = append(threats, threat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating threat finding rows: %w", err)
	}

	return threats, nil
}
