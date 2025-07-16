/**
 * Application constants
 */

// API endpoints
export const API_ENDPOINTS = {
  AUTH: {
    LOGIN: '/auth/login',
    LOGOUT: '/auth/logout',
    REFRESH: '/auth/refresh',
    PROFILE: '/auth/profile'
  },
  DEVICES: {
    LIST: '/api/devices',
    DETAIL: '/api/devices/:id',
    METRICS: '/api/devices/:id/metrics',
    PROCESSES: '/api/devices/:id/processes',
    CONTAINERS: '/api/devices/:id/containers'
  },
  THREATS: {
    LIST: '/api/threats',
    DETAIL: '/api/threats/:id',
    DISMISS: '/api/threats/:id/dismiss'
  },
  HEALTH: {
    SYSTEM: '/api/health/system',
    SERVICES: '/api/health/services'
  },
  MCP: {
    QUERY: '/api/mcp/query',
    CHAT: '/api/mcp/chat'
  }
};

// Polling intervals (in milliseconds)
export const POLLING_INTERVALS = {
  FAST: 5000,    // 5 seconds
  MEDIUM: 30000, // 30 seconds
  SLOW: 60000,   // 1 minute
  VERY_SLOW: 300000 // 5 minutes
};

// Device types
export const DEVICE_TYPES = {
  LAPTOP: 'laptop',
  DESKTOP: 'desktop',
  SERVER: 'server',
  MOBILE: 'mobile',
  TABLET: 'tablet'
};

// Process states
export const PROCESS_STATES = {
  RUNNING: 'running',
  SLEEPING: 'sleeping',
  STOPPED: 'stopped',
  ZOMBIE: 'zombie'
};

// Container states
export const CONTAINER_STATES = {
  RUNNING: 'running',
  STOPPED: 'stopped',
  PAUSED: 'paused',
  RESTARTING: 'restarting',
  REMOVING: 'removing',
  DEAD: 'dead',
  CREATED: 'created'
};

// Threat severities
export const THREAT_SEVERITIES = {
  CRITICAL: 'critical',
  HIGH: 'high',
  MEDIUM: 'medium',
  LOW: 'low',
  INFO: 'info'
};

// System health statuses
export const HEALTH_STATUSES = {
  HEALTHY: 'healthy',
  WARNING: 'warning',
  CRITICAL: 'critical',
  UNKNOWN: 'unknown'
};

// User roles
export const USER_ROLES = {
  ADMIN: 'admin',
  ANALYST: 'analyst',
  VIEWER: 'viewer'
};

// Chart colors
export const CHART_COLORS = {
  PRIMARY: '#3B82F6',
  SUCCESS: '#10B981',
  WARNING: '#F59E0B',
  ERROR: '#EF4444',
  INFO: '#6366F1',
  SECONDARY: '#6B7280'
};

// Default pagination
export const PAGINATION = {
  DEFAULT_PAGE_SIZE: 20,
  PAGE_SIZE_OPTIONS: [10, 20, 50, 100]
};

// Local storage keys
export const STORAGE_KEYS = {
  AUTH_TOKEN: 'smartsec_auth_token',
  REFRESH_TOKEN: 'smartsec_refresh_token',
  USER_PREFERENCES: 'smartsec_user_preferences',
  SELECTED_DEVICE: 'smartsec_selected_device'
};

// Error messages
export const ERROR_MESSAGES = {
  NETWORK: 'Network error. Please check your connection.',
  UNAUTHORIZED: 'You are not authorized to perform this action.',
  FORBIDDEN: 'Access denied.',
  NOT_FOUND: 'Resource not found.',
  SERVER_ERROR: 'Server error. Please try again later.',
  VALIDATION: 'Please check your input and try again.',
  TIMEOUT: 'Request timed out. Please try again.'
};

// Success messages
export const SUCCESS_MESSAGES = {
  LOGIN: 'Successfully logged in!',
  LOGOUT: 'Successfully logged out!',
  SAVE: 'Changes saved successfully!',
  DELETE: 'Item deleted successfully!',
  UPDATE: 'Updated successfully!'
};

// Time formats
export const TIME_FORMATS = {
  DATETIME: 'MMM DD, YYYY HH:mm',
  DATE: 'MMM DD, YYYY',
  TIME: 'HH:mm',
  RELATIVE: 'relative'
};

// Refresh intervals for different data types
export const REFRESH_INTERVALS = {
  REAL_TIME: 1000,     // 1 second
  METRICS: 5000,       // 5 seconds
  DEVICES: 30000,      // 30 seconds
  THREATS: 60000,      // 1 minute
  HEALTH: 30000,       // 30 seconds
  LOGS: 10000          // 10 seconds
};
