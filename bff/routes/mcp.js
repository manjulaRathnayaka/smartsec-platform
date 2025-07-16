const express = require('express');
const axios = require('axios');
const winston = require('winston');
const { body, validationResult } = require('express-validator');

const router = express.Router();

const logger = winston.createLogger({
  level: process.env.LOG_LEVEL || 'info',
  format: winston.format.simple(),
  transports: [new winston.transports.Console()]
});

const MCP_SERVER_URL = process.env.MCP_SERVER_URL || 'http://localhost:8082';

// Natural Language Query endpoint
router.post('/query', [
  body('query').isString().isLength({ min: 1, max: 1000 }).trim()
], async (req, res) => {
  try {
    const errors = validationResult(req);
    if (!errors.isEmpty()) {
      return res.status(400).json({ errors: errors.array() });
    }

    const { query } = req.body;

    logger.info(`Natural language query from user ${req.user.email}: ${query}`);

    // Send query to MCP server
    const response = await axios.post(`${MCP_SERVER_URL}/query`, {
      query,
      user_id: req.user.id,
      user_role: req.user.role
    }, {
      timeout: 30000
    });

    res.json({
      query,
      result: response.data,
      timestamp: new Date().toISOString()
    });
  } catch (error) {
    logger.error('MCP query error:', error);

    if (error.response) {
      res.status(error.response.status).json({
        error: 'MCP server error',
        message: error.response.data?.error || 'Unknown error'
      });
    } else if (error.code === 'ECONNREFUSED') {
      res.status(503).json({
        error: 'Service unavailable',
        message: 'MCP server is not available'
      });
    } else {
      res.status(500).json({
        error: 'Internal server error',
        message: 'Failed to process query'
      });
    }
  }
});

// Get available tools/capabilities
router.get('/tools', async (req, res) => {
  try {
    const response = await axios.get(`${MCP_SERVER_URL}/tools`, {
      timeout: 10000
    });

    res.json(response.data);
  } catch (error) {
    logger.error('Error fetching MCP tools:', error);
    res.status(500).json({
      error: 'Failed to fetch tools',
      message: error.message
    });
  }
});

// Execute specific tool
router.post('/tools/:toolName', [
  body('arguments').isObject().optional()
], async (req, res) => {
  try {
    const errors = validationResult(req);
    if (!errors.isEmpty()) {
      return res.status(400).json({ errors: errors.array() });
    }

    const { toolName } = req.params;
    const { arguments: toolArgs = {} } = req.body;

    logger.info(`Tool execution: ${toolName} by user ${req.user.email}`);

    const response = await axios.post(`${MCP_SERVER_URL}/tools/${toolName}`, {
      arguments: toolArgs,
      user_id: req.user.id,
      user_role: req.user.role
    }, {
      timeout: 30000
    });

    res.json({
      tool: toolName,
      result: response.data,
      timestamp: new Date().toISOString()
    });
  } catch (error) {
    logger.error(`Tool execution error (${req.params.toolName}):`, error);

    if (error.response) {
      res.status(error.response.status).json({
        error: 'Tool execution failed',
        message: error.response.data?.error || 'Unknown error'
      });
    } else {
      res.status(500).json({
        error: 'Internal server error',
        message: 'Failed to execute tool'
      });
    }
  }
});

// Get query history for current user
router.get('/history', async (req, res) => {
  try {
    const { limit = 50, offset = 0 } = req.query;

    const response = await axios.get(`${MCP_SERVER_URL}/history`, {
      params: {
        user_id: req.user.id,
        limit,
        offset
      },
      timeout: 10000
    });

    res.json(response.data);
  } catch (error) {
    logger.error('Error fetching query history:', error);
    res.status(500).json({
      error: 'Failed to fetch history',
      message: error.message
    });
  }
});

module.exports = router;
