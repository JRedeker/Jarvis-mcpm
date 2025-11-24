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
        
        // Simple lookup strategy
        for (const group in config.technologies) {
             if (config.technologies[group][name]) {
                 found = true;
                 const tech = config.technologies[group][name];
                 if (tech.package) {
                     pkgName = tech.package;
                 } else if (tech.repo) {
                     // If it's a repo, we might want to install via git+https
                     pkgName = `git+${tech.repo}.git`;
                 }
             }
        }

        // Check package mappings if not found in tech tree but might be a mapped name
        if (!found && config.package_mappings) {
             // This logic is a bit circular. Usually mappings map pkg -> shortname.
             // Here we are looking up shortname -> pkg.
             // The technologies.toml structure handles shortname -> data.
        }

        if (!found) {
             console.log(chalk.yellow(`Warning: ${name} not found in technologies.toml. Attempting direct npm install...`));
        }

        try {
            // Use --save to ensure it's added to package.json
            execSync(`npm install ${pkgName} --save`, { stdio: 'inherit', cwd: __dirname });
            console.log(chalk.green(`Successfully installed ${name}`));
        } catch (e) {
            console.error(chalk.red(`Failed to install ${name}: ${e.message}`));
            process.exit(1);
        }
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
