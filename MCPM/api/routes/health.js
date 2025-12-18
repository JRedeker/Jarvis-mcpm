/**
 * Health Check Routes
 *
 * GET /api/v1/health - System health check
 */

const express = require('express');
const { successResponse, isDockerAvailable, getConfigDir, REGISTRY_CONFIG_PATH } = require('../helpers');
const fs = require('fs');
const path = require('path');

const router = express.Router();

/**
 * GET /health
 * Returns system health status
 */
router.get('/health', (req, res) => {
    const checks = {
        node: {
            status: 'ok',
            version: process.version
        },
        registry: {
            status: fs.existsSync(REGISTRY_CONFIG_PATH) ? 'ok' : 'missing',
            path: REGISTRY_CONFIG_PATH
        },
        config_dir: {
            status: fs.existsSync(getConfigDir()) ? 'ok' : 'missing',
            path: getConfigDir()
        },
        docker: {
            status: isDockerAvailable() ? 'detected' : 'not_found'
        }
    };

    // Overall health status
    const allHealthy = checks.registry.status === 'ok' && checks.node.status === 'ok';

    res.json(successResponse({
        status: allHealthy ? 'healthy' : 'degraded',
        timestamp: new Date().toISOString(),
        version: '1.0.0',
        checks
    }));
});

module.exports = router;
