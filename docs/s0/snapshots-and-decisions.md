# S0 — D3 staged/push snapshot feasibility and decision register

**Status:** Bounded support **proposal**, not a final D3 ADR or an alteration of contracts.
**Sources:** [Git identity/gating](../specs/git-gating.md), [local events/isolation](../specs/local-ci-execution.md), [D3 gate](../implementation-plan/authority-and-gates.md#3-decisions-and-feasibility-gates).
**Owner choice:** Common cases first; unverified cases block with explicit non-PASS, without pretending original acceptance obligations disappear.

## Exact execution identity — proposed design

1. Resolve repository and worktree through `git rev-parse`/Git-provided paths, including `--git-path`; do not assume a physical `.git` folder or reuse another worktree's state.
2. For pre-commit, honor Git's **effective index** (`GIT_INDEX_FILE` when supplied); obtain its tree object via Git. Generate an isolated checkout/export from that Git tree, **not from the mutable working directory**. Preserve executable bits, symlinks and gitlinks, or reject unverifiable materialization.
3. For push, use each non-deletion update's **new local object-id tip** supplied to `pre-push`, resolve its commit/tree, and materialize an immutable snapshot. Do not use current HEAD when it is not the ref's new tip.
4. Discover workflow YAML and other validation inputs **inside the materialized snapshot**. Record repo/worktree, Git tree, when applicable commit tip, workflow content, generated event semantics, config, runner, `act` version, variable identities and non-plaintext secret digests in the fingerprint.
5. Generated push/PR events must derive truthful refs/branch/base/head identity from Git, not invent GitHub-only facts. Explicitly distinguish an index tree with no final commit SHA from an actual committed tip. If required branch, default/base or dispatcher input is unavailable, mark the affected required validation INCOMPLETE.
6. Execute `act` against only that isolated snapshot, with separately protected temporary event/secret files and per-run directories. Never expose unstaged or untracked source paths to the engine by accident.
7. Post-commit bind PASS only after verifying the actual new commit **tree and all other fingerprint inputs** equal what was validated. Binding failure does not retroactively justify the commit; do not create a reusable PASS.
8. Recheck snapshot identity at relevant gate boundaries and before reusable PASS; unknown/error/collision/corruption is blocking. Never normalize empty/unsupported-only selection to PASS.

## MVP support vs fail-closed plan

| Git/event scenario | S2/S4 plan or explicit temporary boundary | Required verification |
| --- | --- | --- |
| Ordinary branch staged commit; mixed staged/unstaged/untracked | **Support:** stage-only immutable Git tree; ignore unstaged and untracked content | SN01 |
| Ordinary pushed tip; changed HEAD since ref selection; multiple refs | **Support:** commit object per new non-deletion ref, independently evaluated; deletion exempt | SN03, HK03 |
| Post-commit equivalent tree and changed workflow/config/event | **Support:** strict tree + full fingerprint equivalence; no guess from SHA alone | SN02, FP01 |
| `git commit -a`, `--amend`, partial/explicit-path commits | **Investigate effective temporary index:** if exact staged-to-commit correspondence cannot be proven, fail closed; post-commit must not bind a mismatched result | SN02, SN04 |
| Linked worktrees and Git metadata redirect | **Support isolated identity or fail closed:** never share incompatible worktree results | SN04, ST02 |
| Detached HEAD, unborn branch, missing base/default branch | **Support only known truthful event context**; unknown required PR/push context INCOMPLETE, not guessed | EV02, SN04 |
| Submodules, absent Git objects/LFS data, external path/symlink escape | **Reject unverifiable/mutable checkout or unsafe materialization**; no PASS from working directory leakage | SN04 |
| Modified user/tool hooks or hook that edits index after CI See | **D3 unresolved**: ensure validation runs after allowed synchronous index mutations or block unsafe composition; post-commit mismatch alone cannot protect a commit already created | HK05, SN04 |
| Zero executable workflows, unsupported-only workflows or missing required dispatch input | **INCOMPLETE / non-PASS** whenever required validation cannot execute; never vacuous PASS | W02, EV02 |

**Important:** Gating acceptance requires real temp-repo tests for edge cases. "Fail closed" is an interim behavior, **not evidence that a mandatory supported operation was implemented**. Unsupported common-operation behavior that contradicts a focused spec needs reviewed ADR/contract change, not an S0 unilateral narrowing.

## Decisions and blockers

| Gate | S0 status | Required next evidence and owner review |
| --- | --- | --- |
| **D3** snapshot/event | **PARTIAL; not closed.** Exact tree/tip model proposed. | S2 feasibility spike: effective index for path commits and `-a`; submodules/materialization safety; PR event fields; user-hook ordering/index mutation. Review against all relevant acceptance tests; record ADR if a required scenario cannot be supported. |
| **D1** human approval | **OPEN**, separate S6 feasibility | Threat model, authenticated one-time provider, Windows/WSL2 IPC/automation attacks; never enable bypass without verified provider. |
| **D2** remote duplicate CI | **OPEN**, separate S8 gate | Reviewed suppression mechanism satisfying FR-10 without weakening hosted required checks; otherwise no completed MVP. |
| **D4** release/platforms | **OPEN**, separate S9 gate | Explicit GOOS/GOARCH, cancellation strategy, binary naming, versions/dependency minimums; Windows native and WSL2 independently checked. |

No decision authorizes `--force`, caller-controlled approval, silent remote skip, local hook tamper-proof claims, or a false Go/application PASS.
