package service

import (
	"fmt"
	"time"

	"telemetry-service/internal/models"
	"telemetry-service/internal/repository"
)

type TelemetryService struct {
	deviceRepo    *repository.DeviceRepository
	processRepo   *repository.ProcessRepository
	containerRepo *repository.ContainerRepository
	threatRepo    *repository.ThreatRepository
}

func NewTelemetryService(
	deviceRepo *repository.DeviceRepository,
	processRepo *repository.ProcessRepository,
	containerRepo *repository.ContainerRepository,
	threatRepo *repository.ThreatRepository,
) *TelemetryService {
	return &TelemetryService{
		deviceRepo:    deviceRepo,
		processRepo:   processRepo,
		containerRepo: containerRepo,
		threatRepo:    threatRepo,
	}
}

func (s *TelemetryService) ProcessTelemetry(req *models.TelemetryRequest) error {
	// First, create or update the device
	device := &models.Device{
		MacAddress:  req.MacAddress,
		Hostname:    req.HostMetadata.Hostname,
		OS:          req.HostMetadata.OS,
		Platform:    req.HostMetadata.Platform,
		Version:     req.HostMetadata.Version,
		CurrentUser: req.HostMetadata.CurrentUser,
		LastSeenAt:  req.Timestamp,
	}

	err := s.deviceRepo.CreateOrUpdate(device)
	if err != nil {
		return fmt.Errorf("failed to create or update device: %w", err)
	}

	// Clean up old data (keep only last 7 days)
	cleanupThreshold := time.Now().AddDate(0, 0, -7)

	// Clean up old processes
	err = s.processRepo.DeleteOldProcesses(device.ID, cleanupThreshold)
	if err != nil {
		return fmt.Errorf("failed to clean up old processes: %w", err)
	}

	// Clean up old containers
	err = s.containerRepo.DeleteOldContainers(device.ID, cleanupThreshold)
	if err != nil {
		return fmt.Errorf("failed to clean up old containers: %w", err)
	}

	// Process the processes
	for _, processInfo := range req.Processes {
		process := &models.Process{
			DeviceID:    device.ID,
			PID:         processInfo.PID,
			Name:        processInfo.Name,
			Cmdline:     processInfo.Cmdline,
			Username:    processInfo.Username,
			ExePath:     processInfo.ExePath,
			StartTime:   processInfo.StartTime,
			Status:      processInfo.Status,
			SHA256:      processInfo.SHA256,
			Version:     processInfo.Version,
			FileSize:    processInfo.FileSize,
			CollectedAt: req.Timestamp,
		}

		err := s.processRepo.Create(process)
		if err != nil {
			return fmt.Errorf("failed to create process record: %w", err)
		}
	}

	// Process the containers
	for _, containerInfo := range req.Containers {
		container := &models.Container{
			DeviceID:         device.ID,
			ContainerID:      containerInfo.ID,
			Image:            containerInfo.Image,
			Names:            containerInfo.Names,
			Status:           containerInfo.Status,
			Ports:            containerInfo.Ports,
			Labels:           containerInfo.Labels,
			ContainerCreated: containerInfo.Created,
			CollectedAt:      req.Timestamp,
		}

		err := s.containerRepo.Create(container)
		if err != nil {
			return fmt.Errorf("failed to create container record: %w", err)
		}
	}

	return nil
}

func (s *TelemetryService) GetDevices(limit, offset int) ([]*models.Device, error) {
	return s.deviceRepo.List(limit, offset)
}

func (s *TelemetryService) GetDevice(deviceID string) (*models.Device, error) {
	return s.deviceRepo.GetByID(deviceID)
}

func (s *TelemetryService) GetProcesses(deviceID string, limit, offset int) ([]*models.Process, error) {
	return s.processRepo.GetByDeviceID(deviceID, limit, offset)
}

func (s *TelemetryService) GetContainers(deviceID string, limit, offset int) ([]*models.Container, error) {
	return s.containerRepo.GetByDeviceID(deviceID, limit, offset)
}

func (s *TelemetryService) GetThreatFindings(deviceID string, limit, offset int) ([]*models.ThreatFinding, error) {
	return s.threatRepo.GetByDeviceID(deviceID, limit, offset)
}

func (s *TelemetryService) GetThreatFindingsBySeverity(severity string, limit, offset int) ([]*models.ThreatFinding, error) {
	return s.threatRepo.GetBySeverity(severity, limit, offset)
}

func (s *TelemetryService) CreateThreatFinding(threat *models.ThreatFinding) error {
	return s.threatRepo.Create(threat)
}
