# Product requirements — Diagnostics, secrets, cross-platform hosts

**Status:** Draft revision — owner-requested CLI direction
**Canonical authority:** [Product requirements](../prd.md)

### FR-11 — Logs and diagnostics

Expose selection, progress, failures, skipped/excluded jobs, completion status,
missing tools/dependencies, downloads/cache reuse, and compatibility limitations.
Detailed engine logs remain available with secret redaction.

Distinguish workflow failures from setup/download/input errors and give actionable
next steps. Partial-selection PASS must not look like whole-repository or hosted validation.

### FR-12 — Secrets and variables

Support explicitly supplied local secrets/variables without committing them.
Separate committed configuration, local non-secret configuration, and local secrets.
Do not implicitly forward all host environment variables, credentials, or tokens.
Secrets must not appear in normal logs, reports, or diagnostics. Provide input
methods that avoid placing secret values in process arguments.

CI See does not sandbox arbitrary untrusted workflows. Document that repository
workflows/actions run code on local infrastructure; Docker alone does not establish
a complete host-security boundary.

### FR-13 — Cross-platform host support

MVP hosts: **Windows 11, macOS, Linux**, with prebuilt binaries for an explicit
supported OS/architecture matrix and consistent commands, exits, and reports.
Native Windows and Windows through WSL2 are separate tested environments.

Execution initially targets **Linux-container-compatible jobs through Docker**.
Cross-platform hosts do not imply native Windows/macOS runner jobs or identical
hosted environments. Unsupported runner labels must be reported, never silently
remapped and counted as equivalent validation.

Other OSs, additional architectures, native runner execution, and alternative
container engines are later feasibility work, not MVP guarantees.
CI See itself needs no Go/Node runtime with prebuilt binaries; `act`, Docker,
and workflow-specific tools remain dependencies.
