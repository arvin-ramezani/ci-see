# Developer approval specification — UX, security, acceptance 1–23, deferred work

**Status:** Draft (inherited)  
**Canonical authority:** [Developer approval specification](../developer-approval-ux.md)

<!-- BEGIN ORIGINAL CONTRACT -->
## 12. User Experience Rules

- Keep the decision surface small and direct.
- Default focus must not make accidental bypass easy.
- Cancel is the safe/default outcome.
- Continue anyway must be visually explicit.
- Keyboard operation must be supported.
- Closing the UI must never mean approval.
- The terminal/IDE/agent receives a concise final allow/block message.

## 13. Security Boundary

**In-scope threat:** an AI agent can initiate Git/CLI commands, set process arguments/environment, write stdin, read output, and modify its checkout. None of these capabilities alone may authorize a CI See bypass.

**Residual risk:** a same-user agent with desktop automation, accessibility APIs, or control over the provider process/IPC may be able to click or forge approval. A normal Windows popup, including one launched from WSL2, does not by itself establish that a human approved it. Full host compromise and native Git `--no-verify` are outside CI See's enforcement guarantee.

Before shipping bypass, a security review must document the provider implementation, process/IPC permissions, UI automation exposure, caller-versus-approver trust separation, and limitations. Validate that the provider enforces the documented boundary. Where this cannot be demonstrated, disable Continue anyway and fail closed, rather than claiming agent-proof authorization.

The product must state its actual assurance level; never claim tamper-proof or guaranteed human-only enforcement on an uncontrolled local host.

## 14. Required Acceptance Tests

Implementation must prove:

1. PASS never opens approval UI;
2. non-PASS at a configured gate requests developer approval;
3. no CI See CLI flag can approve bypass;
4. stdin/piped input alone cannot approve bypass;
5. Cancel blocks Git;
6. Continue anyway returns ALLOW_BYPASS;
7. bypass applies only to the exact bound state;
8. bypass is one-time and does not become PASS;
9. changing state while approval is open invalidates the request;
10. closing the approval surface blocks Git;
11. unavailable/crashed approval provider blocks Git;
12. headless mode does not wait indefinitely or auto-approve;
13. bypass audit metadata contains no secret values;
14. terminal/IDE/agent receives an unambiguous final result;
15. Windows 11 + WSL2 can present the human-facing approval surface;
16. caller-controlled flags, environment, stdin, output, files, or process exit codes cannot forge approval;
17. provider response authentication, request binding, expiry, single-use consumption, and replay rejection work;
18. adversarial attempts to automate the provider UI or spoof its IPC are tested on Windows 11 + WSL2, with residual risks recorded;
19. a provider that fails trust-boundary verification has Continue anyway disabled and blocks non-PASS Git operations;
20. security review documents precisely which automated-agent capabilities are and are not resisted;
21. `ALLOW_BYPASS` is returned only after the corresponding audit record is fully and durably persisted;
22. disk-full, permission-denied, interrupted/partial write, or verification failure returns `BLOCK_ERROR` and blocks Git;
23. an audit record from one approval cannot authorize another operation, and a recorded approval must not be reported as successful Git completion.

## 15. Deferred

This spec does not define:

- final visual styling of native approval windows;
- remote/team approval;
- policy servers;
- tamper-resistant enterprise enforcement;
- prevention of native Git `--no-verify`.

Next: MVP implementation plan.
<!-- END ORIGINAL CONTRACT -->
