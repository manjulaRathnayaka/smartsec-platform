package main

import (
	"fmt"
	"os"
)

// Example of how to use Choreo connection configuration
// Based on the Choreo UI connection configuration shown
func main() {
	// Import "os" package to get environment variables
	hostname := os.Getenv("CHOREO_TELEMETRYDB_HOSTNAME")
	port := os.Getenv("CHOREO_TELEMETRYDB_PORT")
	username := os.Getenv("CHOREO_TELEMETRYDB_USERNAME")
	password := os.Getenv("CHOREO_TELEMETRYDB_PASSWORD")
	databasename := os.Getenv("CHOREO_TELEMETRYDB_DATABASENAME")

	// Use the environment variables to connect to the database
	fmt.Printf("Connecting to database:\n")
	fmt.Printf("Host: %s\n", hostname)
	fmt.Printf("Port: %s\n", port)
	fmt.Printf("Database: %s\n", databasename)
	fmt.Printf("Username: %s\n", username)

	// Build connection string
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=require",
		username,
		password,
		hostname,
		port,
		databasename,
	)

	fmt.Printf("Connection String: %s\n", connectionString)
}
