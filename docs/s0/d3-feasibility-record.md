**Document status: HISTORICAL — superseded for product scope by [merged PRD #20](../prd.md).**

All later status labels, authority claims, imperative wording and acceptance clauses
describe the earlier design/test scope only. They are not current CLI requirements
or permission to execute old recipes. Original text is preserved below. D3 results
remain PARTIAL/UNVERIFIED; this retirement does not establish runtime safety.

# D3 — Restricted feasibility record

**Date:** 2026-10-12
**Status:** PARTIAL — report records restricted verification results; full D3 acceptance pending
**Production ready:** NO

## Recorded host-test results

- The sanitized report contains **52 test-result records**: **26 PASS** and **26 UNVERIFIED**. The PASS labels are the report's recorded outcomes; the UNVERIFIED records remain unverified and are not counted as passes.
- It contains **174 inspectable command records**: 173 exit code 0 and one exit code 1; all are recorded as not timed out with complete output capture. These records make the commands and captured output inspectable; they do not independently reproduce the host run.
- Git staged-index snapshots, `commit -a`, `--amend`, explicit-path commits, linked worktrees and submodule gitlinks passed in the recorded results.
- Git hook ordering/mutation tests and WSL2 staged-tree cross-check passed. The late-hook mutation case demonstrates detection of divergence from a pre-commit validated tree; it does not prove complete prevention.
- Pinned cached checkout action integrity, required host dependencies, restricted Docker dependency probe and scoped cleanup checks passed in the recorded results.
- Six reviewed inbound helper Block rules were machine-verified. Outbound authorization gate and inspection mocks passed.

## Published evidence and integrity

- [Sanitized D3 evidence report](evidence/d3/evidence.sanitized.json) — SHA-256 `F5BD9C6AD01BEECFAC8DDCC9A20EC6B49A6D8B86474C8FFDCE8E569E808425FC` (264,704 bytes).
- [Sanitization manifest](evidence/d3/evidence.sanitization-manifest.json) — SHA-256 `7F852F73ECF955F279833E0FA2D1F97BEB6AE02C8FB4129D8960909CA1D40478` (2,997 bytes).
- The manifest cross-references the original evidence SHA-256 and byte count. The original report remains unchanged and unpublished. The manifest records `PUBLISHED_WITH_OWNER_AUTHORIZATION`; owner authorization came from the explicit task instruction to publish the complete corrected report and manifest after integrity and security checks passed, including the retained firewall identifiers and fixture metadata. The runner SHA-256 `350AB5BD265F66A18300693E49B721185E5EF92DD113EB5257964A275EF48131` and its locally-verified provenance are retained in the manifest. The runner artifact was not present in the supplied publication directory, so reviewers cannot independently recompute that hash from this PR.
- Source-to-sanitized review replaced 367 occurrences of the source identity with the consistent pseudonym `d3-user`. Nested-string scans found no original username, personal home path, email, credential/token pattern, private key, RFC1918 address, link-local address or IPv6 ULA. The captured host-network output retains only canonical loopback and unspecified address values; no host-specific private subnet address remains.
- The published JSON is independently inspectable evidence, not a claim that PR CI reran the host tests. Documentation and Go CI validate this repository change only.

## Not verified — no PASS claim

- No real `act` metadata, workflow jobs or GitHub Actions event-fidelity tests ran.
- Host-`act` outbound isolation, live job/container networking, runtime credentials, and commit-to-push PASS reuse remain unverified.
- Direct `--no-verify` bypass behavior and production caller enforcement need a separate accepted design. The late-hook mutation test detects divergence but does not prove prevention.
- These are feasibility tests, not completed Go application or production acceptance.

## Owner decisions and handoff

- **S2 development may start in Git-only/fake-`act` scope** with strict non-PASS outcomes for unverified required behavior. Existing acceptance tests and review gates remain mandatory.
- **Real `act` execution remains blocked** until network and credential controls are independently verified and expressly authorized.
- **Dependencies:** acquire missing approved workflow dependencies, verify integrity before execution and reuse cached copies. See the [local-CI specification](../specs/local-ci-execution/results-and-acceptance.md#dependency-acquisition-and-reuse-2026-10-12). The restricted D3 runner remains offline.
- D3 remains PARTIAL; S4 live-`act` acceptance, complete D3 acceptance and production readiness are not granted.
