#!/usr/bin/env node

const { Command } = require('commander');
const fs = require('fs');
const path = require('path');
const os = require('os');
const toml = require('toml');
const { execSync } = require('child_process');
const chalk = require('chalk');

const program = new Command();
const configPath = path.join(__dirname, 'config', 'technologies.toml');

// Helper to read registry config (TOML)
function readRegistryConfig() {
    if (!fs.existsSync(configPath)) {
        console.error(chalk.red(`Config file not found at ${configPath}`));
        process.exit(1);
    }
    try {
        const content = fs.readFileSync(configPath, 'utf-8');
        return toml.parse(content);
    } catch (e) {
        console.error(chalk.red(`Failed to parse config: ${e.message}`));
        process.exit(1);
    }
}

// Helper to get user config path
function getUserConfigPath() {
    return path.join(os.homedir(), '.config', 'mcpm', 'servers.json');
}

// Helper to read user config
function readUserConfig() {
    const p = getUserConfigPath();
    if (!fs.existsSync(p)) return {};
    try {
        return JSON.parse(fs.readFileSync(p, 'utf-8'));
    } catch (e) {
        return {};
    }
}

// Helper to save user config
function saveUserConfig(data) {
    const p = getUserConfigPath();
    const dir = path.dirname(p);
    if (!fs.existsSync(dir)) fs.mkdirSync(dir, { recursive: true });
    fs.writeFileSync(p, JSON.stringify(data, null, 2));
}

// Helper to get all servers (registry + user)
function getAllServers() {
    const registry = readRegistryConfig();
    const user = readUserConfig();

    // Structure: { group: { name: config } }
    const all = { ...registry.technologies };
    if (Object.keys(user).length > 0) {
        all['custom'] = user;
    }
    return all;
}

program
    .name('mcpm')
    .description('Model Context Protocol Manager CLI')
    .version('1.0.0');

program
    .command('ls')
    .description('List installed MCP servers')
    .action(() => {
        const technologies = getAllServers();
        console.log(chalk.bold('Installed MCP Servers:'));

        let hasInstalled = false;
        for (const group in technologies) {
            let groupPrinted = false;
            for (const key in technologies[group]) {
                const tech = technologies[group][key];
                let isInstalled = false;

                // Check if installed in node_modules
                let pkgName = tech.package || key;

                // Handle scoped packages
                if (!tech.package && tech.repo && tech.repo.includes('github.com/')) {
                    try {
                        const pkgJson = require(path.join(__dirname, 'package.json'));
                        if (pkgJson.dependencies && pkgJson.dependencies[pkgName]) {
                            isInstalled = true;
                        }
                    } catch (e) {}
                } else {
                        try {
                        const pkgJson = require(path.join(__dirname, 'package.json'));
                            if (pkgJson.dependencies && pkgJson.dependencies[pkgName]) {
                            isInstalled = true;
                        }
                    } catch (e) {}
                }

                // If it's a user-defined server, it's considered "installed" (registered)
                if (group === 'custom') {
                    isInstalled = true;
                }

                if (isInstalled) {
                    if (!groupPrinted) {
                        console.log(chalk.blue(`\n[${group}]`));
                        groupPrinted = true;
                    }
                    const desc = tech.description ? `: ${tech.description}` : '';
                    console.log(`- ${key}${desc} ${chalk.green('(Installed)')}`);
                    hasInstalled = true;
                }
            }
        }

        if (!hasInstalled) {
            console.log(chalk.gray('No servers currently installed. Use "mcpm install <name>" to add one.'));
        }
    });

program
    .command('install <name>')
    .description('Install an MCP server')
    .action((name) => {
        console.log(chalk.cyan(`Installing ${name}...`));
        const technologies = getAllServers();
        let found = false;
        let pkgName = name;
        let dockerImage = null;

        // Simple lookup strategy
        for (const group in technologies) {
             if (technologies[group][name]) {
                 found = true;
                 const tech = technologies[group][name];
                 if (tech.docker) {
                     dockerImage = tech.docker;
                 }
                 if (tech.package) {
                     pkgName = tech.package;
                 } else if (tech.repo) {
                     pkgName = `git+${tech.repo}.git`;
                 }
             }
        }

        if (!found) {
             console.log(chalk.yellow(`Warning: ${name} not found in registry. Attempting direct npm install...`));
        }

        // Helper to print config snippet
        const printConfigSnippet = (toolName, type, target) => {
            console.log(chalk.magenta(`\nConfiguration Snippet (Add to your Client Config):`));
            let config = {};
            if (type === 'docker') {
                config = {
                    command: "docker",
                    args: ["run", "-i", "--rm", target]
                };
            } else {
                // npm/node
                config = {
                    command: "node",
                    args: [path.join(__dirname, 'node_modules', target, 'index.js')]
                };
            }
            console.log(JSON.stringify({ [toolName]: config }, null, 2));
        };

        // Docker Preference Logic
        if (dockerImage) {
            console.log(chalk.blue(`Docker image found: ${dockerImage}. Preferring Docker installation...`));
            try {
                execSync(`docker pull ${dockerImage}`, { stdio: 'inherit' });
                console.log(chalk.green(`Successfully pulled Docker image for ${name}`));
                printConfigSnippet(name, 'docker', dockerImage);
                return; // Exit after docker install
            } catch (e) {
                console.log(chalk.red(`Failed to pull docker image: ${e.message}`));
                console.log(chalk.yellow(`Falling back to npm/git installation...`));
                // Fallthrough to npm
            }
        }

        try {
            // Use --save to ensure it's added to package.json
            execSync(`npm install ${pkgName} --save`, { stdio: 'inherit', cwd: __dirname });
            console.log(chalk.green(`Successfully installed ${name}`));

            let folderName = pkgName;
            if (pkgName.startsWith('git+')) {
                folderName = name;
            }
            printConfigSnippet(name, 'npm', folderName);
        } catch (e) {
            console.error(chalk.red(`Failed to install ${name}: ${e.message}`));
            process.exit(1);
        }
    });

program
    .command('new <name>')
    .description('Register a new custom MCP server')
    .option('--type <type>', 'Transport type: stdio or streamable-http', 'stdio')
    .option('--command <command>', 'Command to run (for stdio)')
    .option('--args <args>', 'Arguments for command (space-separated)')
    .option('--env <env>', 'Environment variables (KEY=value, comma-separated)')
    .option('--url <url>', 'URL (for streamable-http)')
    .option('--headers <headers>', 'Headers (KEY=value, comma-separated)')
    .option('--force', 'Overwrite existing server')
    .action((name, options) => {
        // Check for removed SSE type with helpful error message
        if (options.type === 'sse') {
            console.error(chalk.red(`Error: Invalid type 'sse'. SSE transport has been removed.`));
            console.error(chalk.yellow(`Use 'streamable-http' instead: mcpm new ${name} --type streamable-http --url <url>`));
            process.exit(1);
        }

        // Validate type
        const validTypes = ['stdio', 'http', 'streamable-http'];
        if (!validTypes.includes(options.type)) {
            console.error(chalk.red(`Invalid type: ${options.type}. Must be one of: ${validTypes.join(', ')}`));
            process.exit(1);
        }

        // Validate required params
        if (options.type === 'stdio' && !options.command) {
            console.error(chalk.red('Error: --command is required for stdio type'));
            process.exit(1);
        }
        if (['http', 'streamable-http'].includes(options.type) && !options.url) {
            console.error(chalk.red(`Error: --url is required for ${options.type} type`));
            process.exit(1);
        }

        const userConfig = readUserConfig();
        if (userConfig[name] && !options.force) {
            console.error(chalk.red(`Error: Server '${name}' already exists. Use --force to overwrite.`));
            process.exit(1);
        }

        // Build server config object
        const serverConfig = {
            type: options.type,
        };

        if (options.command) serverConfig.command = options.command;
        if (options.args) serverConfig.args = options.args.split(' '); // Basic splitting
        if (options.url) serverConfig.url = options.url;

        if (options.env) {
            serverConfig.env = {};
            options.env.split(',').forEach(pair => {
                const [k, v] = pair.split('=');
                if (k && v) serverConfig.env[k.trim()] = v.trim();
            });
        }

        if (options.headers) {
            serverConfig.headers = {};
            options.headers.split(',').forEach(pair => {
                const [k, v] = pair.split('=');
                if (k && v) serverConfig.headers[k.trim()] = v.trim();
            });
        }

        userConfig[name] = serverConfig;
        saveUserConfig(userConfig);
        console.log(chalk.green(`Successfully registered server '${name}'`));
    });

program
    .command('doctor')
    .description('Check system status')
    .action(() => {
        console.log(chalk.bold('MCPM System Status:'));
        console.log(`- Node.js: ${process.version} ${chalk.green('OK')}`);
        console.log(`- Config: ${configPath} ${chalk.green('OK')}`);
        // Check Docker
        try {
            execSync('docker --version', { stdio: 'ignore' });
            console.log(`- Docker: ${chalk.green('Detected')}`);
        } catch {
            console.log(`- Docker: ${chalk.red('Not Found')}`);
        }
        console.log(chalk.green('\nSystem is healthy.'));
    });

program
    .command('search <query>')
    .description('Search for servers')
    .action((query) => {
        const technologies = getAllServers();
        console.log(chalk.bold(`Search results for "${query}":`));
        for (const group in technologies) {
            for (const key in technologies[group]) {
                if (key.includes(query)) {
                     console.log(`- ${key} (${group})`);
                }
            }
        }
    });

program
    .command('info <name>')
    .description('Get server info')
    .action((name) => {
         const technologies = getAllServers();

         // Prioritize custom/user configuration
         if (technologies['custom'] && technologies['custom'][name]) {
             console.log(JSON.stringify(technologies['custom'][name], null, 2));
             return;
         }

         for (const group in technologies) {
            if (technologies[group][name]) {
                console.log(JSON.stringify(technologies[group][name], null, 2));
                return;
            }
        }
        console.error(chalk.red(`Server ${name} not found`));
    });

program
    .command('uninstall <name>')
    .description('Remove an installed or registered server')
    .action((name) => {
        const userConfig = readUserConfig();
        if (userConfig[name]) {
            delete userConfig[name];
            saveUserConfig(userConfig);
            console.log(chalk.green(`Successfully removed custom server '${name}'`));
            return;
        }

        console.log(chalk.cyan(`Attempting to uninstall '${name}' from dependencies...`));
        try {
             execSync(`npm uninstall ${name}`, { stdio: 'inherit', cwd: __dirname });
             console.log(chalk.green(`Successfully uninstalled '${name}'`));
        } catch (e) {
             console.error(chalk.red(`Failed to uninstall '${name}': ${e.message}`));
             process.exit(1);
        }
    });

// Stubs for other commands to prevent Jarvis errors
['edit', 'usage', 'client', 'profile', 'config', 'migrate', 'share'].forEach(cmd => {
    program.command(cmd).action(() => {
        console.log(chalk.yellow(`Command '${cmd}' is not yet fully implemented in this version.`));
    });
});

program.parse(process.argv);
