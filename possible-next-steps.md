# Possible next steps to recover Cipher (cipher-aggregator)

A prioritized list of next steps to fix the cipher-aggregator (Cipher MCP) based on recent logs and attempts.

## Priority A — Essential fixes to get Cipher running

1. **Correct the `cipher.yml` agent config and validate it**
   - _Why:_ Cipher currently fails to load the agent config (`invalid_type expected object received undefined`). The file must be top-level YAML with the `agent:` block properly indented.
   - _Action:_
     - Backup the file (`cp cipher.yml cipher.yml.bak`).
     - Ensure the beginning of `cipher.yml` contains the following block (column 0):
       ```yaml
       agent:
         provider: "openrouter"
         model: "gpt-4.1-mini"
         api_key: "${OPENROUTER_API_KEY}"
         temperature: 0.7
         max_tokens: 4096
       ```
     - Keep the rest of the YAML structure as originally intended (systemPrompt, toolExecution, servers, etc.) without adding leading spaces before the top-level keys.
     - If you need to test locally without a real credential, you may replace `${OPENROUTER_API_KEY}` with `dummy_key_for_dev`, but make sure the YAML remains well-formed (no indentation errors).
   - _Validation:_ Run `python3 -c 'import yaml, json; print(json.dumps(yaml.safe_load(open("cipher.yml")), indent=2))'` to ensure it parses.

2. **Ensure environment variables are loaded into cipher process**
   - _Why:_ Cipher complains about missing API keys when starting; `mcp-manager.sh` now exports `.env`, but the proc still needs to see `OPENROUTER_API_KEY`.
   - _Action:_ `source .env && echo "$OPENROUTER_API_KEY"` should print the expected key.

3. **Restart cipher via `mcp-manager.sh` and monitor logs**
   - _Action:_
     - `./mcp-manager.sh restart`. If it fails, capture the first 30 lines of the new log (`tail -n 60 logs/cipher-aggregator-*.log`).
     - `tail -f logs/cipher-aggregator-*.log /tmp/cipher-mcp.log`
   - _Verification:_ Look for lines such as "Loading agent config from /home/jrede/dev/MCP/cipher.yml" followed by a series of successful tool registrations and the SSE server ready message.

4. **Confirm cipher is running**
   - _Action:_ `./mcp-manager.sh status` report confirms running process and SSE endpoint. You may open `http://127.0.0.1:3020/sse` if desired.

## Priority B — Stability improvements

5. **Fix or reinstall better-sqlite3 native bindings**
   - _Why:_ Logs show "Could not locate the bindings file..." which means sqlite persistence falls back to in-memory storage.
   - _Action:_ Rebuild the native bindings used by the global pnpm/npm install (e.g., `pnpm rebuild better-sqlite3` or reinstall `better-sqlite3` in the environment that supplies cipher).

6. **Install or disable optional MCP servers throwing npm E404**
   - _Why:_ Log noise from `magic-mcp` and `server-web` failing to install (npm 404). They are optional but may clutter logs.
   - _Action:_ Either disable them in `cipher.yml` (set `enabled: false`) or install from a valid source (if you have local packages).

## Phase 1 validation commands (safe to run now)
- `cp cipher.yml cipher.yml.bak`
- `python3 -c 'import yaml, sys; print(yaml.safe_load(open("cipher.yml")))'`
- `source .env && echo "$OPENROUTER_API_KEY"`
- `python3 -c 'import yaml, json; print(json.dumps(yaml.safe_load(open("cipher.yml")), indent=2))'`

## What success looks like
- No YAML parsing errors when loading `cipher.yml`.
- Cipher starts successfully (`./mcp-manager.sh restart`) and stays running.
- `/tmp/cipher-mcp.log` no longer shows "invalid_type" or missing API key errors.
- `cipher.yml` remains valid YAML and `api_key` resolves to a real credential.

Feel free to run these steps now or let me know if you’d like me to take the edits in Act mode.
