package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Sample telemetry data structure
type TelemetryData struct {
	Timestamp    time.Time       `json:"timestamp"`
	MacAddress   string          `json:"mac_address"`
	HostMetadata HostMetadata    `json:"host_metadata"`
	Processes    []ProcessInfo   `json:"processes"`
	Containers   []ContainerInfo `json:"containers"`
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
	SHA256    string   `json:"sha256"`
	Version   string   `json:"version"`
	FileSize  int64    `json:"file_size"`
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

func main() {
	// Create sample telemetry data
	telemetryData := TelemetryData{
		Timestamp:  time.Now(),
		MacAddress: "aa:bb:cc:dd:ee:ff",
		HostMetadata: HostMetadata{
			Hostname:    "test-laptop",
			OS:          "darwin",
			Platform:    "darwin",
			Version:     "14.0.0",
			CurrentUser: "testuser",
			Uptime:      3600,
		},
		Processes: []ProcessInfo{
			{
				PID:       1234,
				Name:      "chrome",
				Cmdline:   []string{"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"},
				Username:  "testuser",
				ExePath:   "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
				StartTime: time.Now().Add(-time.Hour).Unix(),
				Status:    "running",
				SHA256:    "abc123def456ghi789jkl012mno345pqr678stu901vwx234yz567890abcdef12",
				Version:   "120.0.0",
				FileSize:  1024000,
			},
			{
				PID:       5678,
				Name:      "terminal",
				Cmdline:   []string{"/Applications/Utilities/Terminal.app/Contents/MacOS/Terminal"},
				Username:  "testuser",
				ExePath:   "/Applications/Utilities/Terminal.app/Contents/MacOS/Terminal",
				StartTime: time.Now().Add(-30 * time.Minute).Unix(),
				Status:    "running",
				SHA256:    "def456ghi789jkl012mno345pqr678stu901vwx234yz567890abcdef12abc123",
				Version:   "2.12.7",
				FileSize:  2048000,
			},
		},
		Containers: []ContainerInfo{
			{
				ID:      "container123456",
				Image:   "nginx:latest",
				Names:   []string{"/nginx-web"},
				Status:  "running",
				Ports:   []string{"80:8080"},
				Labels:  map[string]string{"app": "web", "environment": "development"},
				Created: time.Now().Add(-2 * time.Hour).Unix(),
			},
			{
				ID:      "container789012",
				Image:   "postgres:15",
				Names:   []string{"/postgres-db"},
				Status:  "running",
				Ports:   []string{"5432:5432"},
				Labels:  map[string]string{"app": "database", "environment": "development"},
				Created: time.Now().Add(-3 * time.Hour).Unix(),
			},
		},
	}

	// Send telemetry to service
	err := sendTelemetry("http://localhost:8080/api/telemetry", telemetryData)
	if err != nil {
		fmt.Printf("Error sending telemetry: %v\n", err)
		return
	}

	fmt.Println("Telemetry sent successfully!")

	// Query telemetry data
	err = queryTelemetry("http://localhost:8080/api/telemetry?device_id=test&type=all&limit=10")
	if err != nil {
		fmt.Printf("Error querying telemetry: %v\n", err)
		return
	}

	// Query devices
	err = queryDevices("http://localhost:8080/api/devices")
	if err != nil {
		fmt.Printf("Error querying devices: %v\n", err)
		return
	}
}

func sendTelemetry(endpoint string, data TelemetryData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	fmt.Printf("Response: %s\n", string(body))
	return nil
}

func queryTelemetry(endpoint string) error {
	fmt.Printf("\nQuerying telemetry data: %s\n", endpoint)

	resp, err := http.Get(endpoint)
	if err != nil {
		return fmt.Errorf("failed to query telemetry: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	fmt.Printf("Telemetry Query Response: %s\n", string(body))
	return nil
}

func queryDevices(endpoint string) error {
	fmt.Printf("\nQuerying devices: %s\n", endpoint)

	resp, err := http.Get(endpoint)
	if err != nil {
		return fmt.Errorf("failed to query devices: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	fmt.Printf("Devices Query Response: %s\n", string(body))
	return nil
}
