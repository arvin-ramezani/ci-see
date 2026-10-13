# Implementation planning — Focused checks, acceptance and review

**Status:** Draft revision aligned with merged PRD #20
**Canonical authority:** [Implementation planning](../implementation-plan.md)

## 5. Required Test Matrix

Use [current traceability](../s0/cli-traceability.md), not retired gating/approval fixtures.

| Category | Required evidence |
| --- | --- |
| CLI/selection | Help/version, explicit execution/listing/options, missing inputs, ambiguity, nonempty scope, consistent exits/JSON |
| Source/events | Current dirty inputs, truthful refs/context, mutation detection, worktrees, path/link safety and selected policy |
| Dependencies/config | Missing tools, cold/warm reuse, trusted identities, corrupt/changed refs, denied/interrupted downloads and ambient override rejection |
| Runtime/reports | Nonzero/zero-but-incomplete results, failed jobs, crash/cancel/timeouts, atomic writes, mixed outcomes and parallel run isolation |
| Credentials/resources | Safe transport, redaction, explicit inheritance, run-owned cleanup and no unintended host/socket/privilege exposure |
| Hosts/release | Windows-native, Linux, macOS and separate WSL2 evidence for supported targets, missing runtime and prebuilt installation |

During implementation run focused meaningful deterministic tests, formatting,
type/build/vet checks and documentation checks appropriate to the change.
Fake act/process fixtures do not need Docker or network.
Record deferred/unrun real-tool tests accurately. No numeric coverage threshold
is introduced by this documentation revision.

Real act/Docker and cross-platform acceptance remain required before release,
performed under an independently reviewed scope and explicit owner authorization.
Do not execute large live fixture suites merely to update documentation.

## 6. Review, Merge and Release Gates

1. Documentation/plan: exact-head independent review and owner-approved merge.
2. Implementation: separate scoped branch/PR; deterministic verification and limitations.
3. Review: independent code review; security review for process, source, credentials,
   network or dependency trust changes. Do not represent implementer verification as independent review.
4. Merge: explicit owner authorization; no automatic merge or production claim.
5. Release: active slices S0–S5/S9/S10, relevant decisions and L-01–L-12 accepted;
   real-tool/host/install/recovery evidence reviewed, then explicit owner go/no-go.

Removed gates/provider/bypass and deferred hosted suppression are not release
conditions. Historical D3 acceptance is not silently promoted to PASS or erased.

## 7. Known Risks, Dependencies and Next Action

- Act does not fully reproduce hosted semantics; ignored required features cannot PASS.
- Checkout actions/source mutation can make HEAD-only reports misleading.
- Action/image verification must bind actual execution bytes, not just cache names.
- Act config/default files, process descendants and Docker resources can escape
  naive argv/cleanup assumptions; SEC-01 needs a tested contract.
- Secrets/permission and cancellation behavior need host-specific evidence.
- Downloading dependencies does not prevent runtime networking or hosted triggers.
- Reusing the old gate framework would add unused models and acceptance obligations.

**Next:** review/merge this documentation alignment, then authorize a separate
cleanup/command-contract slice. Do not change runtime code or close PR #19 here.
