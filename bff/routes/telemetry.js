const express = require('express');
const axios = require('axios');
const winston = require('winston');
const { query, validationResult } = require('express-validator');
const { requireAdmin } = require('../middleware/auth');

const router = express.Router();

const logger = winston.createLogger({
  level: process.env.LOG_LEVEL || 'info',
  format: winston.format.simple(),
  transports: [new winston.transports.Console()]
});

const TELEMETRY_API_URL = process.env.TELEMETRY_API_URL || 'http://localhost:8080';

// Get device activity for current user
router.get('/devices', [
  query('limit').optional().isInt({ min: 1, max: 100 }),
  query('offset').optional().isInt({ min: 0 }),
  query('device_id').optional().isString()
], async (req, res) => {
  try {
    const errors = validationResult(req);
    if (!errors.isEmpty()) {
      return res.status(400).json({ errors: errors.array() });
    }

    const { limit = 20, offset = 0, device_id } = req.query;

    // For regular users, filter by their associated devices
    // For admins, show all devices
    let params = { limit, offset };

    if (req.user.role !== 'admin') {
      params.user_id = req.user.id;
    }

    if (device_id) {
      params.device_id = device_id;
    }

    const response = await axios.get(`${TELEMETRY_API_URL}/devices`, {
      params,
      timeout: 10000
    });

    res.json(response.data);
  } catch (error) {
    logger.error('Error fetching devices:', error);

    if (error.response) {
      res.status(error.response.status).json({
        error: 'Telemetry service error',
        message: error.response.data?.error || 'Unknown error'
      });
    } else {
      res.status(500).json({
        error: 'Internal server error',
        message: 'Failed to fetch devices'
      });
    }
  }
});

// Get containers for a specific device
router.get('/containers', [
  query('device_id').isString(),
  query('limit').optional().isInt({ min: 1, max: 100 }),
  query('offset').optional().isInt({ min: 0 })
], async (req, res) => {
  try {
    const errors = validationResult(req);
    if (!errors.isEmpty()) {
      return res.status(400).json({ errors: errors.array() });
    }

    const { device_id, limit = 20, offset = 0 } = req.query;

    const response = await axios.get(`${TELEMETRY_API_URL}/containers`, {
      params: { device_id, limit, offset },
      timeout: 10000
    });

    res.json(response.data);
  } catch (error) {
    logger.error('Error fetching containers:', error);

    if (error.response) {
      res.status(error.response.status).json({
        error: 'Telemetry service error',
        message: error.response.data?.error || 'Unknown error'
      });
    } else {
      res.status(500).json({
        error: 'Internal server error',
        message: 'Failed to fetch containers'
      });
    }
  }
});

// Get system health metrics
router.get('/health', async (req, res) => {
  try {
    const response = await axios.get(`${TELEMETRY_API_URL}/health`, {
      timeout: 10000
    });

    res.json(response.data);
  } catch (error) {
    logger.error('Error fetching system health:', error);
    res.status(500).json({
      error: 'Failed to fetch system health',
      message: error.message
    });
  }
});

// Get device statistics (admin only)
router.get('/stats', requireAdmin, async (req, res) => {
  try {
    const response = await axios.get(`${TELEMETRY_API_URL}/stats`, {
      timeout: 10000
    });

    res.json(response.data);
  } catch (error) {
    logger.error('Error fetching statistics:', error);
    res.status(500).json({
      error: 'Failed to fetch statistics',
      message: error.message
    });
  }
});

// Get recent activities
router.get('/activities', [
  query('limit').optional().isInt({ min: 1, max: 100 }),
  query('offset').optional().isInt({ min: 0 }),
  query('device_id').optional().isString(),
  query('type').optional().isIn(['process', 'container', 'network'])
], async (req, res) => {
  try {
    const errors = validationResult(req);
    if (!errors.isEmpty()) {
      return res.status(400).json({ errors: errors.array() });
    }

    const { limit = 20, offset = 0, device_id, type } = req.query;

    let params = { limit, offset };

    if (req.user.role !== 'admin') {
      params.user_id = req.user.id;
    }

    if (device_id) params.device_id = device_id;
    if (type) params.type = type;

    const response = await axios.get(`${TELEMETRY_API_URL}/activities`, {
      params,
      timeout: 10000
    });

    res.json(response.data);
  } catch (error) {
    logger.error('Error fetching activities:', error);

    if (error.response) {
      res.status(error.response.status).json({
        error: 'Telemetry service error',
        message: error.response.data?.error || 'Unknown error'
      });
    } else {
      res.status(500).json({
        error: 'Internal server error',
        message: 'Failed to fetch activities'
      });
    }
  }
});

// Get threat findings
router.get('/threats', [
  query('limit').optional().isInt({ min: 1, max: 100 }),
  query('offset').optional().isInt({ min: 0 }),
  query('severity').optional().isIn(['low', 'medium', 'high', 'critical']),
  query('device_id').optional().isString()
], async (req, res) => {
  try {
    const errors = validationResult(req);
    if (!errors.isEmpty()) {
      return res.status(400).json({ errors: errors.array() });
    }

    const { limit = 20, offset = 0, severity, device_id } = req.query;

    let params = { limit, offset };

    if (req.user.role !== 'admin') {
      params.user_id = req.user.id;
    }

    if (severity) params.severity = severity;
    if (device_id) params.device_id = device_id;

    const response = await axios.get(`${TELEMETRY_API_URL}/threats`, {
      params,
      timeout: 10000
    });

    res.json(response.data);
  } catch (error) {
    logger.error('Error fetching threats:', error);

    if (error.response) {
      res.status(error.response.status).json({
        error: 'Telemetry service error',
        message: error.response.data?.error || 'Unknown error'
      });
    } else {
      res.status(500).json({
        error: 'Internal server error',
        message: 'Failed to fetch threats'
      });
    }
  }
});

module.exports = router;
