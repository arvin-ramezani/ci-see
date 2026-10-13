# Local CI execution — Explicit invocation, discovery and event selection

**Status:** Draft aligned with merged PRD #20; command details pending
**Canonical authority:** [Local CI execution](../local-ci-execution.md)

## 1. Purpose

Specify explicit local execution of existing workflows through `act` and Docker.
CI See orchestrates; it does not interpret Actions itself or enforce Git operations.

## 2. Scope

Own discovery/selection, inputs, prerequisites, tested source, engine lifecycle,
dependencies, secrets, reports and compatibility limits.
Exclude hooks, approval/bypass, automatic runs and hosted-run suppression.

Current behavior: only S1 help/version and unavailable placeholders exist.
Expected behavior: explicit `ci-see` executes selected validation; optional init
checks setup without Git integration. Workflow/job listing and event/JSON options
must be available when the command contract is finalized.
No implementation or command availability is claimed by this spec.

## 3. Dependencies

Users install CI See, `act`, Docker/Linux containers and Git for repository context.
Missing/unusable dependencies yield INCOMPLETE with installation guidance.
Help/version remain usable without execution tools. No automatic host-tool installation.
Prebuilt CI See needs no host Go/Node runtime; job runtimes are workflow dependencies.
Record engine version, executable identity, runner digest/platform and configuration.

## 4. Workflow Discovery

Discover GitHub-compatible `.yml`/`.yaml` workflow files directly in
`.github/workflows/`. Nested/non-workflow files do not become implicit pipelines.
Use existing YAML without mandatory edits or manual registration.
Malformed YAML, empty selection and unsupported required jobs cannot yield PASS.

Expose workflow/job listing and explicit workflow/job selection.
Report selection, dependency jobs, matrix scope and exclusions. An ambiguous job
selection across workflows must be rejected rather than silently broadened.
Validate paths stay in the selected repository; reject unsafe path traversal.
No selected work may disappear silently during discovery or execution.

## 5. MVP Event Policy

Support `pull_request`, `push` and `workflow_dispatch` for compatible workloads.
Retain the earlier per-workflow default preference: pull_request, then push,
then workflow_dispatch. Run each eligible workflow once, not once per trigger.
An explicit event overrides that preference and selects matching workflows.

Required missing context/dispatch inputs remain INCOMPLETE; do not fall back to a
different event to hide missing input. Report unsupported-only workflows and
event exclusions distinctly. Unsupported required selected jobs are non-PASS;
an intentional subset PASS is labeled as such.

Normal invocation derives available local context and accepts missing values
through CLI inputs; no hand-written event JSON required.
Exact option spelling, listing syntax and input precedence remain **CLI-01**,
to be finalized in a focused contract before CLI feature implementation.
Advanced payload support must follow the same validation/identity rules.
