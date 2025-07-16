/**
 * Utility functions for formatting data
 */

/**
 * Format bytes to human readable format
 * @param {number} bytes - The number of bytes
 * @param {number} decimals - Number of decimal places
 * @returns {string} Formatted string
 */
export const formatBytes = (bytes, decimals = 2) => {
  if (bytes === 0) return '0 Bytes';

  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
};

/**
 * Format timestamp to relative time
 * @param {string|Date} timestamp - The timestamp to format
 * @returns {string} Relative time string
 */
export const formatRelativeTime = (timestamp) => {
  const now = new Date();
  const time = new Date(timestamp);
  const diff = now - time;

  const minutes = Math.floor(diff / 60000);
  const hours = Math.floor(diff / 3600000);
  const days = Math.floor(diff / 86400000);

  if (minutes < 1) return 'just now';
  if (minutes < 60) return `${minutes}m ago`;
  if (hours < 24) return `${hours}h ago`;
  if (days < 30) return `${days}d ago`;

  return time.toLocaleDateString();
};

/**
 * Format CPU percentage
 * @param {number} cpu - CPU usage as decimal (0-1)
 * @returns {string} Formatted percentage
 */
export const formatCPU = (cpu) => {
  return `${(cpu * 100).toFixed(1)}%`;
};

/**
 * Format memory usage
 * @param {number} used - Used memory in bytes
 * @param {number} total - Total memory in bytes
 * @returns {string} Formatted memory usage
 */
export const formatMemory = (used, total) => {
  const usedFormatted = formatBytes(used);
  const totalFormatted = formatBytes(total);
  const percentage = ((used / total) * 100).toFixed(1);

  return `${usedFormatted} / ${totalFormatted} (${percentage}%)`;
};

/**
 * Format threat severity
 * @param {string} severity - Threat severity level
 * @returns {object} Color classes for the severity
 */
export const formatThreatSeverity = (severity) => {
  const severityMap = {
    critical: {
      bg: 'bg-red-100',
      text: 'text-red-800',
      border: 'border-red-200'
    },
    high: {
      bg: 'bg-orange-100',
      text: 'text-orange-800',
      border: 'border-orange-200'
    },
    medium: {
      bg: 'bg-yellow-100',
      text: 'text-yellow-800',
      border: 'border-yellow-200'
    },
    low: {
      bg: 'bg-green-100',
      text: 'text-green-800',
      border: 'border-green-200'
    },
    info: {
      bg: 'bg-blue-100',
      text: 'text-blue-800',
      border: 'border-blue-200'
    }
  };

  return severityMap[severity.toLowerCase()] || severityMap.info;
};

/**
 * Format status badge
 * @param {string} status - Status string
 * @returns {object} Color classes for the status
 */
export const formatStatus = (status) => {
  const statusMap = {
    active: {
      bg: 'bg-green-100',
      text: 'text-green-800',
      border: 'border-green-200'
    },
    inactive: {
      bg: 'bg-gray-100',
      text: 'text-gray-800',
      border: 'border-gray-200'
    },
    error: {
      bg: 'bg-red-100',
      text: 'text-red-800',
      border: 'border-red-200'
    },
    warning: {
      bg: 'bg-yellow-100',
      text: 'text-yellow-800',
      border: 'border-yellow-200'
    },
    pending: {
      bg: 'bg-blue-100',
      text: 'text-blue-800',
      border: 'border-blue-200'
    }
  };

  return statusMap[status.toLowerCase()] || statusMap.pending;
};

/**
 * Truncate text with ellipsis
 * @param {string} text - Text to truncate
 * @param {number} maxLength - Maximum length
 * @returns {string} Truncated text
 */
export const truncateText = (text, maxLength = 50) => {
  if (!text || text.length <= maxLength) return text;
  return text.substring(0, maxLength) + '...';
};
