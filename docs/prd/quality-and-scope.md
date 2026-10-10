# Product requirements — Reliability, compatibility, and MVP scope

**Status:** Draft (inherited)  
**Canonical authority:** [Product requirements](../prd.md)

<!-- BEGIN ORIGINAL CONTRACT -->
## 7. Reliability and Safety Requirements

1. **Fail closed:** if CI See cannot prove PASS for the exact state, it is not PASS.
2. **No false PASS after crash:** interrupted or partial execution cannot become PASS.
3. **Concurrency safety:** simultaneous CI runs or Git operations cannot corrupt state or associate a result with the wrong code.
4. **Repository isolation:** state from one repository/worktree cannot satisfy another incompatible state.
5. **Workflow invalidation:** workflow/config changes invalidate dependent results.
6. **Durable bypass audit:** every allowed bypass requires a successfully persisted local audit record before continuation; missing or failed audit persistence blocks Git.
7. **Safe recovery:** if CI See crashes while Git is waiting, Git fails safely instead of silently continuing.

## 8. Compatibility Goals

CI See should maximize compatibility with ordinary GitHub Actions by delegating execution to an established engine rather than inventing a new workflow language.

Initial focus:

- ubuntu-latest style jobs;
- common marketplace actions;
- service containers;
- environment variables and secrets;
- job dependencies;
- matrices where supported;
- common GitHub expressions and events.

Full GitHub-hosted runner parity is not an MVP promise.

## 9. MVP Non-Goals

The MVP will not:

- implement a GitHub Actions engine from scratch;
- guarantee exact GitHub-hosted runner-image parity;
- replace GitHub;
- provide hosted CI SaaS;
- fully emulate macOS or Windows hosted runners;
- let AI agents automatically approve failed CI;
- require a second CI configuration;
- become a general Docker orchestration tool.

## 10. MVP Scope

The first usable release must prove the complete primary loop:

1. install CI See;
2. initialize a repository;
3. discover GitHub Actions workflows;
4. run relevant workflows locally;
5. persist exact-state results;
6. automatically gate commit;
7. automatically gate push;
8. let the developer Cancel or Continue anyway after non-PASS state;
9. support Git initiated by AI-agent/IDE processes;
10. provide consistent supported behavior on Windows, Linux, and macOS;
11. provide enough diagnostics to understand failures.

<!-- END ORIGINAL CONTRACT -->
