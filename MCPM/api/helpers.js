/**
 * MCPM API Helper Functions
 *
 * Shared utilities for config management and response formatting
 */

const fs = require('fs');
const path = require('path');
const os = require('os');
const toml = require('toml');
const { execSync } = require('child_process');

// Paths
const MCPM_DIR = path.dirname(__dirname);
const REGISTRY_CONFIG_PATH = path.join(MCPM_DIR, 'config', 'technologies.toml');

/**
 * Get user config directory path
 * @returns {string} Path to MCPM config directory
 */
function getConfigDir() {
    return path.join(os.homedir(), '.config', 'mcpm');
}

/**
 * Get user servers config path
 * @returns {string} Path to servers.json
 */
function getUserConfigPath() {
    return path.join(getConfigDir(), 'servers.json');
}

/**
 * Get profiles config path
 * @returns {string} Path to profiles.json
 */
function getProfilesConfigPath() {
    return path.join(getConfigDir(), 'profiles.json');
}

/**
 * Get clients config path
 * @returns {string} Path to clients.json
 */
function getClientsConfigPath() {
    return path.join(getConfigDir(), 'clients.json');
}

/**
 * Read registry config (TOML format)
 * @returns {Object} Registry configuration
 */
function readRegistryConfig() {
    if (!fs.existsSync(REGISTRY_CONFIG_PATH)) {
        throw new Error(`Registry config not found at ${REGISTRY_CONFIG_PATH}`);
    }
    const content = fs.readFileSync(REGISTRY_CONFIG_PATH, 'utf-8');
    return toml.parse(content);
}

/**
 * Read user config (JSON format)
 * @param {string} configPath - Path to config file
 * @returns {Object} User configuration (empty object if not exists)
 */
function readJsonConfig(configPath) {
    if (!fs.existsSync(configPath)) {
        return {};
    }
    try {
        return JSON.parse(fs.readFileSync(configPath, 'utf-8'));
    } catch (e) {
        console.error(`Failed to parse config at ${configPath}:`, e.message);
        return {};
    }
}

/**
 * Write JSON config file
 * @param {string} configPath - Path to config file
 * @param {Object} data - Data to write
 */
function writeJsonConfig(configPath, data) {
    const dir = path.dirname(configPath);
    if (!fs.existsSync(dir)) {
        fs.mkdirSync(dir, { recursive: true });
    }
    fs.writeFileSync(configPath, JSON.stringify(data, null, 2));
}

/**
 * Read user servers config
 * @returns {Object} User server configurations
 */
function readUserConfig() {
    return readJsonConfig(getUserConfigPath());
}

/**
 * Write user servers config
 * @param {Object} data - Server configurations
 */
function saveUserConfig(data) {
    writeJsonConfig(getUserConfigPath(), data);
}

/**
 * Read profiles config
 * @returns {Object} Profile configurations
 */
function readProfilesConfig() {
    return readJsonConfig(getProfilesConfigPath());
}

/**
 * Write profiles config
 * @param {Object} data - Profile configurations
 */
function saveProfilesConfig(data) {
    writeJsonConfig(getProfilesConfigPath(), data);
}

/**
 * Read clients config
 * @returns {Object} Client configurations
 */
function readClientsConfig() {
    return readJsonConfig(getClientsConfigPath());
}

/**
 * Write clients config
 * @param {Object} data - Client configurations
 */
function saveClientsConfig(data) {
    writeJsonConfig(getClientsConfigPath(), data);
}

/**
 * Get all servers (registry + user-defined)
 * @returns {Object} Combined server configurations grouped by category
 */
function getAllServers() {
    const registry = readRegistryConfig();
    const user = readUserConfig();

    const all = { ...registry.technologies };
    if (Object.keys(user).length > 0) {
        all['custom'] = user;
    }
    return all;
}

/**
 * Find a server by name across all categories
 * @param {string} name - Server name
 * @returns {Object|null} Server config or null if not found
 */
function findServer(name) {
    const servers = getAllServers();

    // Check custom/user first
    if (servers['custom'] && servers['custom'][name]) {
        return { ...servers['custom'][name], source: 'custom' };
    }

    // Then check registry
    for (const group in servers) {
        if (servers[group][name]) {
            return { ...servers[group][name], source: group };
        }
    }

    return null;
}

/**
 * Check if Docker is available
 * @returns {boolean} True if Docker is detected
 */
function isDockerAvailable() {
    try {
        execSync('docker --version', { stdio: 'ignore' });
        return true;
    } catch {
        return false;
    }
}

/**
 * Create a success API response
 * @param {*} data - Response data
 * @returns {Object} Formatted success response
 */
function successResponse(data) {
    return {
        success: true,
        data,
        error: null
    };
}

/**
 * Create an error API response
 * @param {string} code - Error code
 * @param {string} message - Error message
 * @param {Object} [details] - Additional error details
 * @returns {Object} Formatted error response
 */
function errorResponse(code, message, details = null) {
    return {
        success: false,
        data: null,
        error: {
            code,
            message,
            ...(details && { details })
        }
    };
}

/**
 * Synchronize server profile_tags with profiles.json
 * When a server is added/removed from a profile, update the server's profile_tags
 * @param {string} profileName - Profile name
 * @param {string[]} serversToAdd - Servers being added to the profile
 * @param {string[]} serversToRemove - Servers being removed from the profile
 */
function syncServerProfileTags(profileName, serversToAdd = [], serversToRemove = []) {
    const servers = readUserConfig();
    let modified = false;

    // Add profile tag to servers being added
    for (const serverName of serversToAdd) {
        if (servers[serverName]) {
            const tags = servers[serverName].profile_tags || [];
            if (!tags.includes(profileName)) {
                servers[serverName].profile_tags = [...tags, profileName];
                modified = true;
            }
        }
    }

    // Remove profile tag from servers being removed
    for (const serverName of serversToRemove) {
        if (servers[serverName]) {
            const tags = servers[serverName].profile_tags || [];
            if (tags.includes(profileName)) {
                servers[serverName].profile_tags = tags.filter(t => t !== profileName);
                modified = true;
            }
        }
    }

    if (modified) {
        saveUserConfig(servers);
    }

    return modified;
}

/**
 * Audit and reconcile profile_tags in servers.json with profiles.json
 * Returns a report of mismatches and optionally fixes them
 * @param {boolean} autoFix - If true, automatically fix mismatches
 * @returns {Object} Audit report with mismatches and actions taken
 */
function auditProfileSync(autoFix = false) {
    const servers = readUserConfig();
    const profiles = readProfilesConfig();
    const report = {
        mismatches: [],
        fixes: [],
        serverTagsNotInProfile: [],
        profileServersWithoutTag: []
    };

    // Build a map of profile -> servers from profiles.json
    const profileServerMap = {};
    for (const profileName in profiles) {
        profileServerMap[profileName] = new Set(profiles[profileName].servers || []);
    }

    // Check each server's profile_tags against actual profile membership
    for (const serverName in servers) {
        const serverTags = new Set(servers[serverName].profile_tags || []);

        // Check if server claims to be in profiles it's not actually in
        for (const tag of serverTags) {
            if (profileServerMap[tag] && !profileServerMap[tag].has(serverName)) {
                report.serverTagsNotInProfile.push({
                    server: serverName,
                    tag: tag,
                    issue: `Server has tag '${tag}' but is not in profile's server list`
                });
                report.mismatches.push({ type: 'tag_without_membership', server: serverName, profile: tag });

                if (autoFix) {
                    // Remove the tag from the server since it's not in the profile
                    servers[serverName].profile_tags = servers[serverName].profile_tags.filter(t => t !== tag);
                    report.fixes.push(`Removed tag '${tag}' from server '${serverName}'`);
                }
            }
        }
    }

    // Check each profile's servers have the corresponding tag
    for (const profileName in profiles) {
        const profileServers = profiles[profileName].servers || [];
        for (const serverName of profileServers) {
            if (servers[serverName]) {
                const serverTags = servers[serverName].profile_tags || [];
                if (!serverTags.includes(profileName)) {
                    report.profileServersWithoutTag.push({
                        server: serverName,
                        profile: profileName,
                        issue: `Server in profile '${profileName}' but missing the profile tag`
                    });
                    report.mismatches.push({ type: 'membership_without_tag', server: serverName, profile: profileName });

                    if (autoFix) {
                        // Add the tag to the server
                        servers[serverName].profile_tags = [...serverTags, profileName];
                        report.fixes.push(`Added tag '${profileName}' to server '${serverName}'`);
                    }
                }
            }
        }
    }

    if (autoFix && report.fixes.length > 0) {
        saveUserConfig(servers);
    }

    report.summary = {
        totalMismatches: report.mismatches.length,
        fixesApplied: report.fixes.length,
        isInSync: report.mismatches.length === 0
    };

    return report;
}

module.exports = {
    getConfigDir,
    getUserConfigPath,
    getProfilesConfigPath,
    getClientsConfigPath,
    readRegistryConfig,
    readJsonConfig,
    writeJsonConfig,
    readUserConfig,
    saveUserConfig,
    readProfilesConfig,
    saveProfilesConfig,
    readClientsConfig,
    saveClientsConfig,
    getAllServers,
    findServer,
    isDockerAvailable,
    successResponse,
    errorResponse,
    syncServerProfileTags,
    auditProfileSync,
    MCPM_DIR,
    REGISTRY_CONFIG_PATH
};
