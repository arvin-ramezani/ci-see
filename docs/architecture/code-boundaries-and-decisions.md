# System architecture — Go package boundaries, deferred decisions, next specs

**Status:** Draft (inherited)  
**Canonical authority:** [System architecture](../architecture.md)

<!-- BEGIN ORIGINAL CONTRACT -->
## 11. Code Boundaries

```text
cmd/
  ci-see/
    main.go
internal/
  cli/
  core/
  adapters/
    git/
    act/
    state/
    approval/
  hooks/
go.mod
```

Rules:

- `core` must not import infrastructure adapters.
- Git/`act`/filesystem details stay in adapters.
- CLI and hooks call core use cases instead of duplicating rules.
- Use Go modules and idiomatic packages; `internal/core` must remain independent of process, OS, and Git adapter implementations.
- Prefer Go standard library; introduce third-party dependencies only for justified requirements.
- Boundary data must be validated explicitly; use typed result states and errors.
- Run `go test ./...` and `go vet ./...`; add integration tests using isolated temporary Git repositories and stubbed `act` where appropriate.
- Tests must enforce exact-state PASS and fail-closed behavior.

## 12. Deferred Decisions

Focused specs or ADRs will define:

- event selection and generated event payloads;
- exact validation fingerprint inputs;
- developer approval/bypass UX;
- secrets, artifacts, and cache behavior;
- runner-image policy;
- duplicate GitHub-hosted execution control;
- native Windows/macOS runner behavior.

## 13. Next Specifications

1. Local CI execution specification.
2. Git gating specification.
3. Developer approval UX specification.
4. MVP implementation plan.

Implementation should begin only when the relevant behavior is specified enough to test deterministically. The MVP uses Go end to end for the CLI and orchestration; no Node.js/TypeScript application or runtime is required.
<!-- END ORIGINAL CONTRACT -->
