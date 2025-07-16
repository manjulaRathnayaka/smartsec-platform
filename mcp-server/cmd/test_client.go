package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// TestClient for testing the MCP server
type TestClient struct {
	baseURL string
	client  *http.Client
}

// NewTestClient creates a new test client
func NewTestClient(baseURL string) *TestClient {
	return &TestClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// TestSchema tests the schema endpoint
func (tc *TestClient) TestSchema() error {
	fmt.Println("Testing GET /mcp/schema...")
	
	resp, err := tc.client.Get(tc.baseURL + "/mcp/schema")
	if err != nil {
		return fmt.Errorf("failed to get schema: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	
	var schema map[string]interface{}
	if err := json.Unmarshal(body, &schema); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w", err)
	}
	
	fmt.Printf("✓ Schema retrieved successfully (size: %d bytes)\n", len(body))
	return nil
}

// TestEntities tests the entities endpoint
func (tc *TestClient) TestEntities() error {
	fmt.Println("Testing GET /mcp/entities...")
	
	resp, err := tc.client.Get(tc.baseURL + "/mcp/entities")
	if err != nil {
		return fmt.Errorf("failed to get entities: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	
	var entities map[string]interface{}
	if err := json.Unmarshal(body, &entities); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w", err)
	}
	
	fmt.Printf("✓ Entities retrieved successfully\n")
	return nil
}

// TestQuery tests the query endpoint
func (tc *TestClient) TestQuery() error {
	fmt.Println("Testing POST /mcp/query...")
	
	queryReq := map[string]interface{}{
		"entity": "devices",
		"fields": []string{"id", "hostname", "os"},
		"filters": []map[string]interface{}{
			{
				"field":    "os",
				"operator": "eq",
				"value":    "Linux",
			},
		},
		"limit": 10,
	}
	
	jsonData, err := json.Marshal(queryReq)
	if err != nil {
		return fmt.Errorf("failed to marshal query request: %w", err)
	}
	
	resp, err := tc.client.Post(tc.baseURL+"/mcp/query", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to post query: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	
	var queryResp map[string]interface{}
	if err := json.Unmarshal(body, &queryResp); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w", err)
	}
	
	status, ok := queryResp["status"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid status in response")
	}
	
	if status == "failed" {
		if errorMsg, exists := queryResp["error"]; exists {
			return fmt.Errorf("query failed: %v", errorMsg)
		}
		return fmt.Errorf("query failed with unknown error")
	}
	
	fmt.Printf("✓ Query executed successfully (status: %s)\n", status)
	
	// If query has ID, test result retrieval
	if queryID, exists := queryResp["id"]; exists {
		return tc.TestQueryResult(queryID.(string))
	}
	
	return nil
}

// TestQueryResult tests the query result endpoint
func (tc *TestClient) TestQueryResult(queryID string) error {
	fmt.Printf("Testing GET /mcp/query/%s/result...\n", queryID)
	
	resp, err := tc.client.Get(tc.baseURL + "/mcp/query/" + queryID + "/result")
	if err != nil {
		return fmt.Errorf("failed to get query result: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w", err)
	}
	
	fmt.Printf("✓ Query result retrieved successfully\n")
	return nil
}

// TestHealth tests the health endpoint
func (tc *TestClient) TestHealth() error {
	fmt.Println("Testing GET /health...")
	
	resp, err := tc.client.Get(tc.baseURL + "/health")
	if err != nil {
		return fmt.Errorf("failed to get health: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	
	var health map[string]interface{}
	if err := json.Unmarshal(body, &health); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w", err)
	}
	
	status, ok := health["status"].(string)
	if !ok || status != "healthy" {
		return fmt.Errorf("service is not healthy: %v", health)
	}
	
	fmt.Printf("✓ Health check passed\n")
	return nil
}

// RunAllTests runs all tests
func (tc *TestClient) RunAllTests() error {
	fmt.Println("Running MCP Server tests...")
	fmt.Println(strings.Repeat("=", 50))
	
	tests := []struct {
		name string
		fn   func() error
	}{
		{"Health Check", tc.TestHealth},
		{"Schema", tc.TestSchema},
		{"Entities", tc.TestEntities},
		{"Query", tc.TestQuery},
	}
	
	for _, test := range tests {
		fmt.Printf("\n--- %s ---\n", test.name)
		if err := test.fn(); err != nil {
			fmt.Printf("✗ %s failed: %v\n", test.name, err)
			return err
		}
	}
	
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("All tests passed! ✓")
	return nil
}

func main() {
	baseURL := "http://localhost:8082"
	if len(os.Args) > 1 {
		baseURL = os.Args[1]
	}
	
	client := NewTestClient(baseURL)
	
	if err := client.RunAllTests(); err != nil {
		fmt.Printf("Tests failed: %v\n", err)
		os.Exit(1)
	}
}
