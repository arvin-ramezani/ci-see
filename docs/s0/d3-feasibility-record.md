# D3 — Restricted feasibility record

**Date:** 2026-10-12  
**Status:** PARTIAL — restricted verification PASS; full D3 acceptance pending  
**Production ready:** NO

## Verified — disposable Windows 11 / WSL2 host tests

- Git staged-index snapshots, `commit -a`, `--amend`, explicit-path commits, linked worktrees and submodule gitlinks passed.
- Git hook ordering/mutation tests and WSL2 staged-tree cross-check passed; the post-commit mutation case demonstrated divergence from a pre-commit validated tree.
- Pinned cached checkout action integrity, required host dependencies, restricted Docker dependency probe and scoped cleanup checks passed.
- Six reviewed inbound helper Block rules were machine-verified. Outbound authorization gate and inspection mocks passed.

**Evidence:** `host-d3-d550e6f41bf94191a2a1d808b0f2b9a5`  
**Runner SHA-256:** `350AB5BD265F66A18300693E49B721185E5EF92DD113EB5257964A275EF48131`

**Evidence availability:** Externally retained on the owner's Windows test host under the `host-d3-d550e6f41bf94191a2a1d808b0f2b9a5` evidence directory (`evidence.json`). No sanitized immutable report is published in this repository or attached to the PR. The results here summarize owner-provided console output; remote reviewers cannot independently inspect or reproduce the underlying machine report from this reference alone. Do not publish the raw evidence before reviewing it for local paths and sensitive details.

## Not verified — no PASS claim

- No real `act` metadata, workflow jobs or GitHub Actions event-fidelity tests ran.
- Host-`act` outbound isolation, live job/container networking, runtime credentials, and commit-to-push PASS reuse remain unverified.
- Direct `--no-verify` bypass behavior and production caller enforcement need a separate accepted design; hook mutation detected after commit is not prevention.
- These are feasibility tests, not completed Go application or production acceptance.

## Owner decisions and handoff

- **S2 development may start in Git-only/fake-`act` scope** with strict non-PASS outcomes for unverified required behavior. Existing acceptance tests and review gates remain mandatory.
- **Real `act` execution remains blocked** until network and credential controls are independently verified and expressly authorized.
- **Dependencies:** acquire missing approved workflow dependencies, verify integrity before execution and reuse cached copies. See the [local-CI specification](../specs/local-ci-execution/results-and-acceptance.md#dependency-acquisition-and-reuse-2026-10-12). The restricted D3 runner remains offline.
- D3 remains PARTIAL; S4 live-`act` acceptance, complete D3 acceptance and production readiness are not granted.
