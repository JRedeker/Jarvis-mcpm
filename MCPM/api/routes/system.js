/**
 * System Routes
 *
 * GET  /api/v1/usage    - Get usage statistics
 * POST /api/v1/migrate  - Migrate configuration
 * GET  /api/v1/search   - Search servers (top-level alias)
 */

const express = require('express');
const fs = require('fs');
const {
    successResponse,
    errorResponse,
    getConfigDir,
    getUserConfigPath,
    getProfilesConfigPath,
    getClientsConfigPath,
    getAllServers,
    readProfilesConfig,
    readClientsConfig
} = require('../helpers');

const router = express.Router();

/**
 * GET /usage
 * Get usage statistics for MCPM
 */
router.get('/usage', (req, res) => {
    try {
        const servers = getAllServers();
        const profiles = readProfilesConfig();
        const clients = readClientsConfig();

        // Count servers by category
        const serverCounts = {};
        let totalServers = 0;
        let installedCount = 0;

        for (const group in servers) {
            const groupServers = Object.keys(servers[group]).length;
            serverCounts[group] = groupServers;
            totalServers += groupServers;

            // Count installed (custom servers are always installed)
            if (group === 'custom') {
                installedCount += groupServers;
            }
        }

        res.json(successResponse({
            servers: {
                total: totalServers,
                installed: installedCount,
                byCategory: serverCounts
            },
            profiles: {
                total: Object.keys(profiles).length,
                list: Object.keys(profiles)
            },
            clients: {
                configured: Object.keys(clients).length,
                list: Object.keys(clients)
            },
            configDir: getConfigDir()
        }));
    } catch (err) {
        res.status(500).json(errorResponse('USAGE_ERROR', err.message));
    }
});

/**
 * POST /migrate
 * Migrate configuration to latest format
 */
router.post('/migrate', (req, res) => {
    try {
        const migrations = [];
        const configDir = getConfigDir();

        // Ensure config directory exists
        if (!fs.existsSync(configDir)) {
            fs.mkdirSync(configDir, { recursive: true });
            migrations.push('Created config directory');
        }

        // Check/create servers.json
        const serversPath = getUserConfigPath();
        if (!fs.existsSync(serversPath)) {
            fs.writeFileSync(serversPath, '{}');
            migrations.push('Created servers.json');
        }

        // Check/create profiles.json
        const profilesPath = getProfilesConfigPath();
        if (!fs.existsSync(profilesPath)) {
            fs.writeFileSync(profilesPath, '{}');
            migrations.push('Created profiles.json');
        }

        // Check/create clients.json
        const clientsPath = getClientsConfigPath();
        if (!fs.existsSync(clientsPath)) {
            fs.writeFileSync(clientsPath, '{}');
            migrations.push('Created clients.json');
        }

        // Check for legacy SSE configurations and convert to HTTP
        let legacyConversions = 0;
        const serversConfig = JSON.parse(fs.readFileSync(serversPath, 'utf-8'));
        for (const name in serversConfig) {
            if (serversConfig[name].type === 'sse') {
                serversConfig[name].type = 'streamable-http';
                legacyConversions++;
            }
        }

        if (legacyConversions > 0) {
            fs.writeFileSync(serversPath, JSON.stringify(serversConfig, null, 2));
            migrations.push(`Converted ${legacyConversions} SSE servers to streamable-http`);
        }

        if (migrations.length === 0) {
            migrations.push('Configuration is already up to date');
        }

        res.json(successResponse({
            migrations,
            configDir,
            message: 'Migration complete'
        }));
    } catch (err) {
        res.status(500).json(errorResponse('MIGRATE_ERROR', err.message));
    }
});

/**
 * GET /search
 * Top-level search endpoint (alias for /servers/search)
 */
router.get('/search', (req, res) => {
    const query = req.query.q || req.query.query || '';

    if (!query.trim()) {
        return res.status(400).json(errorResponse(
            'VALIDATION_ERROR',
            'Search query is required. Use ?q=<query>'
        ));
    }

    try {
        const servers = getAllServers();
        const results = [];
        const lowerQuery = query.toLowerCase();

        for (const group in servers) {
            for (const name in servers[group]) {
                const server = servers[group][name];
                const lowerName = name.toLowerCase();
                const lowerDesc = (server.description || '').toLowerCase();

                // Match on name or description
                if (lowerName.includes(lowerQuery) || lowerDesc.includes(lowerQuery)) {
                    results.push({
                        name,
                        group,
                        description: server.description || null,
                        type: server.type || 'stdio',
                        relevance: lowerName === lowerQuery ? 1.0 :
                                   lowerName.startsWith(lowerQuery) ? 0.8 :
                                   lowerName.includes(lowerQuery) ? 0.6 : 0.4
                    });
                }
            }
        }

        // Sort by relevance
        results.sort((a, b) => b.relevance - a.relevance);

        res.json(successResponse({
            query,
            count: results.length,
            results
        }));
    } catch (err) {
        res.status(500).json(errorResponse('SEARCH_ERROR', err.message));
    }
});

module.exports = router;
