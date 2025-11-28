#!/usr/bin/env node

const { Command } = require('commander');
const fs = require('fs');
const path = require('path');
const toml = require('toml');
const { execSync } = require('child_process');
const chalk = require('chalk');

const program = new Command();
const configPath = path.join(__dirname, 'config', 'technologies.toml');

// Helper to read config
function readConfig() {
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

program
    .name('mcpm')
    .description('Model Context Protocol Manager CLI')
    .version('1.0.0');

program
    .command('ls')
    .description('List installed MCP servers')
    .action(() => {
        const config = readConfig();
        console.log(chalk.bold('Installed MCP Servers:'));

        let hasInstalled = false;
        if (config.technologies) {
            for (const group in config.technologies) {
                let groupPrinted = false;
                for (const key in config.technologies[group]) {
                    const tech = config.technologies[group][key];
                    let isInstalled = false;

                    // Check if installed in node_modules
                    // If 'package' is defined, check that. If 'repo' is defined, check the key name?
                    // Simple check: does node_modules/<package_name> exist?
                    let pkgName = tech.package || key; // heuristic

                    // Handle scoped packages
                    if (!tech.package && tech.repo && tech.repo.includes('github.com/')) {
                        // For git repos without explicit package name, we might not know the folder name easily
                        // unless we enforced a convention.
                        // For this dev version, let's just check if the key exists in package.json dependencies
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

                    if (isInstalled) {
                        if (!groupPrinted) {
                            console.log(chalk.blue(`\n[${group}]`));
                            groupPrinted = true;
                        }
                        console.log(`- ${key}: ${tech.description || ''} ${chalk.green('(Installed)')}`);
                        hasInstalled = true;
                    }
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
        const config = readConfig();
        let found = false;
        let pkgName = name;
        let dockerImage = null;

        // Simple lookup strategy
        for (const group in config.technologies) {
             if (config.technologies[group][name]) {
                 found = true;
                 const tech = config.technologies[group][name];
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
             console.log(chalk.yellow(`Warning: ${name} not found in technologies.toml. Attempting direct npm install...`));
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
            // For npm, the target is the package name (or mapped name)
            // We need to know the actual folder in node_modules. usually pkgName unless it's a git url.
            // For git urls, npm usually installs to the repo name.
            // For this prototype, we'll assume the user provided 'name' matches the folder or we'd need to parse pkgName.
            let folderName = pkgName;
            if (pkgName.startsWith('git+')) {
                // simplistic fallback
                folderName = name;
            }
            printConfigSnippet(name, 'npm', folderName);
        } catch (e) {
            console.error(chalk.red(`Failed to install ${name}: ${e.message}`));
            process.exit(1);
        }
    });

program
    .command('doctor')
    .description('Check system status')
    .action(() => {
        console.log(chalk.bold.blue('🔍 Checking System Vital Signs...'));
        console.log(chalk.gray('----------------------------------------'));

        let allGood = true;

        // 1. Check Node.js
        if (process.version) {
            console.log(`✅ Node.js Runtime    : ${chalk.green(process.version)}`);
        } else {
            console.log(`❌ Node.js Runtime    : ${chalk.red('Error')}`);
            allGood = false;
        }

        // 2. Check Config
        if (fs.existsSync(configPath)) {
             console.log(`✅ Registry File      : ${chalk.green('Connected')}`);
        } else {
             console.log(`❌ Registry File      : ${chalk.red('Missing')}`);
             allGood = false;
        }

        // 3. Check Docker
        try {
            const dockerVer = execSync('docker --version', { encoding: 'utf8' }).trim();
            console.log(`✅ Docker Engine      : ${chalk.green(dockerVer)}`);
        } catch {
            console.log(`⚠️ Docker Engine      : ${chalk.yellow('Not Found (Remote/Docker-based tools will fail)')}`);
            // Not fatal for everything, but good to know
        }

        // 4. Check Jarvis Binary
        const jarvisPath = path.join(__dirname, '..', 'Jarvis', 'jarvis');
        if (fs.existsSync(jarvisPath)) {
             console.log(`✅ Jarvis Core        : ${chalk.green('Online')}`);
        } else {
             console.log(`❌ Jarvis Core        : ${chalk.red('Offline (Binary not found)')}`);
             allGood = false;
        }

        console.log(chalk.gray('----------------------------------------'));

        if (allGood) {
            console.log(chalk.bold.green('\n🚀 ALL SYSTEMS GO! 🚀\n'));
            console.log(chalk.cyan('Jarvis is ready to assist.'));
        } else {
            console.log(chalk.bold.red('\n⚠️  SYSTEM CHECK FAILED'));
            console.log('Please resolve the issues above.');
            process.exit(1);
        }
    });

program
    .command('search <query>')
    .description('Search for servers')
    .action((query) => {
        const config = readConfig();
        console.log(chalk.bold(`Search results for "${query}":`));
        for (const group in config.technologies) {
            for (const key in config.technologies[group]) {
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
         const config = readConfig();
         for (const group in config.technologies) {
            if (config.technologies[group][name]) {
                console.log(JSON.stringify(config.technologies[group][name], null, 2));
                return;
            }
        }
        console.error(chalk.red(`Server ${name} not found`));
    });

// Stubs for other commands to prevent Jarvis errors
['uninstall', 'edit', 'new', 'usage', 'client', 'profile', 'config', 'migrate', 'share'].forEach(cmd => {
    program.command(cmd).action(() => {
        console.log(chalk.yellow(`Command '${cmd}' is not yet fully implemented in this version.`));
    });
});

program.parse(process.argv);
