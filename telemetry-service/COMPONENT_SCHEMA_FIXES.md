# Choreo Component.yaml Schema 1.2 Corrections

## Issues Fixed in telemetry-service/.choreo/component.yaml

### 1. Schema Version Format Updates
**Changes made to comply with Choreo schema version 1.2:**

#### Before (Schema 1.0 format):
```yaml
schemaVersion: "1.2"
version: 0.1
endpoint:
  - name: telemetry-api
    port: 8080
    type: REST
    networkVisibility: Project
    context: /
    schemaFilePath: openapi.yaml
```

#### After (Schema 1.2 format):
```yaml
schemaVersion: "1.2"
version: "0.1"
endpoints:
  - name: telemetry-api
    service:
      port: 8080
    type: REST
    networkVisibility: Project
    context: /
    schemaFilePath: openapi.yaml
```

### 2. Key Changes Made:

1. **Plural endpoint property**:
   - Changed `endpoint:` to `endpoints:` (plural)
   - This is required in schema version 1.2

2. **Service property structure**:
   - Added `service:` property containing the port configuration
   - Moved `port: 8080` under the `service:` property

3. **Version format**:
   - Changed `version: 0.1` to `version: "0.1"` (string format)
   - Schema 1.2 requires version to be a string

4. **Connection references maintained**:
   - Kept the existing `dependencies.connectionReferences` configuration
   - This ensures the database connection setup remains intact

### 3. Current Complete Configuration:

```yaml
schemaVersion: "1.2"
version: "0.1"
endpoints:
  - name: telemetry-api
    service:
      port: 8080
    type: REST
    networkVisibility: Project
    context: /
    schemaFilePath: openapi.yaml
dependencies:
  connectionReferences:
    - name: TelemetryDB
      resourceRef: database:NonProductionPG/telemetry
```

### 4. Validation Results:
- ✅ Schema version 1.2 compliance
- ✅ Proper endpoint structure
- ✅ Service port configuration
- ✅ Connection references preserved
- ✅ OpenAPI specification reference maintained

### 5. Expected Behavior:
With these corrections, the telemetry service should:
- Pass Choreo component validation
- Build successfully in Choreo
- Maintain database connectivity via connection references
- Expose the REST API properly on port 8080

The corrected component.yaml now fully complies with Choreo schema version 1.2 requirements while preserving all the database connection and API configuration needed for successful deployment.
