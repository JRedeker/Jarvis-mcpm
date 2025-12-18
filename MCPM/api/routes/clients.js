/**
 * Client Management Routes
 *
 * GET /api/v1/clients          - List all clients
 * GET /api/v1/clients/:name    - Get client info
 * PUT /api/v1/clients/:name    - Edit client config
 */

const express = require('express');
const fs = require('fs');
const os = require('os');
const path = require('path');
const {
    successResponse,
    errorResponse,
    readClientsConfig,
    saveClientsConfig
} = require('../helpers');

const router = express.Router();

// Known MCP clients and their default config paths
const KNOWN_CLIENTS = {
    'opencode': {
        displayName: 'OpenCode',
        configPaths: [
            path.join(os.homedir(), '.config', 'opencode', 'opencode.json'),
            './opencode.json'
        ]
    },
    'claude-code': {
        displayName: 'Claude Code CLI',
        configPaths: [
            path.join(os.homedir(), '.claude.json')
        ]
    },
    'claude-desktop': {
        displayName: 'Claude Desktop',
        configPaths: [
            path.join(os.homedir(), '.config', 'Claude', 'claude_desktop_config.json'),
            path.join(os.homedir(), 'Library', 'Application Support', 'Claude', 'claude_desktop_config.json')
        ]
    }
};

/**
 * GET /clients
 * List all known clients and their status
 */
router.get('/', (req, res) => {
    try {
        const clientsConfig = readClientsConfig();
        const result = [];

        for (const name in KNOWN_CLIENTS) {
            const client = KNOWN_CLIENTS[name];
            const customConfig = clientsConfig[name] || {};

            // Try to detect the client
            let detected = false;
            let configPath = customConfig.configPath || null;

            if (!configPath) {
                for (const p of client.configPaths) {
                    if (fs.existsSync(p)) {
                        detected = true;
                        configPath = p;
                        break;
                    }
                }
            } else {
                detected = fs.existsSync(configPath);
            }

            result.push({
                name,
                displayName: client.displayName,
                detected,
                configPath,
                profiles: customConfig.profiles || [],
                servers: customConfig.servers || []
            });
        }

        res.json(successResponse({
            count: result.length,
            clients: result
        }));
    } catch (err) {
        res.status(500).json(errorResponse('LIST_ERROR', err.message));
    }
});

/**
 * GET /clients/:name
 * Get detailed client info
 */
router.get('/:name', (req, res) => {
    const { name } = req.params;

    try {
        const client = KNOWN_CLIENTS[name];
        if (!client) {
            return res.status(404).json(errorResponse(
                'CLIENT_NOT_FOUND',
                `Unknown client '${name}'. Known clients: ${Object.keys(KNOWN_CLIENTS).join(', ')}`
            ));
        }

        const clientsConfig = readClientsConfig();
        const customConfig = clientsConfig[name] || {};

        // Try to detect and read config
        let detected = false;
        let configPath = customConfig.configPath || null;
        let configContents = null;

        if (!configPath) {
            for (const p of client.configPaths) {
                if (fs.existsSync(p)) {
                    detected = true;
                    configPath = p;
                    break;
                }
            }
        } else {
            detected = fs.existsSync(configPath);
        }

        if (detected && configPath) {
            try {
                configContents = JSON.parse(fs.readFileSync(configPath, 'utf-8'));
            } catch (e) {
                // Config exists but couldn't be parsed
            }
        }

        res.json(successResponse({
            name,
            displayName: client.displayName,
            detected,
            configPath,
            possiblePaths: client.configPaths,
            profiles: customConfig.profiles || [],
            servers: customConfig.servers || [],
            config: configContents
        }));
    } catch (err) {
        res.status(500).json(errorResponse('INFO_ERROR', err.message));
    }
});

/**
 * PUT /clients/:name
 * Edit client configuration
 */
router.put('/:name', (req, res) => {
    const { name } = req.params;
    const {
        config_path,
        add_server,
        remove_server,
        add_profile,
        remove_profile
    } = req.body;

    try {
        const client = KNOWN_CLIENTS[name];
        if (!client) {
            return res.status(404).json(errorResponse(
                'CLIENT_NOT_FOUND',
                `Unknown client '${name}'. Known clients: ${Object.keys(KNOWN_CLIENTS).join(', ')}`
            ));
        }

        const clientsConfig = readClientsConfig();
        if (!clientsConfig[name]) {
            clientsConfig[name] = { profiles: [], servers: [] };
        }

        const changes = [];

        // Set config path
        if (config_path) {
            clientsConfig[name].configPath = config_path;
            changes.push(`Set config path to: ${config_path}`);
        }

        // Add server
        if (add_server) {
            const servers = Array.isArray(add_server) ? add_server : [add_server];
            clientsConfig[name].servers = [...new Set([...(clientsConfig[name].servers || []), ...servers])];
            changes.push(`Added servers: ${servers.join(', ')}`);
        }

        // Remove server
        if (remove_server) {
            const servers = Array.isArray(remove_server) ? remove_server : [remove_server];
            clientsConfig[name].servers = (clientsConfig[name].servers || []).filter(s => !servers.includes(s));
            changes.push(`Removed servers: ${servers.join(', ')}`);
        }

        // Add profile
        if (add_profile) {
            const profiles = Array.isArray(add_profile) ? add_profile : add_profile.split(',').map(p => p.trim());
            clientsConfig[name].profiles = [...new Set([...(clientsConfig[name].profiles || []), ...profiles])];
            changes.push(`Added profiles: ${profiles.join(', ')}`);
        }

        // Remove profile
        if (remove_profile) {
            const profiles = Array.isArray(remove_profile) ? remove_profile : remove_profile.split(',').map(p => p.trim());
            clientsConfig[name].profiles = (clientsConfig[name].profiles || []).filter(p => !profiles.includes(p));
            changes.push(`Removed profiles: ${profiles.join(', ')}`);
        }

        if (changes.length === 0) {
            return res.status(400).json(errorResponse(
                'NO_CHANGES',
                'No changes specified. Use config_path, add_server, remove_server, add_profile, or remove_profile.'
            ));
        }

        saveClientsConfig(clientsConfig);

        res.json(successResponse({
            name,
            displayName: client.displayName,
            ...clientsConfig[name],
            changes,
            message: `Successfully updated client '${name}'`
        }));
    } catch (err) {
        res.status(500).json(errorResponse('EDIT_ERROR', err.message));
    }
});

module.exports = router;
