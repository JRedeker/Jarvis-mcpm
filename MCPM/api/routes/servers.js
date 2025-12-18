/**
 * Server Management Routes
 *
 * GET    /api/v1/servers           - List all servers
 * GET    /api/v1/servers/:name     - Get server info
 * POST   /api/v1/servers/:name/install - Install server
 * DELETE /api/v1/servers/:name     - Uninstall server
 * POST   /api/v1/servers           - Create custom server
 * PUT    /api/v1/servers/:name     - Edit server config
 * GET    /api/v1/search            - Search servers
 */

const express = require('express');
const { execSync } = require('child_process');
const path = require('path');
const {
    successResponse,
    errorResponse,
    getAllServers,
    findServer,
    readUserConfig,
    saveUserConfig,
    MCPM_DIR
} = require('../helpers');

const router = express.Router();

/**
 * GET /search
 * Search servers by query
 * IMPORTANT: Must be defined BEFORE /:name to avoid route collision
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

/**
 * GET /servers
 * List all servers (registry + custom)
 */
router.get('/', (req, res) => {
    try {
        const servers = getAllServers();
        const result = [];

        for (const group in servers) {
            for (const name in servers[group]) {
                const server = servers[group][name];
                result.push({
                    name,
                    group,
                    description: server.description || null,
                    type: server.type || 'stdio',
                    installed: isServerInstalled(name, server, group)
                });
            }
        }

        res.json(successResponse({
            count: result.length,
            servers: result
        }));
    } catch (err) {
        res.status(500).json(errorResponse('LIST_ERROR', err.message));
    }
});

/**
 * GET /servers/:name
 * Get detailed info for a specific server
 */
router.get('/:name', (req, res) => {
    const { name } = req.params;

    try {
        const server = findServer(name);
        if (!server) {
            return res.status(404).json(errorResponse(
                'SERVER_NOT_FOUND',
                `Server '${name}' not found in registry or custom servers`
            ));
        }

        res.json(successResponse({
            name,
            ...server,
            installed: isServerInstalled(name, server, server.source)
        }));
    } catch (err) {
        res.status(500).json(errorResponse('INFO_ERROR', err.message));
    }
});

/**
 * POST /servers/:name/install
 * Install a server from registry
 */
router.post('/:name/install', (req, res) => {
    const { name } = req.params;

    try {
        const servers = getAllServers();
        let found = false;
        let pkgName = name;
        let dockerImage = null;

        // Find server in registry
        for (const group in servers) {
            if (servers[group][name]) {
                found = true;
                const tech = servers[group][name];
                if (tech.docker) {
                    dockerImage = tech.docker;
                }
                if (tech.package) {
                    pkgName = tech.package;
                } else if (tech.repo) {
                    pkgName = `git+${tech.repo}.git`;
                }
                break;
            }
        }

        if (!found) {
            return res.status(404).json(errorResponse(
                'SERVER_NOT_FOUND',
                `Server '${name}' not found in registry. Use POST /servers to create a custom server.`
            ));
        }

        // Try Docker first if available
        if (dockerImage) {
            try {
                execSync(`docker pull ${dockerImage}`, { stdio: 'pipe' });
                return res.json(successResponse({
                    name,
                    method: 'docker',
                    image: dockerImage,
                    message: `Successfully pulled Docker image for ${name}`
                }));
            } catch (dockerErr) {
                // Fall through to npm
                console.log(`Docker pull failed, falling back to npm: ${dockerErr.message}`);
            }
        }

        // Install via npm
        execSync(`npm install ${pkgName} --save`, {
            stdio: 'pipe',
            cwd: MCPM_DIR
        });

        res.json(successResponse({
            name,
            method: 'npm',
            package: pkgName,
            message: `Successfully installed ${name}`
        }));
    } catch (err) {
        res.status(500).json(errorResponse('INSTALL_ERROR', err.message));
    }
});

/**
 * DELETE /servers/:name
 * Uninstall a server
 */
router.delete('/:name', (req, res) => {
    const { name } = req.params;

    try {
        // Check if it's a custom server
        const userConfig = readUserConfig();
        if (userConfig[name]) {
            delete userConfig[name];
            saveUserConfig(userConfig);
            return res.json(successResponse({
                name,
                method: 'custom',
                message: `Successfully removed custom server '${name}'`
            }));
        }

        // Check if server is installed via npm before trying to uninstall
        let isInstalled = false;
        try {
            const pkgJson = require(path.join(MCPM_DIR, 'package.json'));
            isInstalled = !!(pkgJson.dependencies && pkgJson.dependencies[name]);
        } catch (e) {
            // package.json doesn't exist or is invalid
        }

        if (!isInstalled) {
            return res.status(404).json(errorResponse(
                'SERVER_NOT_INSTALLED',
                `Server '${name}' is not installed`
            ));
        }

        // Uninstall from npm
        try {
            execSync(`npm uninstall ${name}`, {
                stdio: 'pipe',
                cwd: MCPM_DIR
            });
            res.json(successResponse({
                name,
                method: 'npm',
                message: `Successfully uninstalled '${name}'`
            }));
        } catch (npmErr) {
            res.status(500).json(errorResponse(
                'UNINSTALL_ERROR',
                `Failed to uninstall '${name}': ${npmErr.message}`
            ));
        }
    } catch (err) {
        res.status(500).json(errorResponse('UNINSTALL_ERROR', err.message));
    }
});

/**
 * POST /servers
 * Create a new custom server
 */
router.post('/', (req, res) => {
    const { name, type, command, args, url, env, headers, force } = req.body;

    // Validate required fields
    if (!name) {
        return res.status(400).json(errorResponse(
            'VALIDATION_ERROR',
            'Server name is required'
        ));
    }

    // SSE is no longer supported - check this FIRST for a specific error message
    if (type === 'sse') {
        return res.status(400).json(errorResponse(
            'VALIDATION_ERROR',
            "SSE transport has been removed. Use 'streamable-http' instead."
        ));
    }

    if (!type || !['stdio', 'http', 'streamable-http'].includes(type)) {
        return res.status(400).json(errorResponse(
            'VALIDATION_ERROR',
            "Type must be one of: 'stdio', 'http', 'streamable-http'"
        ));
    }

    // Validate type-specific requirements
    if (type === 'stdio' && !command) {
        return res.status(400).json(errorResponse(
            'VALIDATION_ERROR',
            'Command is required for stdio type'
        ));
    }

    if (['http', 'streamable-http'].includes(type) && !url) {
        return res.status(400).json(errorResponse(
            'VALIDATION_ERROR',
            `URL is required for ${type} type`
        ));
    }

    try {
        const userConfig = readUserConfig();

        if (userConfig[name] && !force) {
            return res.status(409).json(errorResponse(
                'SERVER_EXISTS',
                `Server '${name}' already exists. Use force=true to overwrite.`
            ));
        }

        // Build server config
        const serverConfig = { type };
        if (command) serverConfig.command = command;
        if (args) serverConfig.args = Array.isArray(args) ? args : args.split(' ');
        if (url) serverConfig.url = url;
        if (env) serverConfig.env = parseKeyValue(env);
        if (headers) serverConfig.headers = parseKeyValue(headers);

        userConfig[name] = serverConfig;
        saveUserConfig(userConfig);

        res.status(201).json(successResponse({
            name,
            ...serverConfig,
            message: `Successfully registered server '${name}'`
        }));
    } catch (err) {
        res.status(500).json(errorResponse('CREATE_ERROR', err.message));
    }
});

/**
 * PUT /servers/:name
 * Edit an existing server configuration
 */
router.put('/:name', (req, res) => {
    const { name } = req.params;
    const { command, args, url, env, headers } = req.body;

    try {
        const userConfig = readUserConfig();

        if (!userConfig[name]) {
            return res.status(404).json(errorResponse(
                'SERVER_NOT_FOUND',
                `Custom server '${name}' not found. Only custom servers can be edited.`
            ));
        }

        // Update fields if provided
        if (command !== undefined) userConfig[name].command = command;
        if (args !== undefined) {
            userConfig[name].args = Array.isArray(args) ? args : args.split(' ');
        }
        if (url !== undefined) userConfig[name].url = url;
        if (env !== undefined) {
            userConfig[name].env = parseKeyValue(env);
        }
        if (headers !== undefined) {
            userConfig[name].headers = parseKeyValue(headers);
        }

        saveUserConfig(userConfig);

        res.json(successResponse({
            name,
            ...userConfig[name],
            message: `Successfully updated server '${name}'`
        }));
    } catch (err) {
        res.status(500).json(errorResponse('EDIT_ERROR', err.message));
    }
});

/**
 * Check if a server is installed
 * @param {string} name - Server name
 * @param {Object} tech - Server configuration
 * @param {string} group - Server group
 * @returns {boolean} True if installed
 */
function isServerInstalled(name, tech, group) {
    // Custom servers are always "installed" (registered)
    if (group === 'custom') {
        return true;
    }

    // Check package.json dependencies
    try {
        const pkgJson = require(path.join(MCPM_DIR, 'package.json'));
        const pkgName = tech.package || name;
        return !!(pkgJson.dependencies && pkgJson.dependencies[pkgName]);
    } catch (e) {
        return false;
    }
}

/**
 * Parse KEY=value pairs from string or object
 * @param {string|Object} input - Input to parse
 * @returns {Object} Parsed key-value object
 */
function parseKeyValue(input) {
    if (typeof input === 'object') return input;
    if (typeof input !== 'string') return {};

    const result = {};
    input.split(',').forEach(pair => {
        const [key, ...valueParts] = pair.split('=');
        if (key && valueParts.length > 0) {
            result[key.trim()] = valueParts.join('=').trim();
        }
    });
    return result;
}

module.exports = router;
