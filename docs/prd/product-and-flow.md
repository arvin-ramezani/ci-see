# CI See — Product intent and explicit CLI flow

**Status:** Draft revision — owner-requested CLI direction
**Canonical authority:** [Product requirements](../prd.md)
**Last updated:** 2026-10-13

## 1. Product Summary

CI See is a lightweight Go CLI that runs existing GitHub Actions workflows locally
through `act`. Developers or AI agents explicitly invoke it, receive clear results,
and reuse downloaded dependencies to reduce setup time and hosted CI dependence.

## 2. Problem

Push-and-wait validation adds latency and can exhaust included hosted Actions
minutes, particularly for private repositories on GitHub Free. Standard
GitHub-hosted runners are free for public repositories; faster feedback remains useful there.

`act` already executes workflows locally. CI See simplifies setup diagnostics,
workflow/event selection, dependency reuse, and human/agent result reporting.
It must not reproduce the workflow engine or introduce a second CI definition.

Local execution does not prevent GitHub workflows from running after a push or PR.
The MVP cannot promise elimination of hosted usage while remote triggers remain enabled.

## 3. Target Users

- Solo developers and small teams with GitHub Actions workflows.
- GitHub Free users seeking to reduce hosted CI usage.
- Developers using Codex, Claude, Cursor, or other coding agents.
- Contributors validating compatible workflows before PRs.

## 4. Product Principles

1. **One CI definition:** use existing GitHub Actions YAML.
2. **Explicit invocation:** the developer owns when to run verification.
3. **Small product:** Go CLI around `act`; no background daemon or approval GUI.
4. **Cache reuse:** download missing dependencies and reuse verified compatible copies.
5. **Honest results:** unknown, unsupported, interrupted, incomplete, or failed validation is not PASS.
6. **Local context:** PASS describes the selected local run, not hosted-runner parity.
7. **Cross-platform hosts:** consistent commands on Windows, macOS, and Linux.
8. **Low adoption cost:** minimal configuration and no mandatory workflow rewrites.
9. **Lightweight installation:** prebuilt binaries need no host Go/Node runtime
   for CI See itself. Workflow steps may require their own runtimes.

## 5. Recommended User Flow

1. Install CI See, `act`, and Docker; initialize only if local configuration is needed.
2. Change code and explicitly invoke `ci-see` from the repository.
3. Select a workflow/job/event when needed; acquire missing dependencies and run locally.
4. Inspect the result/logs and fix failures or decide the next development action.
5. Commit, push, or create a PR using ordinary Git/GitHub commands.

There is no automatic run, commit/push interception, or protection against forgetting
verification. The CLI must not commit or push on the user's behalf.

### Commit gate

**Removed.** CI See neither installs commit hooks nor permits/blocks commits.
This heading remains for legacy links; it is not an implementation requirement.

### AI-agent flow

An optional small skill can tell an agent to invoke CI See as a verification step,
inspect the exit code and JSON report, and report PASS, FAIL, or INCOMPLETE.
The same CLI works without a skill or agent-specific service.

Agents must report actual selected scope and limitations. Choosing not to run
verification cannot be represented as a successful CI See run.

### Push gate

**Removed.** CI See does not intercept pushes or require a result/bypass to push.
Existing GitHub workflow triggers remain unchanged.
