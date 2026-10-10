# Product requirements — Success criteria, decisions, technical validation, open questions

**Status:** Draft (inherited)  
**Canonical authority:** [Product requirements](../prd.md)

<!-- BEGIN ORIGINAL CONTRACT -->
## 11. Success Criteria

A developer can use an ordinary GitHub Actions repository and experience:

~~~text
change code
→ attempt commit
→ CI See verifies local CI state
→ PASS: commit proceeds
→ non-PASS: developer decides
→ push receives a second exact-state check
→ normal PR workflow continues
~~~

The MVP is successful when:

- users do not need to remember a separate CI command;
- AI agents cannot silently bypass a failed gate;
- GitHub Actions remains the CI source;
- PASS cannot be reused for changed code;
- the same core UX works across supported host operating systems.

## 12. Confirmed Product Decisions

- Product name: **CI See**.
- MVP implementation language: **Go**, compiled CLI with prebuilt cross-platform binaries; no Node.js/TypeScript app or Node.js runtime dependency. The Go toolchain is required only when building from source.
- GitHub Actions workflows remain the primary CI definitions.
- Local execution is a core capability.
- Default development flow validates locally before commit/push.
- Commit and push must be automatically interceptable.
- Failed or incomplete CI requires developer control over bypass.
- AI agents must not silently make that bypass decision.
- Cross-platform host support is required.
- Normal Git/GitHub developer workflow should remain familiar.

## 13. Current MVP Direction Requiring Technical Validation

- Go modules and prebuilt Go executables for the CLI, with the exact build/release toolchain pinned during implementation;
- nektos/act as the local workflow execution engine;
- Docker-compatible runtime as the execution substrate;
- Git hooks or equivalent interception for commit/push;
- CI See-owned native/desktop approval surface for agent-initiated Git operations.

## 14. Open Questions for Later Specs

1. What exact mechanism prevents unnecessary duplicate GitHub-hosted execution?
2. Which event should CI See choose when workflows expose multiple triggers?
3. What exact inputs form the validation fingerprint?
4. Should default execution run all eligible workflows or a policy-defined set?
5. How should the approval UI be implemented cross-platform?
6. What happens in headless sessions?
7. Which act runner-image strategy gives the best reliability/performance balance?
8. How should artifacts and caches behave locally?
9. How should local secrets be stored cross-platform?
10. How should bypass paths such as git --no-verify be handled?

<!-- END ORIGINAL CONTRACT -->
