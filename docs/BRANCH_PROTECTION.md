# Branch Protection Configuration

This document describes the recommended branch protection settings for the Jarvis repository.

## Main Branch Protection

Configure these settings in GitHub: `Settings → Branches → Add rule`

### Basic Settings

**Branch name pattern:** `main`

### Protect Matching Branches

| Setting | Value | Reason |
|:--------|:------|:-------|
| **Require a pull request before merging** | ✅ Enabled | All changes via PR |
| **Require approvals** | 1 | At least one reviewer |
| **Dismiss stale pull request approvals** | ✅ Enabled | New commits require re-review |
| **Require review from Code Owners** | Optional | If CODEOWNERS file exists |

### Status Checks

**Require status checks to pass before merging:** ✅ Enabled

**Required checks:**
- `Go Tests` - Core handler tests must pass
- `Go Build` - Binary must compile
- `Shell Tests (bats)` - Script tests must pass

**Optional checks (continue-on-error):**
- `Go Lint` - Linting warnings don't block
- `Docker Build` - Docker image builds

### Additional Settings

| Setting | Value |
|:--------|:------|
| **Require branches to be up to date** | ✅ Enabled |
| **Require conversation resolution** | ✅ Enabled |
| **Require signed commits** | Optional |
| **Include administrators** | ✅ Enabled |
| **Allow force pushes** | ❌ Disabled |
| **Allow deletions** | ❌ Disabled |

---

## Develop Branch (Optional)

For a `develop` branch workflow:

**Branch name pattern:** `develop`

| Setting | Value |
|:--------|:------|
| Require a pull request | ✅ Enabled |
| Require approvals | 0 (self-merge OK) |
| Require status checks | ✅ Enabled (same as main) |
| Allow force pushes | ❌ Disabled |

---

## Feature Branches

Feature branches (`feat/*`, `fix/*`, etc.) should not have protection rules to allow developers to work freely.

---

## CI/CD Status Badges

The README includes these badges that reflect CI status:

```markdown
[![Tests](https://img.shields.io/github/actions/workflow/status/JRedeker/Jarvis-mcpm/test.yml?branch=main&label=tests&logo=github)](https://github.com/JRedeker/Jarvis-mcpm/actions/workflows/test.yml)
[![Coverage](https://img.shields.io/badge/coverage-60%25-yellow?logo=codecov)](Jarvis/)
[![Go Report Card](https://goreportcard.com/badge/github.com/JRedeker/Jarvis-mcpm)](https://goreportcard.com/report/github.com/JRedeker/Jarvis-mcpm)
```

---

## Recommended Workflow

1. Create feature branch from `main`:
   ```bash
   git checkout -b feat/my-feature main
   ```

2. Make changes, commit with conventional format:
   ```bash
   git commit -m "feat: Add new handler"
   ```

3. Push and create PR:
   ```bash
   git push -u origin feat/my-feature
   gh pr create
   ```

4. Wait for CI checks to pass

5. Request review, address feedback

6. Squash and merge to `main`

---

## Troubleshooting CI Failures

### Go Tests Failing

```bash
cd Jarvis
go test -v ./...
```

### Go Build Failing

```bash
cd Jarvis
go mod tidy
go build -v .
```

### Shell Tests Failing

```bash
# Install bats locally
npm install -g bats
# Run tests
bats scripts/tests/*.bats
```

### Docker Build Failing

```bash
docker build -t mcpm-daemon ./mcpm-daemon
```
