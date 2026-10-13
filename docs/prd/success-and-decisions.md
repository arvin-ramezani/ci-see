# Product requirements — Success, decisions, validation, open questions

**Status:** Draft revision — owner-requested CLI direction
**Canonical authority:** [Product requirements](../prd.md)

## 11. Success Criteria

A developer runs compatible existing workflows with one explicit command,
understands the result, and continues normal Git independently.

Acceptance must demonstrate:

- no required YAML rewrites, hooks, approval UI, or background service;
- cold-cache acquisition and warm-cache reuse without redundant compatible downloads;
- missing tools, corrupt dependencies, unsupported selected jobs, job failures,
  and interruptions yield clear non-PASS results;
- selection, inputs, logs, exits, and JSON work non-interactively for humans/agents;
- reports identify actual tested local inputs, dirty changes, and selected scope;
- consistent core CLI behavior across the published Windows/macOS/Linux matrix.

Measure setup steps, orchestration overhead, cold/warm duration, disk use, and
downloaded bytes on a documented representative workflow. Set numeric budgets
in the revised plan after a baseline. Do not invent hosted-minute savings while
remote triggers remain enabled.

## 12. Confirmed Product Decisions

Owner-requested revision, 2026-10-13:

- **CI See** remains a Go CLI with prebuilt binaries; no host Go/Node dependency
  for CI See itself. Building from source requires Go.
- GitHub Actions YAML remains the CI definition.
- **Explicit CLI invocation** replaces automatic execution and all Git gates.
- Developers own when to run it; agents may invoke the same CLI as verification.
- `act` is the MVP engine; Docker/Linux containers are the execution substrate.
- Download missing dependencies, verify them, and reuse existing caching facilities;
  normal operation is not permanently offline.
- Windows, macOS, and Linux hosts are required within a tested support matrix.
- Hosted suppression and broader OS/native runner support are deferred.

## 13. Current MVP Direction Requiring Technical Validation

- Validate `act` against representative workflows and supported host platforms.
- Select engine versions and runner images balancing size and compatibility.
- Verify acquisition/cache integrity, package-manager boundaries, explicit secrets,
  cancellation, concurrency, and reliable result capture.
- Publish tested host/architecture and workflow-feature matrices.

This revision is not runtime acceptance. Existing D3 remains PARTIAL; live-`act`
testing retains its separate security authorization. No firewall change, host
provisioning, or workflow execution is authorized by this documentation change.

## 14. Open Questions for Later Specs

1. Default event and workflow/job selection policy.
2. Derived event fields versus required explicit inputs.
3. Supported engine versions, host architectures, and runner images.
4. Capturing/fingerprinting actual tested files and reporting input changes.
5. Stable exit codes and versioned JSON/report schema.
6. Existing `act` configuration, cache locations, and mutable-reference handling.
7. Cross-platform secret input methods.
8. Network/trust policies for acquisition versus execution.
9. Cancellation, cleanup, and concurrency bounds.
10. Feasibility of future opt-in hosted suppression.
