/**
 * Profile Management Routes
 *
 * GET    /api/v1/profiles          - List all profiles
 * GET    /api/v1/profiles/:name    - Get profile info
 * POST   /api/v1/profiles          - Create profile
 * PUT    /api/v1/profiles/:name    - Edit profile
 * DELETE /api/v1/profiles/:name    - Delete profile
 */

const express = require('express');
const {
    successResponse,
    errorResponse,
    readProfilesConfig,
    saveProfilesConfig
} = require('../helpers');

const router = express.Router();

/**
 * GET /profiles
 * List all profiles
 */
router.get('/', (req, res) => {
    try {
        const profiles = readProfilesConfig();
        const result = [];

        for (const name in profiles) {
            const profile = profiles[name];
            result.push({
                name,
                servers: profile.servers || [],
                description: profile.description || null,
                port: profile.port || null
            });
        }

        res.json(successResponse({
            count: result.length,
            profiles: result
        }));
    } catch (err) {
        res.status(500).json(errorResponse('LIST_ERROR', err.message));
    }
});

/**
 * GET /profiles/:name
 * Get profile details
 */
router.get('/:name', (req, res) => {
    const { name } = req.params;

    try {
        const profiles = readProfilesConfig();

        if (!profiles[name]) {
            return res.status(404).json(errorResponse(
                'PROFILE_NOT_FOUND',
                `Profile '${name}' not found`
            ));
        }

        res.json(successResponse({
            name,
            ...profiles[name]
        }));
    } catch (err) {
        res.status(500).json(errorResponse('INFO_ERROR', err.message));
    }
});

/**
 * POST /profiles
 * Create a new profile
 */
router.post('/', (req, res) => {
    const { name, servers, description, port } = req.body;

    if (!name) {
        return res.status(400).json(errorResponse(
            'VALIDATION_ERROR',
            'Profile name is required'
        ));
    }

    try {
        const profiles = readProfilesConfig();

        if (profiles[name]) {
            return res.status(409).json(errorResponse(
                'PROFILE_EXISTS',
                `Profile '${name}' already exists`
            ));
        }

        profiles[name] = {
            servers: servers || [],
            ...(description && { description }),
            ...(port && { port })
        };

        saveProfilesConfig(profiles);

        res.status(201).json(successResponse({
            name,
            ...profiles[name],
            message: `Successfully created profile '${name}'`
        }));
    } catch (err) {
        res.status(500).json(errorResponse('CREATE_ERROR', err.message));
    }
});

/**
 * PUT /profiles/:name
 * Edit an existing profile
 */
router.put('/:name', (req, res) => {
    const { name } = req.params;
    const { new_name, add_servers, remove_servers, description, port } = req.body;

    try {
        const profiles = readProfilesConfig();

        if (!profiles[name]) {
            return res.status(404).json(errorResponse(
                'PROFILE_NOT_FOUND',
                `Profile '${name}' not found`
            ));
        }

        let profile = profiles[name];

        // Add servers
        if (add_servers) {
            const serversToAdd = Array.isArray(add_servers) ? add_servers : add_servers.split(',').map(s => s.trim());
            profile.servers = [...new Set([...(profile.servers || []), ...serversToAdd])];
        }

        // Remove servers
        if (remove_servers) {
            const serversToRemove = Array.isArray(remove_servers) ? remove_servers : remove_servers.split(',').map(s => s.trim());
            profile.servers = (profile.servers || []).filter(s => !serversToRemove.includes(s));
        }

        // Update description
        if (description !== undefined) {
            profile.description = description || null;
        }

        // Update port
        if (port !== undefined) {
            profile.port = port;
        }

        // Handle rename
        let finalName = name;
        if (new_name && new_name !== name) {
            if (profiles[new_name]) {
                return res.status(409).json(errorResponse(
                    'PROFILE_EXISTS',
                    `Profile '${new_name}' already exists`
                ));
            }
            profiles[new_name] = profile;
            delete profiles[name];
            finalName = new_name;
        } else {
            profiles[name] = profile;
        }

        saveProfilesConfig(profiles);

        res.json(successResponse({
            name: finalName,
            ...profiles[finalName],
            message: `Successfully updated profile '${finalName}'`
        }));
    } catch (err) {
        res.status(500).json(errorResponse('EDIT_ERROR', err.message));
    }
});

/**
 * DELETE /profiles/:name
 * Delete a profile
 */
router.delete('/:name', (req, res) => {
    const { name } = req.params;

    try {
        const profiles = readProfilesConfig();

        if (!profiles[name]) {
            return res.status(404).json(errorResponse(
                'PROFILE_NOT_FOUND',
                `Profile '${name}' not found`
            ));
        }

        delete profiles[name];
        saveProfilesConfig(profiles);

        res.json(successResponse({
            name,
            message: `Successfully deleted profile '${name}'`
        }));
    } catch (err) {
        res.status(500).json(errorResponse('DELETE_ERROR', err.message));
    }
});

module.exports = router;
