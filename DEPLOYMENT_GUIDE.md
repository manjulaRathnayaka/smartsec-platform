# SmartSec Platform - Choreo Deployment Guide

This guide will help you deploy the SmartSec Platform to WSO2 Choreo.

## Prerequisites

1. **WSO2 Choreo Account**: Sign up at [console.choreo.dev](https://console.choreo.dev)
2. **GitHub Repository**: Your code must be in a GitHub repository
3. **Database Access**: PostgreSQL and Redis instances (can be created in Choreo)

## Step-by-Step Deployment

### 1. Prepare Your GitHub Repository

1. **Fork or clone** this repository to your GitHub account
2. **Update environment variables** in `.env.example` with your values
3. **Commit and push** all changes to your GitHub repository

### 2. Create a Choreo Project

1. Log in to [Choreo Console](https://console.choreo.dev)
2. Click **Create Project**
3. Enter project details:
   - **Name**: SmartSec Platform
   - **Description**: Comprehensive cybersecurity monitoring platform
   - **Visibility**: Private (or Public if desired)

### 3. Set Up Databases

#### PostgreSQL Database
1. In Choreo Console, go to **Dependencies** → **Databases**
2. Click **Add Database**
3. Select **PostgreSQL**
4. Configure:
   - **Name**: smartsec-postgres
   - **Version**: 15
   - **Storage**: 20GB (or as needed)
   - **Backup**: Enable
5. Note the connection details for later use

#### Redis Cache
1. In **Dependencies** → **Databases**
2. Click **Add Database**
3. Select **Redis**
4. Configure:
   - **Name**: smartsec-redis
   - **Version**: 7
   - **Memory**: 1GB (or as needed)
5. Note the connection details for later use

### 4. Deploy Components

Deploy components in the following order:

#### 4.1 Telemetry Service (Go Service)

1. Click **Create Component**
2. Select **Service** type
3. Connect your GitHub repository
4. Configure:
   - **Name**: telemetry-service
   - **Description**: Telemetry data collection service
   - **Build Path**: `/telemetry-service`
   - **Dockerfile**: Use existing Dockerfile
   - **Port**: 8080
5. Set environment variables:
   ```
   DATABASE_URL=<PostgreSQL connection string>
   REDIS_URL=<Redis connection string>
   PORT=8080
   ```
6. Deploy the component

#### 4.2 MCP Server (Node.js Service)

1. Click **Create Component**
2. Select **Service** type
3. Configure:
   - **Name**: mcp-server
   - **Description**: AI integration service
   - **Build Path**: `/mcp-server`
   - **Dockerfile**: Use existing Dockerfile
   - **Port**: 8081
4. Set environment variables:
   ```
   DATABASE_URL=<PostgreSQL connection string>
   TELEMETRY_SERVICE_URL=<Telemetry service URL>
   PORT=8081
   ```
5. Deploy the component

#### 4.3 BFF Service (Node.js Service)

1. Click **Create Component**
2. Select **Service** type
3. Configure:
   - **Name**: bff
   - **Description**: Backend for Frontend API
   - **Build Path**: `/bff`
   - **Dockerfile**: Use existing Dockerfile
   - **Port**: 3001
4. Set environment variables:
   ```
   DATABASE_URL=<PostgreSQL connection string>
   REDIS_URL=<Redis connection string>
   TELEMETRY_SERVICE_URL=<Telemetry service URL>
   MCP_SERVER_URL=<MCP server URL>
   JWT_SECRET=<your-jwt-secret>
   SESSION_SECRET=<your-session-secret>
   PORT=3001
   ```
5. Deploy the component

#### 4.4 Frontend (React Web App)

1. Click **Create Component**
2. Select **Web App** type
3. Configure:
   - **Name**: frontend
   - **Description**: React web portal
   - **Build Path**: `/frontend`
   - **Dockerfile**: Use existing Dockerfile
   - **Port**: 3000
4. Set environment variables:
   ```
   REACT_APP_API_BASE_URL=<BFF service URL>
   REACT_APP_TELEMETRY_URL=<Telemetry service URL>
   REACT_APP_MCP_URL=<MCP server URL>
   BFF_SERVICE_URL=<BFF service internal URL>
   ```
5. Deploy the component

#### 4.5 Laptop Agent (Scheduled Task)

1. Click **Create Component**
2. Select **Scheduled Task** type
3. Configure:
   - **Name**: laptop-agent
   - **Description**: Device monitoring agent
   - **Build Path**: `/laptop-agent`
   - **Dockerfile**: Use existing Dockerfile
   - **Schedule**: Every 30 seconds (or as needed)
4. Set environment variables:
   ```
   TELEMETRY_SERVICE_URL=<Telemetry service URL>
   AGENT_POLL_INTERVAL=30s
   ```
5. Deploy the component

### 5. Configure Networking

1. In Choreo Console, go to **Networking** → **Service Mesh**
2. Configure service-to-service communication:
   - **BFF** → **Telemetry Service**
   - **BFF** → **MCP Server**
   - **MCP Server** → **Telemetry Service**
   - **Laptop Agent** → **Telemetry Service**
3. Set up load balancing and health checks

### 6. Set Up Monitoring

1. Enable **Observability** in Choreo Console
2. Configure:
   - **Logging**: Enable structured logging
   - **Metrics**: Enable Prometheus metrics
   - **Tracing**: Enable distributed tracing
   - **Alerts**: Set up alerts for key metrics

### 7. Configure Security

1. **Authentication**: Set up OAuth2 integration
2. **Authorization**: Configure JWT settings
3. **Network Security**: Set up VPC and firewall rules
4. **Secrets Management**: Use Choreo's secret management

### 8. Database Initialization

1. Connect to your PostgreSQL database
2. Run the initialization script:
   ```bash
   psql -h <postgres-host> -U <username> -d <database> -f database/schema.sql
   ```

### 9. Test the Deployment

1. **Access the frontend** at the provided URL
2. **Test API endpoints** using the BFF service URL
3. **Verify data collection** from the telemetry service
4. **Check logs** for any errors

### 10. Production Optimization

1. **Scaling**: Configure auto-scaling for services
2. **Caching**: Optimize Redis usage
3. **CDN**: Set up CDN for static assets
4. **Performance**: Monitor and optimize performance

## Environment Variables Reference

### Database
- `DATABASE_URL`: PostgreSQL connection string
- `REDIS_URL`: Redis connection string

### Authentication
- `JWT_SECRET`: Secret key for JWT tokens
- `SESSION_SECRET`: Secret key for sessions
- `OAUTH_CLIENT_ID`: OAuth client ID
- `OAUTH_CLIENT_SECRET`: OAuth client secret

### Service URLs
- `BFF_BASE_URL`: BFF service URL
- `TELEMETRY_SERVICE_URL`: Telemetry service URL
- `MCP_SERVER_URL`: MCP server URL

### Frontend
- `REACT_APP_API_BASE_URL`: BFF service URL for frontend
- `REACT_APP_TELEMETRY_URL`: Telemetry service URL for frontend
- `REACT_APP_MCP_URL`: MCP server URL for frontend

## Troubleshooting

### Common Issues

1. **Database Connection Errors**
   - Check connection strings
   - Verify database is running
   - Check firewall rules

2. **Service Communication Issues**
   - Verify service mesh configuration
   - Check service discovery
   - Verify port configurations

3. **Frontend Not Loading**
   - Check environment variables
   - Verify API endpoints
   - Check CORS configuration

4. **Authentication Issues**
   - Verify JWT secret
   - Check OAuth configuration
   - Verify session settings

### Debugging Commands

```bash
# Check service health
curl -f <service-url>/health

# View logs
choreo logs --component=<component-name>

# Check database connection
psql -h <host> -U <user> -d <database> -c "SELECT 1"

# Test Redis connection
redis-cli -h <host> -p <port> ping
```

## Support

For deployment issues:
1. Check [Choreo Documentation](https://wso2.com/choreo/docs/)
2. Review component logs in Choreo Console
3. Check GitHub repository issues
4. Contact support at support@smartsec.com

## Next Steps

After successful deployment:
1. Set up monitoring dashboards
2. Configure alerting rules
3. Set up CI/CD pipelines
4. Plan for scaling and backup strategies
5. Implement additional security measures
