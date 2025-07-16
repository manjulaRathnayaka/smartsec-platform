# SmartSec Platform - Choreo Deployment Ready

A comprehensive cybersecurity monitoring platform built with modern technologies and deployed on WSO2 Choreo.

## 🏗️ Architecture Overview

The SmartSec Platform consists of multiple microservices designed for cloud-native deployment:

### Components

1. **Frontend (React SPA)** - Web portal for monitoring and management
2. **BFF (Backend for Frontend)** - API gateway and authentication layer
3. **Telemetry Service** - Data collection and processing
4. **MCP Server** - Model Context Protocol server for AI integration
5. **Laptop Agent** - Device monitoring agent

## 🚀 Choreo Deployment

This repository is structured for deployment on WSO2 Choreo with the following components:

### Service Components
- `frontend/` - React web application (WebApp component)
- `bff/` - Backend for Frontend API (Service component)
- `telemetry-service/` - Telemetry data service (Service component)
- `mcp-server/` - AI integration service (Service component)

### Task Components
- `laptop-agent/` - Device monitoring agent (Scheduled Task component)

## 📦 Quick Start

### Prerequisites
- Node.js 18+
- Go 1.21+
- PostgreSQL 15+
- Redis 7+

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd smartsec-platform
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start all services**
   ```bash
   # Start databases (Docker)
   docker-compose up -d postgres redis

   # Start backend services
   cd telemetry-service && npm start &
   cd mcp-server && npm start &
   cd bff && npm start &

   # Start frontend
   cd frontend && npm start
   ```

4. **Access the application**
   - Frontend: http://localhost:3000
   - BFF API: http://localhost:3001
   - Telemetry Service: http://localhost:8080
   - MCP Server: http://localhost:8081

## 🔧 Configuration

### Environment Variables

Each component has its own `.env` file. Key variables include:

```env
# Database
DATABASE_URL=postgresql://user:password@localhost:5432/smartsec
REDIS_URL=redis://localhost:6379

# Authentication
JWT_SECRET=your-jwt-secret
OAUTH_CLIENT_ID=your-oauth-client-id
OAUTH_CLIENT_SECRET=your-oauth-client-secret

# API Endpoints
BFF_BASE_URL=http://localhost:3001
TELEMETRY_SERVICE_URL=http://localhost:8080
MCP_SERVER_URL=http://localhost:8081
```

## 🏛️ Database Schema

The platform uses PostgreSQL with the following main tables:

- `devices` - Device information and status
- `telemetry_data` - Time-series monitoring data
- `threats` - Security findings and alerts
- `users` - User authentication and profiles
- `organizations` - Multi-tenant support

## 🔐 Security Features

- JWT-based authentication
- OAuth2 integration
- Role-based access control (RBAC)
- API rate limiting
- CORS protection
- Input validation and sanitization
- SQL injection prevention

## 📊 Monitoring & Observability

- Health check endpoints
- Structured logging
- Metrics collection
- Distributed tracing ready
- Performance monitoring

## 🧪 Testing

```bash
# Run all tests
npm run test:all

# Run integration tests
npm run test:integration

# Run e2e tests
npm run test:e2e
```

## 🚀 Deployment

### Choreo Deployment

1. **Connect your GitHub repository** to Choreo
2. **Create components** for each service:
   - Frontend: WebApp component
   - BFF: Service component
   - Telemetry Service: Service component
   - MCP Server: Service component
   - Laptop Agent: Scheduled Task component

3. **Configure databases**:
   - PostgreSQL for persistent data
   - Redis for caching and sessions

4. **Set environment variables** in Choreo console

5. **Deploy components** in the correct order:
   1. Databases
   2. Backend services
   3. BFF
   4. Frontend
   5. Laptop Agent

### Manual Deployment

```bash
# Build all components
npm run build:all

# Deploy to your infrastructure
npm run deploy:production
```

## 🔗 API Documentation

- BFF API: `/api/docs`
- Telemetry Service: `/docs`
- MCP Server: `/api/docs`

## 📈 Scaling

The platform is designed for horizontal scaling:

- **Stateless services** for easy scaling
- **Database connection pooling**
- **Redis clustering** support
- **Load balancer** ready

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

For support and questions:
- GitHub Issues
- Documentation: `/docs`
- Email: support@smartsec.com

## 🗂️ Project Structure

```
smartsec-platform/
├── frontend/                 # React web application
│   ├── src/
│   ├── public/
│   ├── package.json
│   └── .choreo/
├── bff/                      # Backend for Frontend
│   ├── routes/
│   ├── middleware/
│   ├── server.js
│   └── .choreo/
├── telemetry-service/        # Telemetry data service
│   ├── internal/
│   ├── main.go
│   └── .choreo/
├── mcp-server/               # AI integration service
│   ├── src/
│   ├── package.json
│   └── .choreo/
├── laptop-agent/             # Device monitoring agent
│   ├── main.go
│   └── .choreo/
├── database/                 # Database schemas and migrations
├── docker-compose.yml        # Local development setup
├── .env.example             # Environment variables template
└── README.md                # This file
```
