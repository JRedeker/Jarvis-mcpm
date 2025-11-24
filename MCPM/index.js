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
        // In a real implementation, this would check node_modules or a registry file
        // For now, we'll list what's in package.json dependencies if possible, or just mock it
        // matching the technologies.toml
        const config = readConfig();
        console.log(chalk.bold('Available Technologies (from config):'));
        if (config.technologies) {
            for (const group in config.technologies) {
                console.log(chalk.blue(`\n[${group}]`));
                for (const key in config.technologies[group]) {
                    const tech = config.technologies[group][key];
                    console.log(`- ${key}: ${tech.description || ''}`);
                }
            }
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
        
        // Find package name from mapping or config
        if (config.package_mappings) {
             // Reverse lookup or direct? The toml has "@pkg" = "name"
             // We want to find the tech entry for "name"
        }
        
        // Simple lookup strategy
        for (const group in config.technologies) {
             if (config.technologies[group][name]) {
                 found = true;
                 const tech = config.technologies[group][name];
                 if (tech.package) {
                     pkgName = tech.package;
                 } else if (tech.repo) {
                     // For git repos, we might strictly need a package name or install from git
                     console.log(chalk.yellow(`Note: ${name} is a repository-based tool. Installing from source not fully implemented in this skeleton.`));
                     return; 
                 }
             }
        }

        if (!found) {
             console.log(chalk.yellow(`Warning: ${name} not found in technologies.toml. Attempting direct npm install...`));
        }

        try {
            execSync(`npm install ${pkgName}`, { stdio: 'inherit', cwd: __dirname });
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
