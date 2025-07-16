package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

type TelemetryData struct {
	Timestamp    time.Time       `json:"timestamp"`
	HostMetadata HostMetadata    `json:"host_metadata"`
	Processes    []ProcessInfo   `json:"processes"`
	Containers   []ContainerInfo `json:"containers"`
	MacAddress   string          `json:"mac_address"`
}

type HostMetadata struct {
	Hostname    string `json:"hostname"`
	OS          string `json:"os"`
	Platform    string `json:"platform"`
	Version     string `json:"version"`
	CurrentUser string `json:"current_user"`
	Uptime      uint64 `json:"uptime"`
}

type ProcessInfo struct {
	PID       int32    `json:"pid"`
	Name      string   `json:"name"`
	Cmdline   []string `json:"cmdline"`
	Username  string   `json:"username"`
	ExePath   string   `json:"exe_path"`
	StartTime int64    `json:"start_time"`
	Status    string   `json:"status"`
	SHA256    string   `json:"sha256,omitempty"`
	Version   string   `json:"version,omitempty"`
	FileSize  int64    `json:"file_size,omitempty"`
}

type ContainerInfo struct {
	ID      string            `json:"id"`
	Image   string            `json:"image"`
	Names   []string          `json:"names"`
	Status  string            `json:"status"`
	Ports   []string          `json:"ports"`
	Labels  map[string]string `json:"labels"`
	Created int64             `json:"created"`
}

type Config struct {
	APIEndpoint        string
	CollectionInterval time.Duration
	LogOnly            bool
}

type ProcessFileInfo struct {
	SHA256   string
	FileSize int64
	Version  string
}

func main() {
	config := Config{
		APIEndpoint:        getEnvOrDefault("API_ENDPOINT", "http://localhost:8080/api/telemetry"),
		CollectionInterval: time.Duration(getEnvOrDefaultInt("COLLECTION_INTERVAL", 60)) * time.Second,
		LogOnly:            getEnvOrDefault("LOG_ONLY", "false") == "true",
	}

	if config.LogOnly {
		log.Printf("Starting laptop agent in LOG_ONLY mode with collection interval: %v", config.CollectionInterval)
	} else {
		log.Printf("Starting laptop agent with collection interval: %v, endpoint: %s", config.CollectionInterval, config.APIEndpoint)
	}

	ticker := time.NewTicker(config.CollectionInterval)
	defer ticker.Stop()

	// Collect and send initial telemetry
	collectAndSend(config)

	// Continue collecting at intervals
	for range ticker.C {
		collectAndSend(config)
	}
}

func collectAndSend(config Config) {
	log.Println("Collecting telemetry data...")

	telemetry := TelemetryData{
		Timestamp: time.Now(),
	}

	// Collect host metadata
	hostInfo, err := collectHostMetadata()
	if err != nil {
		log.Printf("Error collecting host metadata: %v", err)
	} else {
		telemetry.HostMetadata = hostInfo
	}

	// Collect MAC address
	macAddr, err := getMacAddress()
	if err != nil {
		log.Printf("Error getting MAC address: %v", err)
	} else {
		telemetry.MacAddress = macAddr
	}

	// Collect processes
	processes, err := collectProcesses()
	if err != nil {
		log.Printf("Error collecting processes: %v", err)
	} else {
		telemetry.Processes = processes
	}

	// Collect containers
	containers, err := collectContainers()
	if err != nil {
		log.Printf("Error collecting containers: %v", err)
	} else {
		telemetry.Containers = containers
	}

	// Send telemetry
	err = sendTelemetryOrLog(config, telemetry)
	if err != nil {
		log.Printf("Error sending telemetry: %v", err)
	} else {
		if config.LogOnly {
			log.Printf("Telemetry logged successfully: %d processes, %d containers",
				len(telemetry.Processes), len(telemetry.Containers))
		} else {
			log.Printf("Telemetry sent successfully: %d processes, %d containers",
				len(telemetry.Processes), len(telemetry.Containers))
		}
	}
}

func collectHostMetadata() (HostMetadata, error) {
	info, err := host.Info()
	if err != nil {
		return HostMetadata{}, err
	}

	currentUser := os.Getenv("USER")
	if currentUser == "" {
		currentUser = os.Getenv("USERNAME")
	}

	return HostMetadata{
		Hostname:    info.Hostname,
		OS:          info.OS,
		Platform:    info.Platform,
		Version:     info.PlatformVersion,
		CurrentUser: currentUser,
		Uptime:      info.Uptime,
	}, nil
}

func getMacAddress() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		// Skip loopback interfaces and get first available MAC
		if iface.Name != "lo" && iface.HardwareAddr != "" {
			return iface.HardwareAddr, nil
		}
	}

	return "", fmt.Errorf("no MAC address found")
}

func collectProcesses() ([]ProcessInfo, error) {
	pids, err := process.Pids()
	if err != nil {
		return nil, err
	}

	var processes []ProcessInfo
	processedFiles := make(map[string]ProcessFileInfo) // Cache for file info

	log.Printf("Collecting information for %d processes...", len(pids))
	processCount := 0

	for _, pid := range pids {
		proc, err := process.NewProcess(pid)
		if err != nil {
			continue // Skip processes we can't access
		}

		name, _ := proc.Name()
		cmdlineStr, _ := proc.Cmdline()
		username, _ := proc.Username()
		exe, _ := proc.Exe()
		createTime, _ := proc.CreateTime()
		status, _ := proc.Status()

		// Convert cmdline string to slice
		var cmdline []string
		if cmdlineStr != "" {
			cmdline = []string{cmdlineStr}
		}

		// Convert status slice to string
		var statusStr string
		if len(status) > 0 {
			statusStr = status[0]
		}

		processInfo := ProcessInfo{
			PID:       pid,
			Name:      name,
			Cmdline:   cmdline,
			Username:  username,
			ExePath:   exe,
			StartTime: createTime,
			Status:    statusStr,
		}

		// Enhance with security information (with caching and timeout)
		if exe != "" && exe != "0" {
			if fileInfo, exists := processedFiles[exe]; exists {
				// Use cached information
				processInfo.SHA256 = fileInfo.SHA256
				processInfo.FileSize = fileInfo.FileSize
				processInfo.Version = fileInfo.Version
			} else {
				// Collect file information with timeout
				fileInfo := collectFileInfoWithTimeout(exe, 2*time.Second)
				processedFiles[exe] = fileInfo
				processInfo.SHA256 = fileInfo.SHA256
				processInfo.FileSize = fileInfo.FileSize
				processInfo.Version = fileInfo.Version
			}
		}

		processes = append(processes, processInfo)
		processCount++

		// Log progress every 100 processes
		if processCount%100 == 0 {
			log.Printf("Processed %d/%d processes...", processCount, len(pids))
		}
	}

	log.Printf("Finished collecting information for %d processes", len(processes))
	return processes, nil
}

func collectContainers() ([]ContainerInfo, error) {
	// Try to connect to Docker daemon with different socket paths
	var cli *client.Client
	var err error

	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = os.Getenv("HOME")
		if homeDir == "" {
			homeDir = "/tmp" // fallback
		}
	}

	// Try common Docker socket locations
	socketPaths := []string{
		"unix:///var/run/docker.sock",                                                           // Standard Docker (Linux)
		fmt.Sprintf("unix://%s/.rd/docker.sock", homeDir),                                       // Rancher Desktop
		fmt.Sprintf("unix://%s/.docker/run/docker.sock", homeDir),                               // Docker Desktop
		fmt.Sprintf("unix://%s/.docker/desktop/docker.sock", homeDir),                           // Docker Desktop alternative
		fmt.Sprintf("unix://%s/Library/Containers/com.docker.docker/Data/docker.sock", homeDir), // Docker Desktop macOS
		"unix:///tmp/docker.sock",                                                               // Alternative location
		"unix:///usr/local/var/run/docker.sock",                                                 // Homebrew Docker
	}

	// Filter socket paths to only include those that actually exist
	var validSocketPaths []string
	for _, socketPath := range socketPaths {
		// Extract the file path from the unix:// prefix
		filePath := strings.TrimPrefix(socketPath, "unix://")
		if _, err := os.Stat(filePath); err == nil {
			validSocketPaths = append(validSocketPaths, socketPath)
		}
	}

	// First try with environment variables
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err == nil {
		// Test connection with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		_, pingErr := cli.Ping(ctx)
		cancel()
		if pingErr == nil {
			log.Printf("Connected to Docker daemon via environment settings")
		} else {
			cli.Close()
			cli = nil
		}
	}

	// If environment connection failed, try different socket paths
	if cli == nil {
		// Use valid socket paths if any were found, otherwise fall back to all paths
		pathsToTry := validSocketPaths
		if len(pathsToTry) == 0 {
			pathsToTry = socketPaths
			log.Printf("No existing Docker socket files found, trying all paths")
		} else {
			log.Printf("Found %d existing Docker socket files", len(pathsToTry))
		}

		for _, socketPath := range pathsToTry {
			cli, err = client.NewClientWithOpts(client.WithHost(socketPath), client.WithAPIVersionNegotiation())
			if err != nil {
				log.Printf("Failed to create Docker client for %s: %v", socketPath, err)
				continue
			}

			// Test connection with timeout
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			_, pingErr := cli.Ping(ctx)
			cancel()

			if pingErr == nil {
				log.Printf("Connected to Docker daemon at: %s", socketPath)
				break
			} else {
				log.Printf("Failed to connect to Docker daemon at %s: %v", socketPath, pingErr)
				cli.Close()
				cli = nil
			}
		}
	}

	if cli == nil {
		log.Printf("Docker daemon not accessible at any known location")
		return []ContainerInfo{}, nil // Return empty slice instead of nil
	}
	defer cli.Close()

	// Test connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = cli.Ping(ctx)
	if err != nil {
		log.Printf("Docker daemon not accessible: %v", err)
		return []ContainerInfo{}, nil // Return empty slice instead of nil
	}

	containers, err := cli.ContainerList(ctx, container.ListOptions{
		All: true,
	})
	if err != nil {
		log.Printf("Failed to list containers: %v", err)
		return []ContainerInfo{}, nil // Return empty slice instead of nil
	}

	var containerInfos []ContainerInfo

	for _, container := range containers {
		var ports []string
		for _, port := range container.Ports {
			if port.PublicPort != 0 {
				ports = append(ports, fmt.Sprintf("%d:%d", port.PublicPort, port.PrivatePort))
			}
		}

		containerInfo := ContainerInfo{
			ID:      container.ID,
			Image:   container.Image,
			Names:   container.Names,
			Status:  container.Status,
			Ports:   ports,
			Labels:  container.Labels,
			Created: container.Created,
		}

		containerInfos = append(containerInfos, containerInfo)
	}

	log.Printf("Successfully collected %d containers", len(containerInfos))
	return containerInfos, nil
}

func sendTelemetry(endpoint string, data TelemetryData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal telemetry data: %w", err)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "SmartSec-Laptop-Agent/1.0")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status: %d", resp.StatusCode)
	}

	return nil
}

func sendTelemetryOrLog(config Config, data TelemetryData) error {
	if config.LogOnly {
		// Log the telemetry data instead of sending it
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal telemetry data for logging: %w", err)
		}

		log.Printf("=== TELEMETRY DATA ===\n%s\n=== END TELEMETRY ===", string(jsonData))
		return nil
	}

	return sendTelemetry(config.APIEndpoint, data)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvOrDefaultInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := time.ParseDuration(value + "s"); err == nil {
			return int(intValue.Seconds())
		}
	}
	return defaultValue
}

// getFileSHA256 calculates the SHA256 hash of a file
func getFileSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// getFileSize returns the size of a file in bytes
func getFileSize(filePath string) (int64, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// getFileVersion attempts to extract version information from a file
func getFileVersion(filePath string) (string, error) {
	// Try different approaches based on the OS and file type

	// For macOS/Linux, try to get version from the binary itself
	if version := getMacOSVersion(filePath); version != "" {
		return version, nil
	}

	// Try to get version from file command
	if version := getVersionFromFile(filePath); version != "" {
		return version, nil
	}

	// Try to get version from strings in the binary
	if version := getVersionFromStrings(filePath); version != "" {
		return version, nil
	}

	return "", fmt.Errorf("version information not available")
}

// getMacOSVersion tries to extract version from macOS binaries
func getMacOSVersion(filePath string) string {
	// Try otool command for macOS binaries
	cmd := exec.Command("otool", "-l", filePath)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "version") {
			parts := strings.Fields(line)
			for i, part := range parts {
				if strings.Contains(part, "version") && i+1 < len(parts) {
					return parts[i+1]
				}
			}
		}
	}

	return ""
}

// getVersionFromFile tries to get version using the file command
func getVersionFromFile(filePath string) string {
	cmd := exec.Command("file", filePath)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	outputStr := string(output)
	// Look for version patterns in file output
	if strings.Contains(outputStr, "version") {
		parts := strings.Split(outputStr, " ")
		for i, part := range parts {
			if strings.Contains(part, "version") && i+1 < len(parts) {
				return parts[i+1]
			}
		}
	}

	return ""
}

// getVersionFromStrings tries to extract version from strings in the binary
func getVersionFromStrings(filePath string) string {
	cmd := exec.Command("strings", filePath)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Look for common version patterns
		if strings.Contains(strings.ToLower(line), "version") && len(line) < 50 {
			// Extract version number patterns
			if strings.Contains(line, ".") && (strings.Contains(line, "v") || strings.Contains(line, "V")) {
				return line
			}
		}
		// Look for semantic version patterns (x.y.z)
		if len(line) > 3 && len(line) < 20 {
			parts := strings.Split(line, ".")
			if len(parts) >= 3 {
				// Check if it looks like a version number
				hasDigits := false
				for _, part := range parts {
					if len(part) > 0 && part[0] >= '0' && part[0] <= '9' {
						hasDigits = true
						break
					}
				}
				if hasDigits {
					return line
				}
			}
		}
	}

	return ""
}

// collectFileInfoWithTimeout collects file information with a timeout to prevent hanging
func collectFileInfoWithTimeout(filePath string, timeout time.Duration) ProcessFileInfo {
	result := ProcessFileInfo{}

	// Create a channel to receive the result
	done := make(chan struct{})

	go func() {
		defer close(done)

		// Get SHA256 hash
		if sha256Hash, err := getFileSHA256(filePath); err == nil {
			result.SHA256 = sha256Hash
		}

		// Get file size
		if fileSize, err := getFileSize(filePath); err == nil {
			result.FileSize = fileSize
		}

		// Get version information (this is the most time-consuming part)
		if version, err := getFileVersion(filePath); err == nil && version != "" {
			result.Version = version
		}
	}()

	// Wait for completion or timeout
	select {
	case <-done:
		return result
	case <-time.After(timeout):
		// Timeout occurred, return what we have so far
		return result
	}
}
