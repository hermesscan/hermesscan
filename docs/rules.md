# HermesScan Rule Reference

Generated for HermesScan 0.8.0.

| ID | Severity | Category | Name |
|---|---|---|---|
| `HMS0001` | Medium | reliability | Sleep-based synchronization |
| `HMS0002` | High | isolation | Fixed PostgreSQL port exposure |
| `HMS0003` | High | orchestration | Manual shell background orchestration |
| `HMS0004` | Medium | orchestration | Manual wait orchestration |
| `HMS0005` | High | orchestration | PowerShell background job in CI path |
| `HMS0006` | Medium | lifecycle | PowerShell process lifecycle risk |
| `HMS0007` | Medium | isolation | Docker Compose service startup |
| `HMS0008` | Medium | isolation | Shared temp path |
| `HMS0009` | Low | supply-chain | Mutable GitHub Action reference |
| `HMS0010` | Low | cache | Package install cache collision risk |
| `HMS0011` | Medium | cache | Docker build cache contention risk |
| `HMS0012` | Medium | reliability | Native command exit code may be ignored |
| `HMS0013` | High | supply-chain | GitHub pull_request_target trigger |
| `HMS0014` | Medium | supply-chain | GitHub Actions write-all permissions |
| `HMS0015` | Medium | isolation | GitHub Actions self-hosted runner |
| `HMS0016` | Medium | cache | GitHub Actions broad cache key |

## HMS0001 - Sleep-based synchronization

**Severity:** Medium  
**Category:** reliability  
**Tags:** `synchronization`, `flake-risk`  
**File types:** `bash`, `powershell`, `yaml`, `makefile`  

Sleep-based synchronization can make CI jobs flaky or unnecessarily slow.

**Recommendation:** Prefer explicit readiness checks, health probes, retry loops with deadlines, or CI job dependencies.

```text
(?i)\b(sleep\s+[0-9]+|Start-Sleep\s+(-Seconds\s+)?[0-9]+)
```

## HMS0002 - Fixed PostgreSQL port exposure

**Severity:** High  
**Category:** isolation  
**Tags:** `ports`, `shared-runner`, `postgresql`  
**File types:** `bash`, `powershell`, `yaml`, `docker`, `makefile`  

Fixed PostgreSQL port exposure can collide on shared runners or concurrent local test runs.

**Recommendation:** Use isolated containers, dynamic ports, unique project names, or per-job network namespaces; avoid binding host port 5432 in shared CI.

```text
(?i)(ports?:|--publish|-p|PGPORT|POSTGRES_PORT|DATABASE_URL|localhost:|127\.0\.0\.1:|0\.0\.0\.0:).{0,80}\b5432\b|\b5432\s*:\s*5432\b|\b5432/tcp\b
```

## HMS0003 - Manual shell background orchestration

**Severity:** High  
**Category:** orchestration  
**Tags:** `parallelism`, `shell`  
**File types:** `bash`, `makefile`  

Background jobs inside CI scripts may indicate manual orchestration that is difficult to observe and clean up.

**Recommendation:** Move parallel work into the CI orchestrator DAG or add explicit lifecycle management and cleanup.

```text
(^|[^&])&\s*($|#)
```

## HMS0004 - Manual wait orchestration

**Severity:** Medium  
**Category:** orchestration  
**Tags:** `parallelism`, `shell`  
**File types:** `bash`, `makefile`  

Manual wait logic can hide job dependencies from the CI control plane.

**Recommendation:** Represent dependencies as CI jobs or document why process-level orchestration is required.

```text
(?i)\bwait\b
```

## HMS0005 - PowerShell background job in CI path

**Severity:** High  
**Category:** orchestration  
**Tags:** `powershell`, `background-job`  
**File types:** `powershell`  

PowerShell background jobs can outlive the intended script scope and complicate CI cleanup.

**Recommendation:** Prefer CI jobs, use try/finally cleanup, and explicitly stop/remove jobs when background work is necessary.

```text
(?i)\bStart-Job\b
```

## HMS0006 - PowerShell process lifecycle risk

**Severity:** Medium  
**Category:** lifecycle  
**Tags:** `powershell`, `process`  
**File types:** `powershell`  

Started processes can leak across build steps if they are not tracked and cleaned up.

**Recommendation:** Capture the process object, enforce timeouts, and clean up in a finally block.

```text
(?i)\bStart-Process\b
```

## HMS0007 - Docker Compose service startup

**Severity:** Medium  
**Category:** isolation  
**Tags:** `docker`, `compose`  
**File types:** `bash`, `powershell`, `yaml`, `makefile`  

Docker Compose can collide across concurrent jobs when project names, networks, volumes, or ports are shared.

**Recommendation:** Set a unique COMPOSE_PROJECT_NAME per CI job and ensure compose down runs during cleanup.

```text
(?i)docker\s+compose\s+up|docker-compose\s+up
```

## HMS0008 - Shared temp path

**Severity:** Medium  
**Category:** isolation  
**Tags:** `filesystem`, `temp`  
**File types:** `bash`, `powershell`, `yaml`, `makefile`  

Shared temporary paths can cause cross-job contamination on shared runners.

**Recommendation:** Use a unique per-job workspace or temporary directory derived from the CI run ID.

```text
(?i)(/tmp/|C:\\Temp|\$env:TEMP|\$env:TMP)
```

## HMS0009 - Mutable GitHub Action reference

**Severity:** Low  
**Category:** supply-chain  
**Tags:** `github-actions`, `pinning`  
**File types:** `yaml`  

Mutable action references can change without review and reduce workflow reproducibility.

**Recommendation:** Pin GitHub Actions to a specific version tag or commit SHA.

```text
uses:\s+[^@\s]+@(?:main|master|latest)
```

## HMS0010 - Package install cache collision risk

**Severity:** Low  
**Category:** cache  
**Tags:** `package-manager`  
**File types:** `bash`, `powershell`, `yaml`, `makefile`  

Package managers may use shared caches that can behave poorly under concurrent jobs. This rule is advisory unless paired with shared cache paths or manual parallelism.

**Recommendation:** Use CI-managed caches with scoped keys and avoid shared writable cache directories across concurrent jobs.

```text
(?i)\b(npm\s+install|npm\s+ci|yarn\s+install|pip\s+install|dotnet\s+restore|mvn\s+dependency|mvn\s+install)\b
```

## HMS0011 - Docker build cache contention risk

**Severity:** Medium  
**Category:** cache  
**Tags:** `docker`, `buildkit`  
**File types:** `bash`, `powershell`, `yaml`, `makefile`  

Concurrent Docker builds may contend over build cache or daemon state on shared runners.

**Recommendation:** Use isolated builders, unique cache scopes, or CI-provided build isolation.

```text
(?i)docker\s+build|docker\s+buildx\s+build
```

## HMS0012 - Native command exit code may be ignored

**Severity:** Medium  
**Category:** reliability  
**Tags:** `powershell`, `native-command`  
**File types:** `powershell`  

Windows PowerShell 5.1 does not automatically throw when native commands return non-zero exit codes.

**Recommendation:** Check $LASTEXITCODE after native commands or wrap native invocations in a helper that throws on failure.

```text
(?i)^\s*(docker|npm|yarn|pip|dotnet|mvn|git)\s+
```

## HMS0013 - GitHub pull_request_target trigger

**Severity:** High  
**Category:** supply-chain  
**Tags:** `github-actions`, `pull-request`, `permissions`  
**File types:** `yaml`  

The pull_request_target trigger can expose privileged repository context to pull-request workflows when used incorrectly.

**Recommendation:** Use pull_request for untrusted code or isolate pull_request_target jobs so they do not check out or execute attacker-controlled content.

```text
(?i)\bpull_request_target\b
```

## HMS0014 - GitHub Actions write-all permissions

**Severity:** Medium  
**Category:** supply-chain  
**Tags:** `github-actions`, `permissions`  
**File types:** `yaml`  

Broad write-all workflow permissions increase blast radius if a job or action is compromised.

**Recommendation:** Grant the minimum required permissions explicitly, such as contents: read or pull-requests: write.

```text
(?i)permissions:\s*write-all
```

## HMS0015 - GitHub Actions self-hosted runner

**Severity:** Medium  
**Category:** isolation  
**Tags:** `github-actions`, `runner`, `self-hosted`  
**File types:** `yaml`  

Self-hosted runners can share filesystem, Docker, network, and cache state across jobs when not isolated carefully.

**Recommendation:** Use ephemeral self-hosted runners where possible, enforce per-job cleanup, and avoid shared writable state between jobs.

```text
(?i)runs-on:\s*(\[[^\]]*)?self-hosted\b
```

## HMS0016 - GitHub Actions broad cache key

**Severity:** Medium  
**Category:** cache  
**Tags:** `github-actions`, `cache`  
**File types:** `yaml`  

Broad cache keys can cause unrelated branches or dependency states to reuse the same writable cache namespace.

**Recommendation:** Include dependency lockfile hashes or job-specific inputs in cache keys, such as hashFiles('**/package-lock.json').

```text
(?i)key:\s*\$\{\{\s*runner\.os\s*\}\}\s*$
```

