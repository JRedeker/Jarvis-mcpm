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
    MCPM_DIR,
    REGISTRY_CONFIG_PATH
};
