# System architecture — Execution lifecycle, trust and platforms

**Status:** Draft aligned with merged PRD #20
**Canonical authority:** [System architecture](../architecture.md)

## 7. Runtime Flows

### Local run

1. Explicit invocation resolves repository, workflows, selected jobs/events and inputs.
2. Validate prerequisites and actual tested-source preparation.
3. Resolve trusted action/image identities; acquire missing compatible artifacts.
4. Verify actual execution artifacts; create run-owned payloads and secret transport.
5. Invoke `act`; record per-job scope, outcome and runtime context.
6. Finalize report, redact logs, clean run-owned resources and return an exit/result.

Acquisition and execution are distinct credential/network boundaries, but normal
product use may need networking in both. No mandatory machine-wide air gap,
signed D3 manifest, VM or detached approval is part of normal CLI architecture.

### Commit gate

**Removed.** Ordinary commits neither invoke nor wait for CI See.

### Push gate

**Removed.** Ordinary pushes neither invoke nor wait for CI See.

## 8. Approval Boundary

**Retired product component.** There is no developer bypass, approval UI/token,
pending Git operation or audit-before-allow protocol.
Developer responsibility for invocation does not permit dishonest result reporting.

Historical D3 test procedures retain their explicit live-test authorization;
this design update does not authorize executing those procedures or workflows.

## 9. Reliability and Security

- Use `os/exec` argv, explicit executable discovery, bounded process cancellation.
  No shell-string construction or caller-supplied shell fragments.
- Capture tested inputs accurately; detected ambiguity/change prevents PASS.
- Do not implicitly inherit all host credentials, environment or engine configuration.
- Verify source/commit/digest identity against actual downloaded/cached artifacts
  before execution. Engine cache presence alone is not trust evidence.
- Apply explicit network policy; offline flags do not enforce network isolation.
- Reject accidental host/socket/privileged exposure through engine options.
- Redact secrets from reports/logs; file permissions and transport require host tests.
- Cleanup only resources whose ownership is proven for this run; preserve diagnostics.
- Local workflows run code. Docker does not make arbitrary untrusted actions safe.
- Unsupported required behavior and missing evidence cannot become PASS.

Detailed command/config, source-capture and network policies must be resolved
before their owning slice executes workflows. See the [execution contract](../specs/local-ci-execution.md).

## 10. Platform Strategy

Build separate native binaries for Windows, macOS and Linux. WSL2 uses a Linux
binary and needs distinct runtime evidence from native Windows.
Do not infer macOS support from cross-compilation or Linux-container success.

Keep OS-specific path handling, secret permissions and process-tree termination
behind small consumer-owned seams. Confirm supported OS versions/architectures,
tool versions and runner images before release; no broader OS guarantee.

The CLI needs no host Go/Node runtime when prebuilt; workflow jobs may need their
own tools. Keep Go build requirements in `go.mod` and development CI.
