# Local CI execution — Actual source, execution, credentials and logs

**Status:** Draft aligned with merged PRD #20; source/security details pending
**Canonical authority:** [Local CI execution](../local-ci-execution.md)

## 6. Event Payloads

Generate run-owned temporary payloads from verified local Git context and explicit
inputs. Push includes known ref/commit/repository identity; PR includes known
head/base refs and identities. Never invent remote PR numbers, tokens or GitHub-only facts.
Unborn/detached branches or unknown base context are supported only when required
fields can be supplied truthfully; otherwise INCOMPLETE.

Validate required dispatch inputs/types/defaults. Reject conflicting/ambiguous inputs
and malformed payloads. Keep event inputs in tested-context identity, excluding secrets.
A dirty tested source is recorded separately from its available HEAD commit.
Generated payloads are not committed and cannot overwrite user files.

## 7. Execution Flow

1. Validate explicit selection, context and tools; prepare actual tested inputs.
2. Acquire/verify approved missing artifacts using existing engine/runtime caches.
3. Write RUNNING and create private run-owned execution inputs.
4. Launch `act` with argv and explicit effective configuration.
5. Collect all selected workflow/job/matrix outcomes, runtime errors and exclusions.
6. Clean run-owned resources, finalize durable report and return human/JSON results.

CI See never invokes Git commit/push or changes hooks/triggers.
Never launch workflow code before its required artifacts and effective configuration
are verified. Do not accept implicit pulls/cache substitutions that bypass that check.
Bound execution and cleanup; unresolved process ownership or incomplete termination
cannot become a successful completed run.

## 8. Execution Isolation

Each run binds one repository/worktree, actual tested files, selected scope,
event inputs, configuration, engine version and dependency identities.
Use distinct temporary/report directories; parallel runs cannot overwrite each other.

Default source is current local files, including uncommitted changes, rather than
an index-tree or pushed-tip gate snapshot. Do not attribute a dirty run solely to HEAD.
Prepare a stable execution input or detect changes reliably; ambiguous inputs are non-PASS.
No automatic reuse of historical PASS; every invocation requests a new execution.

**SRC-01 remains unresolved:** exact portable capture mechanism, ignored/untracked
selection, required Git metadata, symlink/submodule handling and mutation detection.
Resolve before source execution implementation; acceptance must prove the report
describes what the jobs actually consumed, including any checkout action behavior.
Do not blindly copy local credential files or materialize links outside the boundary.

**SEC-01 remains unresolved:** minimal effective network, engine configuration,
mount/socket/privilege, timeout and cancellation policy for supported trusted workloads.
Inspect all relevant `.actrc` and default env/secret/variable inputs; no silent
ambient override may change selection, networking, credentials or artifact identity.
Either control them or report a conflicting configuration as INCOMPLETE.
Docker/offline flags are not proof of complete host isolation.

No general untrusted-workflow sandbox is promised. The old D3 VM/signed-approval
procedure is historical test infrastructure, not the normal product runtime.
It remains unverified and unapproved for live execution in this docs-only task.

## 9. Secrets and Variables

Separate committed configuration, local non-secret settings and explicit secrets.
Use allowlisted environment names or explicitly selected private files; never
implicitly forward all host credentials or fetch a GitHub CLI token.
Do not pass plaintext secrets in argv or record their values/digests in normal reports.

Use supported engine transport with private run-owned files outside job source;
restrict permissions/ACLs, redact logs and remove transport files on all exit paths.
Missing required secrets or unsafe transport produce non-PASS with safe reason codes.
Scope download credentials separately from job secrets; no implicit shared credential store.
Exact cross-platform input/transport contract belongs to **SEC-01**.

## 10. Logs

Show selected workflows/jobs, progress, outcomes, exclusions, acquisition/cache reuse,
setup failures and compatibility warnings. Distinguish FAIL from INCOMPLETE.
Detailed engine logs must be available after secret redaction; raw secret output
must not be persisted or streamed unredacted as a diagnostic bypass.

JSON mode emits machine-readable results without engine log contamination.
Define stream placement and schema under **CLI-01**; human logs and JSON agree.
A report cannot label a partial selection as whole-repository/hosted validation.
