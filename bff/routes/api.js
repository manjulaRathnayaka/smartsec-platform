const express = require('express');
const winston = require('winston');
const { requireAdmin } = require('../middleware/auth');

const router = express.Router();

const logger = winston.createLogger({
  level: process.env.LOG_LEVEL || 'info',
  format: winston.format.simple(),
  transports: [new winston.transports.Console()]
});

// Dashboard summary endpoint
router.get('/dashboard', async (req, res) => {
  try {
    // This would typically aggregate data from multiple sources
    const dashboardData = {
      user: {
        id: req.user.id,
        name: req.user.name,
        role: req.user.role,
        department: req.user.department
      },
      summary: {
        totalDevices: req.user.role === 'admin' ? 25 : 3,
        activeContainers: req.user.role === 'admin' ? 142 : 8,
        threatLevel: 'medium',
        systemHealth: 'good'
      },
      recentActivities: [
        {
          id: 1,
          type: 'process',
          description: 'New process started: nginx',
          timestamp: new Date().toISOString(),
          severity: 'info'
        },
        {
          id: 2,
          type: 'container',
          description: 'Container deployment: web-app-v2',
          timestamp: new Date().toISOString(),
          severity: 'info'
        }
      ],
      threats: [
        {
          id: 1,
          title: 'Suspicious network activity',
          severity: 'medium',
          device: 'web-server-01',
          timestamp: new Date().toISOString()
        }
      ]
    };

    res.json(dashboardData);
  } catch (error) {
    logger.error('Dashboard error:', error);
    res.status(500).json({
      error: 'Failed to load dashboard',
      message: error.message
    });
  }
});

// Admin-only fleet overview
router.get('/fleet', requireAdmin, async (req, res) => {
  try {
    const fleetData = {
      overview: {
        totalDevices: 25,
        onlineDevices: 23,
        offlineDevices: 2,
        totalContainers: 142,
        runningContainers: 138,
        stoppedContainers: 4
      },
      riskySystems: [
        {
          deviceId: 'web-server-01',
          riskScore: 85,
          issues: ['Outdated packages', 'Suspicious processes']
        },
        {
          deviceId: 'db-server-02',
          riskScore: 72,
          issues: ['High CPU usage', 'Unusual network traffic']
        }
      ],
      topThreats: [
        {
          type: 'malware',
          count: 5,
          severity: 'high'
        },
        {
          type: 'anomaly',
          count: 12,
          severity: 'medium'
        }
      ]
    };

    res.json(fleetData);
  } catch (error) {
    logger.error('Fleet overview error:', error);
    res.status(500).json({
      error: 'Failed to load fleet overview',
      message: error.message
    });
  }
});

module.exports = router;
