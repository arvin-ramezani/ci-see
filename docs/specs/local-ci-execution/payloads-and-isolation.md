# Local CI execution specification — Event payloads, execution, isolation, secrets, logs

**Status:** Draft (inherited)  
**Canonical authority:** [Local CI execution specification](../local-ci-execution.md)

<!-- BEGIN ORIGINAL CONTRACT -->
## 6. Event Payloads

Users must not need to create event JSON for normal usage.

CI See generates the temporary event payload needed by `act` from local Git/repository context.

For `pull_request`, the generated context should include at minimum:

- current/head branch;
- base/default branch when known;
- head commit SHA when available;
- base commit SHA when available;
- repository identity.

For `push`, it should include at minimum:

- current ref;
- current commit SHA;
- repository identity.

Generated payloads are temporary execution inputs and must not be committed.

If required event context cannot be determined safely, the affected workflow is INCOMPLETE rather than guessed.

## 7. Execution Flow

```text
ci-see
→ verify repository
→ discover workflows
→ verify act/runtime
→ resolve event for each eligible workflow
→ generate temporary event context
→ record RUNNING
→ invoke act
→ capture result/logs
→ persist normalized result
→ remove temporary sensitive files
→ return final status
```

The compiled Go CLI invokes `act` as a child process using Go `os/exec` and argument arrays, not shell-built command strings.

## 8. Execution Isolation

Each run is bound to:

- one repository/worktree;
- one exact local validation state;
- the selected workflows/events;
- the detected `act` version;
- relevant CI See execution configuration.

Concurrent runs must not overwrite each other's temporary files or result records.

Detailed validation fingerprint rules belong in the Git gating specification.

## 9. Secrets and Variables

CI See must not place secret values directly in generated shell command strings.

MVP sources:

- environment variables;
- optional repository-local private secret file under CI See Git metadata;
- optional repository-local private variable file under CI See Git metadata.

Conceptually:

```text
<git-dir>/ci-see/
  secrets.env
  vars.env
```

These files are outside tracked repository content.

CI See passes them to `act` through supported secret/variable mechanisms.

Rules:

- secret values must not be written to normal CI See logs;
- CI See must not copy secrets into result metadata;
- missing required secrets must produce a clear non-PASS result;
- file permissions should be restricted where the host platform supports it.

## 10. Logs

Normal output should show:

- workflow/job currently running;
- PASS/FAIL/SKIPPED outcome;
- engine/runtime failures;
- unsupported workflow/event reasons;
- final normalized result.

Detailed `act` output remains available for diagnosis.

CI See should preserve enough local run metadata to explain later why a result passed, failed, or was incomplete.

<!-- END ORIGINAL CONTRACT -->
