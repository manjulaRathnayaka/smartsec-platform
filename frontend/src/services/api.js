import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:3001';

const api = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true,
});

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('authToken');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor to handle auth errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('authToken');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export const authAPI = {
  login: (credentials) => api.post('/auth/login', credentials),
  loginSSO: (provider) => api.get(`/auth/sso/${provider}`),
  logout: () => api.post('/auth/logout'),
  getProfile: () => api.get('/auth/profile'),
  refreshToken: () => api.post('/auth/refresh'),
};

export const telemetryAPI = {
  getDevices: () => api.get('/telemetry/devices'),
  getDeviceActivity: (deviceId) => api.get(`/telemetry/devices/${deviceId}/activity`),
  getContainers: (deviceId) => api.get(`/telemetry/devices/${deviceId}/containers`),
  getThreats: (deviceId) => api.get(`/telemetry/devices/${deviceId}/threats`),
  getSystemHealth: (deviceId) => api.get(`/telemetry/devices/${deviceId}/health`),
  getFleetActivity: () => api.get('/telemetry/fleet/activity'),
  getFleetThreats: () => api.get('/telemetry/fleet/threats'),
  getFleetHealth: () => api.get('/telemetry/fleet/health'),
  getDashboardData: (deviceId, timeRange) => api.get('/api/dashboard', { params: { device_id: deviceId, time_range: timeRange } }),
};

export const mcpAPI = {
  submitQuery: (query) => api.post('/mcp/query', { query }),
  getQueryHistory: () => api.get('/mcp/history'),
  executeTool: (toolName, params) => api.post('/mcp/execute', { toolName, params }),
};

export const dashboardAPI = {
  getDashboard: () => api.get('/api/dashboard'),
  getAdminOverview: () => api.get('/api/admin/overview'),
};

export default api;
