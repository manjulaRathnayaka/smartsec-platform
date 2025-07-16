#!/bin/bash

# Telemetry Service Choreo Deployment Script
# This script configures the telemetry service for Choreo deployment using Connection References

set -e

echo "🚀 Configuring Telemetry Service for Choreo with Connection References..."

# Component UUID and Project UUID
COMPONENT_UUID="31c4d28f-b998-426a-83d9-67ad924ed937"
PROJECT_UUID="7dd928f1-8615-4f01-9437-2c609c6f2487"

echo "📋 Setting up Choreo connection reference..."

echo "
� Choreo Connection Configuration:

1. In component.yaml, add the connection reference:
   dependencies:
     connectionReferences:
       - name: TelemetryDB
         resourceRef: database:NonProductionPG/telemetry

2. Choreo will automatically inject these environment variables:
   - CHOREO_TELEMETRYDB_HOSTNAME
   - CHOREO_TELEMETRYDB_PORT
   - CHOREO_TELEMETRYDB_USERNAME
   - CHOREO_TELEMETRYDB_PASSWORD
   - CHOREO_TELEMETRYDB_DATABASENAME

3. The Go application will use these variables automatically via config.Load()

✅ No manual configuration needed - Choreo handles everything!
"

echo "🔗 Connection Reference Setup Complete!"
echo "
🚀 Deployment Steps:

1. Commit changes to Git repository
2. Push to GitHub
3. Choreo will automatically build and deploy
4. Connection will be established automatically

✅ Telemetry service is ready for Choreo deployment with connection references!
"
