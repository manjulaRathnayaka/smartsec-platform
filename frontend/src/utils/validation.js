/**
 * Validation utilities
 */

/**
 * Validate email format
 * @param {string} email - Email to validate
 * @returns {boolean} True if valid
 */
export const validateEmail = (email) => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
};

/**
 * Validate password strength
 * @param {string} password - Password to validate
 * @returns {object} Validation result with strength and messages
 */
export const validatePassword = (password) => {
  const result = {
    isValid: false,
    strength: 'weak',
    messages: []
  };

  if (!password) {
    result.messages.push('Password is required');
    return result;
  }

  if (password.length < 8) {
    result.messages.push('Password must be at least 8 characters long');
  }

  if (!/[A-Z]/.test(password)) {
    result.messages.push('Password must contain at least one uppercase letter');
  }

  if (!/[a-z]/.test(password)) {
    result.messages.push('Password must contain at least one lowercase letter');
  }

  if (!/\d/.test(password)) {
    result.messages.push('Password must contain at least one number');
  }

  if (!/[!@#$%^&*(),.?":{}|<>]/.test(password)) {
    result.messages.push('Password must contain at least one special character');
  }

  // Determine strength
  if (result.messages.length === 0) {
    result.isValid = true;
    result.strength = 'strong';
  } else if (result.messages.length <= 2) {
    result.strength = 'medium';
  }

  return result;
};

/**
 * Validate IP address
 * @param {string} ip - IP address to validate
 * @returns {boolean} True if valid
 */
export const validateIP = (ip) => {
  const ipRegex = /^(\d{1,3}\.){3}\d{1,3}$/;
  if (!ipRegex.test(ip)) return false;

  const parts = ip.split('.');
  return parts.every(part => {
    const num = parseInt(part, 10);
    return num >= 0 && num <= 255;
  });
};

/**
 * Validate URL format
 * @param {string} url - URL to validate
 * @returns {boolean} True if valid
 */
export const validateURL = (url) => {
  try {
    new URL(url);
    return true;
  } catch {
    return false;
  }
};

/**
 * Validate port number
 * @param {number|string} port - Port number to validate
 * @returns {boolean} True if valid
 */
export const validatePort = (port) => {
  const portNum = parseInt(port, 10);
  return !isNaN(portNum) && portNum >= 1 && portNum <= 65535;
};

/**
 * Validate required fields
 * @param {object} data - Object with form data
 * @param {string[]} requiredFields - Array of required field names
 * @returns {object} Validation result
 */
export const validateRequired = (data, requiredFields) => {
  const errors = {};

  requiredFields.forEach(field => {
    if (!data[field] || (typeof data[field] === 'string' && data[field].trim() === '')) {
      errors[field] = `${field.charAt(0).toUpperCase() + field.slice(1)} is required`;
    }
  });

  return {
    isValid: Object.keys(errors).length === 0,
    errors
  };
};
