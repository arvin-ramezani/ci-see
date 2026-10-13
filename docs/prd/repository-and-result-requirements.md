# Product requirements — Setup, execution, dependencies, results

**Status:** Draft revision — owner-requested CLI direction
**Canonical authority:** [Product requirements](../prd.md)

## 6. Functional Requirements

### FR-1 — Repository initialization

`ci-see init` detects the repository, workflows, and execution dependencies.
Create only minimal configuration when necessary. Default execution should not
require initialization or a second pipeline definition.

Do not install hooks, alter Git configuration, rewrite workflows, automatically
install host tools, or change firewall/Docker-wide settings.
Missing tools produce actionable installation guidance.

### FR-2 — Workflow discovery

Discover `.github/workflows/` automatically without manual registration.
Provide workflow/job listing and selection. Report invalid YAML, no matching
workflows, and unsupported selections; empty selection is not successful validation.

### FR-3 — One-command local CI

`ci-see` explicitly runs local validation in the current repository. CLI options
select workflows, jobs, events, and inputs without changing workflow files.

Use a compiled Go CLI, `nektos/act` as the MVP engine, and Docker with Linux containers.
CI See adds orchestration, setup diagnostics, and results; it does not build another
Actions interpreter, cache engine, or container runtime.

#### Dependency acquisition and reuse

- Normal use may access the network to acquire missing actions, runner images,
  and workflow package dependencies. The product is not permanently offline.
- Reuse compatible cached actions/images through `act` and Docker.
  Package caching follows the workflow's package manager and configured caches;
  CI See does not promise automatic caching of every package.
- Verify actions against trusted resolved commits and images against trusted
  platform-specific registry digests before execution/reuse; record identities.
  A cache directory, mutable tag, or local image ID alone is insufficient.
- Resolve mutable references through the configured trusted source. A changed
  commit/digest is a new dependency; never silently substitute versions.
- Respect package-manager lockfiles/integrity checks where available. If the
  selected trust policy cannot verify a required dependency, stop as INCOMPLETE.
- Failed downloads, corrupt caches, failed integrity checks, or unavailable
  required networking produce INCOMPLETE, never PASS.
- Scope acquisition credentials to their task. Workflow steps may need runtime
  networking; cache reuse alone does not guarantee offline execution.
- An optional cache-only/offline policy reports missing dependencies. Offline
  flags alone are not network-security boundaries.

### FR-4 — Event handling

Support `push`, `pull_request`, and `workflow_dispatch` for compatible workflows.
Document the default event and expose explicit selection. Derive local metadata
and accept CLI inputs for missing fields; normal use should not need manual JSON.
Advanced callers may supply an event payload.

Required unresolved inputs produce INCOMPLETE. Local simulation must not claim
that all GitHub-generated context/event values are reproduced.

### FR-5 — Durable result identity

Provide human-readable results, stable process exit semantics, and optional JSON
for agents. Persist a local report containing:

- run and repository/worktree identity;
- actual tested source identity, including local changes when present;
- selected workflows/jobs, event and relevant non-secret inputs;
- start/end timestamps and overall/per-workflow/job outcomes;
- CI See, `act`, and runner-image versions/identities;
- configuration and verified dependency identities;
- skipped/excluded work, compatibility warnings, and non-PASS reasons.

Default validation uses current local files, including uncommitted changes,
rather than exclusively the staged index. A commit SHA alone must not describe
a dirty tree as though it were the tested source. The execution spec must define
reliable tested-input capture and identification.

Required outcomes: **PASS**, **FAIL**, **INCOMPLETE** (including cancellation).
No run is **NOT_RUN**, never PASS. Historical reports describe earlier runs;
they are not automatically reusable validation for changed inputs.
Detailed exit codes and JSON schema belong in the execution specification.
