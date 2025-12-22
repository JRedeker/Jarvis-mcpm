# Implementation Tasks

1.  [ ] **Create New Profiles**
    *   Use `jarvis_profile` to create `essentials`, `dev-core`, `research`.
    *   Assign appropriate tools to each (migrating from `toolbox`).
    *   Validation: `jarvis_server(action="list")` shows correct distribution.

2.  [ ] **Configure Port Mappings**
    *   Update `MCPM/api/helpers.js` or config to ensure `research` gets port 6281 (or assign manually via `jarvis_profile`).
    *   Ensure `dev-core` uses 6278.

3.  [ ] **Update Client Configuration (OpenCode)**
    *   Use `jarvis_client(action="edit")` to set the profile stack: `essentials,memory,dev-core,research,data`.
    *   Validation: Check `opencode.json`.

4.  [ ] **Update Documentation**
    *   Update `docs/CONFIGURATION_STRATEGY.md` with the new architecture diagram/description.
    *   Update `AGENTS.md` (OpenSpec header) to reflect the new profile structure in instructions.
    *   Update `README.md` if it references the old `toolbox`.

5.  [ ] **Verification**
    *   Restart all profiles.
    *   Run `jarvis_diagnose(action="full")` to confirm all 5 profiles are running and healthy.
    *   Test one tool from each profile (e.g., `time`, `brave-search`, `basic-memory`).

6.  [ ] **Cleanup**
    *   Delete the legacy `toolbox` profile once migration is confirmed.
