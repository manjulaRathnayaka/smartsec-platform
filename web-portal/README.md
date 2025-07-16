# SmartSec Platform Web Portal

A React-based web portal for the SmartSec cybersecurity platform with natural language query interface.

## Architecture

- **Frontend**: React 18 with Tailwind CSS
- **BFF (Backend for Frontend)**: Node.js/Express API
- **Authentication**: JWT + OAuth2/OpenID Connect
- **State Management**: React Query + Context API
- **Styling**: Tailwind CSS with Heroicons

## Features

### User Features
- 🔐 SSO Login (OAuth2/OpenID Connect)
- 📊 Device Activity Dashboard
- 🔍 Process and Container Monitoring
- ⚠️ Threat Detection and Findings
- 💚 System Health Monitoring
- 🤖 AI-Powered Natural Language Queries

### Admin Features
- 👑 Fleet-wide Activity Overview
- 📈 Cross-device Analytics
- 🚨 Organization-wide Threat Management
- 📋 User and Device Management

### AI Assistant
- 🗣️ Natural Language Query Interface
- 🔗 MCP Server Integration
- 📝 Query History and Results
- 💡 Suggested Security Queries

## Quick Start

### Prerequisites
- Node.js 18+
- PostgreSQL database
- Redis (for sessions)
- MCP Server running on port 8082
- Telemetry Service running on port 8080

### 1. Install Dependencies

```bash
# Install BFF dependencies
cd bff
npm install

# Install frontend dependencies
cd ../frontend
npm install
```

### 2. Configuration

```bash
# Copy environment files
cp bff/.env.example bff/.env
cp frontend/.env.example frontend/.env

# Edit the .env files with your configuration
```

### 3. Start Development Servers

```bash
# Terminal 1: Start BFF server
cd bff
npm run dev

# Terminal 2: Start frontend
cd frontend
npm start
```

The application will be available at:
- Frontend: http://localhost:3000
- BFF API: http://localhost:3001

## API Endpoints

### Authentication
- `POST /auth/login` - Local login
- `GET /auth/oauth2` - OAuth2 login
- `GET /auth/callback` - OAuth2 callback
- `POST /auth/logout` - Logout
- `GET /auth/profile` - Get user profile

### Telemetry
- `GET /telemetry/devices` - Get user devices
- `GET /telemetry/devices/:id/activity` - Get device processes
- `GET /telemetry/devices/:id/containers` - Get device containers
- `GET /telemetry/devices/:id/threats` - Get device threats
- `GET /telemetry/devices/:id/health` - Get device health

### MCP (AI Assistant)
- `POST /mcp/query` - Submit natural language query
- `GET /mcp/history` - Get query history
- `POST /mcp/execute` - Execute MCP tool

### Dashboard
- `GET /api/dashboard` - Get dashboard data
- `GET /api/admin/overview` - Get admin overview (admin only)

## Demo Credentials

```
Email: admin@smartsec.com
Password: password123
Role: admin

Email: user@smartsec.com
Password: password123
Role: user
```

## Environment Variables

### BFF (.env)
```
PORT=3001
NODE_ENV=development
CORS_ORIGIN=http://localhost:3000
DATABASE_URL=postgres://...
JWT_SECRET=your-jwt-secret
MCP_SERVER_URL=http://localhost:8082
TELEMETRY_API_URL=http://localhost:8080
```

### Frontend (.env)
```
REACT_APP_API_URL=http://localhost:3001
REACT_APP_APP_NAME=SmartSec Platform
```

## Components Structure

```
src/
├── components/
│   ├── auth/           # Login, OAuth callback
│   ├── dashboard/      # Main dashboard
│   ├── devices/        # Device activity, processes, containers
│   ├── threats/        # Threat findings and management
│   ├── health/         # System health monitoring
│   ├── ai/            # AI assistant and natural language queries
│   ├── admin/         # Admin dashboard and fleet management
│   ├── layout/        # Header, sidebar, main layout
│   └── common/        # Shared components
├── contexts/          # React contexts (Auth, Device)
├── services/          # API services
└── hooks/            # Custom React hooks
```

## Key Features

### Natural Language Query Interface
The AI Assistant allows users to query their security data using natural language:

- "Show all processes that ran with root privileges in the last 24 hours"
- "Which containers were started by non-admin users?"
- "List the most common commands run by marketing laptops"
- "Show me any suspicious network activity"

### Real-time Monitoring
- Process activity monitoring
- Container lifecycle tracking
- Threat detection alerts
- System health metrics

### Role-based Access
- Regular users see only their own devices
- Admins can view fleet-wide data
- Granular permission controls

## Development

### Adding New Components
1. Create component in appropriate directory
2. Add to router in `App.js`
3. Update navigation in `Sidebar.js`
4. Add API endpoints in `services/api.js`

### Styling Guidelines
- Use Tailwind CSS utility classes
- Follow responsive design principles
- Use Heroicons for consistency
- Implement dark mode support (future)

## Deployment

### Production Build
```bash
# Build frontend
cd frontend
npm run build

# The build files will be in frontend/build/
```

### Docker Support
```bash
# Build and run with Docker
docker-compose up -d
```

## Integration

### MCP Server Integration
The portal integrates with the MCP server for natural language queries:
- Sends user queries to `/mcp/query`
- Receives structured responses
- Displays results in chat interface
- Maintains query history

### Telemetry Service Integration
Real-time data from the telemetry service:
- Device registration and status
- Process and container monitoring
- Threat detection results
- System health metrics

## Security

- JWT-based authentication
- OAuth2/OpenID Connect support
- Session management with Redis
- CORS protection
- Input validation and sanitization
- Role-based authorization

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License - see LICENSE file for details
