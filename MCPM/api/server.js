/**
 * MCPM API Server
 *
 * Provides a REST API for MCPM operations, replacing CLI subprocess calls.
 * This enables structured JSON communication instead of text parsing.
 *
 * Base URL: http://localhost:6275/api/v1
 */

const express = require('express');
const healthRoutes = require('./routes/health');
const serversRoutes = require('./routes/servers');
const profilesRoutes = require('./routes/profiles');
const clientsRoutes = require('./routes/clients');
const systemRoutes = require('./routes/system');

/**
 * Creates and configures the Express application
 * @returns {express.Application} Configured Express app
 */
function createApp() {
    const app = express();

    // Middleware
    app.use(express.json());

    // Request logging (simple middleware)
    app.use((req, res, next) => {
        const start = Date.now();
        res.on('finish', () => {
            const duration = Date.now() - start;
            console.log(`${req.method} ${req.path} ${res.statusCode} ${duration}ms`);
        });
        next();
    });

    // API Routes (versioned)
    const apiRouter = express.Router();
    apiRouter.use('/', healthRoutes);
    apiRouter.use('/servers', serversRoutes);
    apiRouter.use('/profiles', profilesRoutes);
    apiRouter.use('/clients', clientsRoutes);
    apiRouter.use('/', systemRoutes);

    app.use('/api/v1', apiRouter);

    // Root redirect to health
    app.get('/', (req, res) => {
        res.redirect('/api/v1/health');
    });

    // 404 handler
    app.use((req, res) => {
        res.status(404).json({
            success: false,
            data: null,
            error: {
                code: 'NOT_FOUND',
                message: `Endpoint not found: ${req.method} ${req.path}`
            }
        });
    });

    // Error handler
    app.use((err, req, res, next) => {
        console.error('Unhandled error:', err);
        res.status(500).json({
            success: false,
            data: null,
            error: {
                code: 'INTERNAL_ERROR',
                message: err.message || 'An unexpected error occurred'
            }
        });
    });

    return app;
}

/**
 * Starts the API server
 * @param {number} port - Port to listen on (default: 6275)
 * @param {string} host - Host to bind to (default: 0.0.0.0)
 * @returns {Promise<http.Server>} Running HTTP server
 */
function startServer(port = 6275, host = '0.0.0.0') {
    const app = createApp();

    return new Promise((resolve, reject) => {
        const server = app.listen(port, host, () => {
            console.log(`MCPM API Server listening on http://${host}:${port}`);
            console.log(`API Base: http://${host}:${port}/api/v1`);
            resolve(server);
        });

        server.on('error', reject);
    });
}

module.exports = { createApp, startServer };

// Run if executed directly (e.g., node api/server.js)
if (require.main === module) {
    const port = parseInt(process.env.MCPM_API_PORT || '6275', 10);
    const host = process.env.MCPM_API_HOST || '0.0.0.0';

    startServer(port, host)
        .then(() => {
            console.log('MCPM API Server started successfully');
        })
        .catch((err) => {
            console.error('Failed to start MCPM API Server:', err.message);
            process.exit(1);
        });
}
