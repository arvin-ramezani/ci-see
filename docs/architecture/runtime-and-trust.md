# System architecture — Execution/gating flows, approvals, security, platforms

**Status:** Draft (inherited)  
**Canonical authority:** [System architecture](../architecture.md)

<!-- BEGIN ORIGINAL CONTRACT -->
## 7. Runtime Flows

### Local run

```text
ci-see
→ inspect repository/workflows
→ determine validation identity
→ record RUNNING
→ execute act
→ persist PASS / FAIL / INCOMPLETE
→ return a clear exit result
```

### Commit gate

```text
git commit
→ hook invokes CI See
→ calculate staged-state identity
→ matching PASS: continue
→ otherwise: block and explain
```

### Push gate

```text
git push
→ hook invokes CI See
→ identify pushed commit state
→ matching acceptable result: continue
→ otherwise: block
```

## 8. Approval Boundary

Approval remains a separate core boundary.

The MVP does not require a long-lived daemon or desktop application, but the CLI may invoke a short-lived human approval surface when a developer decision is required. The exact mechanism belongs in the approval UX spec.

CI See must not expose a trivial machine-only `--force` or equivalent bypass that an AI agent could silently use. Non-interactive Git operations must fail safely if a valid developer approval cannot be obtained.

This keeps the core CLI architecture small while preserving the PRD requirement for developer-controlled bypass across terminal, IDE, and AI-agent initiated Git operations.

## 9. Reliability and Security

- Unknown, stale, corrupt, or interrupted state is not PASS.
- Invoke `git` and `act` via Go `os/exec` with argument arrays and bounded cancellation; avoid shell-string construction.
- Validate external/process/file data at boundaries.
- Do not persist secrets in result metadata.
- Logs must explain why a gate passed, failed, or became stale.
- Local `act` success is evidence, not proof of exact GitHub-hosted runner parity.

## 10. Platform Strategy

Primary development is Windows 11 + WSL2. Use a Linux binary inside WSL2 and a Windows binary for native Windows terminals; do not assume a binary for one host can run on the other.

Publish Go CLI builds for each supported OS/architecture; installing a prebuilt release must not require Go, Node.js, or npm. Keep the Go version declared in `go.mod` and pin it in build/release workflows when implementation begins.

Architecture must avoid WSL-specific assumptions. Platform-specific behavior stays behind adapters so native Windows, Linux, and macOS can be validated independently.

The MVP focuses on Linux-container execution through `act`; full Windows/macOS hosted-runner emulation is not required.

<!-- END ORIGINAL CONTRACT -->
