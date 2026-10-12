# Issue #17 — D3 Phase A real-act verification plan

**Status:** Phase A plan; draft for independent security and architecture review  
**Authority:** This plan maps the [D3 feasibility record](d3-feasibility-record.md), [fixture catalog](fixtures.md), [snapshot decisions](snapshots-and-decisions.md), [local CI specification](../specs/local-ci-execution.md), [Git gating specification](../specs/git-gating.md), [architecture](../architecture.md), and [implementation gates](../implementation-plan/authority-and-gates.md). Canonical contract text remains authoritative.  
**Baseline:** main at aa1aa266c41f26dd1a9b78002cc887a97bd592a7 (recheck before execution).  
**Current decision:** **D3 PARTIAL; real act BLOCKED; PRODUCTION_READY = NO.**

This is a plan, not execution evidence. It adds no product implementation, does not close any acceptance requirement, and does not authorize real act. D3 feasibility observations describe a reviewed environment and run. Product acceptance requires tests of CI See behavior after its S4/S5 implementation exists. Every requirement not independently demonstrated remains **OPEN** or **INCOMPLETE**.

## Status and evidence rules

- **D3 feasibility** means characterize tools, host boundaries, event/snapshot behavior, or evidence quality in an approved disposable environment.
- **Product implementation acceptance** means prove CI See enforces the required behavior in its implemented S4/S5 paths.
- **Both** means Phase B may gather feasibility observations, while the product claim stays open until a separate S4/S5 implementation test passes.
- Phase B starts only after independent security review PASS, resolution of required architecture choices, and explicit owner approval for the exact hosts, workflows, endpoints, and commands. Phase B uses only reviewed, pinned, disposable workflow sources. Never run arbitrary or unreviewed workflows.
- A positive control succeeding does not prove isolation. A negative control failing does not prove enforcement. Record configured policy and independent observations separately for each process/network boundary.
- Isolation is **INCOMPLETE** if enforcement configuration or independent observation is missing, ambiguous, bypassable, or not attributable to the tested component.
- No successful feasibility row changes product acceptance to PASS. Keep product rows OPEN until CI See enforcement is independently tested against canonical acceptance clauses.
- Preserve private original logs. Publish only a privacy and integrity reviewed sanitized report with source identity, command/run identity, timestamps, exit/timeout outcomes, hashes, byte counts, and a redaction manifest. Never publish secrets, user paths, or internal addressing.

## Environment preflight

No preflight or host mutation occurs in Phase A. In Phase B, inventory the selected disposable VM before any test. Do not install, download, upgrade, or reconfigure host tools as part of restricted verification.

| Prerequisite | Verification and evidence | Blocker behavior |
| --- | --- | --- |
| Disposable VM capability | Record VM platform/configuration and guest evidence that nested virtualization extensions are available. Confirm WSL2 is supported and usable in this VM without changing host/VM configuration. Record provider documentation or administrator evidence where guest inspection cannot prove capability. | Missing/unclear nested virtualization or WSL2 support blocks that host leg; do not silently substitute another environment. |
| Windows-native leg | Record Windows build, architecture, native Git path/version, act path/hash/version, container runtime client/server identity, daemon endpoint, runner image/digest/platform, and test account. | Missing/mismatched tool is an environmental blocker. No installation or download. |
| WSL2 leg | Separately record WSL/kernel and distribution identity, Linux binary/tool versions, Git, act, runtime endpoint, image digest/platform, path translation and process boundaries. Do not infer WSL2 behavior from the Windows-native leg. | Missing/mismatched WSL2 prerequisite blocks this leg independently. |
| Go restricted verification | Require Go **1.27.2**. Run with GOTOOLCHAIN=local and record go version, go env GOTOOLCHAIN, and host identity. | If local Go is not exactly 1.27.2, mark Go verification BLOCKED by environment. Do not allow Go to fetch a toolchain or modify host tools. Phase A documentation checks do not require Go. |
| Clean test ownership | Record starting processes, containers, networks, listeners, ports and test-owned labels/IDs. Create a unique run identity; cleanup may target only resources proved to belong to this run. | Unknown ownership, failed inventory, or inability to clean safely blocks execution. |

Record tool hashes and runtime/image identities before use. For actions, record immutable commit and verified tree/materialized content. For images, require the trusted OCI manifest digest for the selected platform and evidence the runtime will use that digest. A mutable tag, local image ID, or directory name alone is insufficient.


## Runtime security profile and go/no-go checks

The target is an owner-selected, disposable Windows 11 VM with a dedicated disposable container daemon. Treat the daemon API as VM-administrator access: only the reviewed act host process may reach that daemon endpoint, and no job, service, action, or helper container may receive a daemon socket, Windows named pipe, or daemon credential. A shared host/remote daemon or an endpoint that cannot be attributed and confined to the disposable VM is a NO-GO.

| Surface | Allowed profile | Forbidden exposure / negative control | Go/no-go evidence |
| --- | --- | --- | --- |
| Daemon socket / Windows named pipe | Reachable only by the reviewed act host process from the disposable VM account, through the dedicated VM daemon endpoint. | No socket, named pipe, daemon credential, or bind mount in any job/service/action container. From each workflow container, attempt a harmless API version request against the known socket/pipe path; it must fail. | Save daemon endpoint/config and ACL/owner inspection separately from container mount/process inspection and the failed in-container probe. Any reachable API or unknown endpoint blocks before a job starts. |
| Host filesystem and WSL interop | Mount only the disposable fixture worktree and run-owned temp/cache paths required by act; fixture input read-only where supported. | No Windows drive roots, user profile, SSH/config/credential directories, arbitrary host paths, /mnt/c or WSL interop shares, host home, or secret files. Attempt reads of pre-created non-secret sentinel files outside the allowlist; they must be absent/unreadable. | Compare reviewed workflow/action mounts with actual daemon mount metadata, path translation, WSL interop settings, and sentinel probe results. Any extra or uninspectable mount blocks. |
| Privilege, capabilities, namespaces | Dedicated non-admin workflow account; unprivileged container profile; no added capabilities; read-only root where compatible, no-new-privileges, restrictive runtime security profile, per-run network namespace. | Privileged mode, added capabilities, devices, host network/PID/IPC namespaces, namespace sharing with unrelated runs, or container options that weaken isolation. A harmless identity/capability/namespace probe must show the allowed profile and fail to elevate. | Check workflow and effective act/container options before start, then inspect actual container HostConfig, security options, namespace IDs and capabilities. Unsupported or mismatched enforcement blocks before execution. |
| Ports and local routes | No published ports. Only explicitly approved intra-run bridge paths, if a fixture requires them. | Published host ports, host loopback, unrelated bridge peers, VM services, and host gateway paths. Probe each from job and service containers; prohibited paths must fail. | Compare requested ports/networks with actual daemon inspection and independent VM-side listener/connection observations. Any unexpected listener or reachable host service blocks. |
| Credential helpers, stores and secrets | Dedicated test account with no personal or production credentials; only generated dummy canaries explicitly required by a fixture. An explicit dummy GITHUB_TOKEN prevents implicit token lookup. | No gh auth, cloud credentials, Git credential helper/store, SSH agent/key, Docker credential helper/store, actual token, or inherited secret in the VM test account or child environment. | Inspect helper configuration and credential-source inventory without printing values; trace access attempts; check sanitized environment/argv/logs and canary containment. Any real credential present, helper/store access, or canary escape blocks. Never use a real credential to test masking or isolation. |

Run a static workflow/config and effective-runtime-profile review before any workflow execution. A failed negative control, unexpected exposure, unsupported profile control, or inability to inspect effective state is NO-GO. Do not weaken runtime protections during a run; stop and resolve the design first.

## Network isolation controls

Before Phase B, the owner and independent security reviewer must select a disposable Windows 11 VM and approve a written destination policy for acquisition, engine metadata, workflow execution, and teardown. Enforcement must be independent of the guest: place the VM behind an owner-controlled, outside-guest deny-by-default virtual-switch or upstream network policy, and capture traffic at that boundary. The selected platform and exact commands remain an architecture decision; if it cannot enforce and capture VM traffic independently, Phase B stays BLOCKED. Do not change a personal host firewall, global Docker settings, or shared network policy.

The outside-guest boundary is the egress backstop, not proof of internal isolation. Before any act network attempt, demonstrate its effective deny-by-default policy with non-act controls. Then separately demonstrate enforcement and observation for native act, child/helper processes, WSL2, the dedicated daemon, and each job/service container. Use process/runtime controls inside the disposable VM plus independently captured outside-guest traffic. External capture must reconcile all VM egress; local observations must attribute allowed/denied attempts to the tested process/container. Loopback, bridge, daemon API, named-pipe, and WSL interop paths do not necessarily cross the VM egress boundary and require runtime-profile checks.

For every applicable component below, test both an **approved positive control** and a **prohibited negative control** for IPv4 and IPv6, and separately inspect DNS and proxy behavior. Use controlled endpoints/sinks with independent logs or packet/connection observations. Record configured enforcement policy, effective state, observed source process/container, destination, protocol, timestamp, and result as distinct evidence. Test loopback, bridge, host-gateway, daemon endpoint, proxy bypass and DNS paths explicitly. Include action fetch and image pull as acquisition paths, distinct from workflow runtime egress. The same process/container boundary must generate each control; a separate shell probe is only a transport control, not evidence about act or a workflow container.

| Component/boundary | Positive control | Negative control | Independent observation and policy evidence |
| --- | --- | --- | --- |
| Host act process | Connect to each explicitly approved host-process destination for its authorized phase. | Attempt controlled IPv4/IPv6 and DNS access to prohibited destinations; test approved proxy versus prohibited proxy bypass. | Process-attributed connection trace plus effective host policy/rule inspection; record act process identity and descendants. |
| act child processes/helpers | Exercise an approved helper/control connection only where policy permits. | Attempt prohibited direct, proxy, DNS, IPv4 and IPv6 routes from each helper class. | Per-process tree and independently observed socket/connection events; effective per-process policy. |
| Docker daemon | In the acquisition phase, resolve and fetch only an approved pinned image from the approved registry. In execution, exercise only explicitly allowed daemon destinations. | Attempt an unapproved registry, DNS destination, proxy, and prohibited IPv4/IPv6 destination. | Daemon logs/connection trace and effective daemon policy/config; independently verify pulled manifest digest and platform. |
| Job containers | Reach an approved test endpoint only if the execution policy explicitly allows that class of workflow traffic. | Attempt prohibited IPv4/IPv6, DNS, proxy, bridge, loopback-to-host, and published-port paths from inside the job. | In-container probe result plus host/daemon-side connection observation; inspect actual network mode/routes and effective policy. |
| Service containers | Reach only an explicitly approved peer/control endpoint if required by the fixture. | Attempt prohibited egress and prohibited host/bridge/peer paths. | Probe within the service plus independent runtime/network observation and inspected network configuration. |
| Action/dependency acquisition | Fetch a reviewed immutable action/image from an approved source during the separately authorized acquisition leg. | Attempt unapproved source, mutable identity, cache substitution, and prohibited destination. | Source-side request evidence, resolved commit/digest, local object/content verification, and proof execution selected that verified artifact. |

If a component cannot be isolated and observed separately, record that gap and classify the security result **INCOMPLETE**. Docker internal networks, offline flags, disabled pulls, an allowlist, or firewall rules alone are not proof of host-act isolation. Positive and negative controls must run through the same component boundary as the claim they support.


## Abort, rollback, and emergency cleanup

Use these stop points for every Phase B leg. The operator records UTC time, run ID, trigger, process/resource identities, and last safe checkpoint in the private evidence log. Any stop condition makes that leg BLOCKED/INCOMPLETE; never continue to a later act stage to investigate.

1. **Before acquisition:** verify written approval, VM snapshot/restore point, outer deny-by-default policy and independent capture, runtime go/no-go profile, and run-owned resource inventory. If any proof is absent, do not invoke act. If policy or observation is lost, stop before network access.
2. **During acquisition or metadata:** on an unexpected destination, failed control, policy drift, unapproved action/image identity, secret/helper access, timeout, or inconclusive capture, stop the command and block all workflow execution. Preserve logs and hashes privately; do not retry through another path.
3. **During a job:** on a prohibited connection succeeding, canary exposure, unexpected mount/privilege/namespace/port, missing telemetry, timeout, cancellation failure, or runaway child, request bounded graceful cancellation for the run; then terminate only process IDs verified as descendants of this run, accounting for PID reuse. Do not kill an ambiguously owned process or shared daemon.
4. **Resource rollback:** enumerate resources using the unique run ID/labels. Stop and remove only containers, networks, volumes, temp files, listeners, and processes whose ownership is independently confirmed. Recheck process trees, daemon resources, bridge networks, and VM listeners against the pre-run inventory. If ownership is uncertain, do not delete; quarantine the VM and mark cleanup INCOMPLETE.
5. **VM rollback/discard:** if a child escaped the run boundary, evidence is inconclusive, a real credential may have been exposed, policy could not be restored, or post-stop checks differ unexpectedly, disconnect or suspend the disposable VM using the owner-controlled hypervisor. Retain it only for approved private forensics, then restore a known clean snapshot or discard the VM. Do not perform cleanup on the physical host or shared daemon. Record snapshot identity and restore/discard result.
6. **Closeout:** seal raw evidence in the approved private location, hash the bundle, record redaction and custody, and publish nothing until a separate reviewer approves sanitization. Failed cleanup, an unaccounted resource, absent post-stop listener checks, or lost evidence custody leaves isolation INCOMPLETE and bars Phase B completion.

Do not attempt emergency cleanup of unrelated resources. If an outer boundary cannot contain an active runaway process, suspend/disconnect the disposable VM through the hypervisor; if unavailable, this is a preflight blocker and Phase B must not start.

## Acceptance matrix for the 26 UNVERIFIED records

**Host shorthand:** Windows-native and WSL2 are separate legs wherever applicable; each row records which leg ran.  
**Preflight shorthand:** disposable VM, tool identities, approved endpoint policy, dummy-only credentials, run ownership, and independent evidence capture above. Missing prerequisites mean **BLOCKED/INCOMPLETE**, never PASS.  
**Authority shorthand:** LCI = local CI spec and its linked discovery/events, payloads/isolation, and results/acceptance topics; Git = Git gating spec and linked modes/commit, identity/outcomes, and bypass/acceptance topics; Arch = architecture and D3 gate; S0 = fixtures, snapshot decisions, and D3 feasibility record.

| Evidence ID | Class; executable phase | Authority / fixture | Verification, independent oracle, retained evidence, and non-PASS rule |
| --- | --- | --- | --- |
| ACT-CREDENTIAL-ISOLATION-PREFLIGHT | D3 feasibility; Phase B preflight before any act call | LCI §§3, 9; Arch §9; preflight | Inspect child environment allowlist, isolated home/config, explicit empty token behavior, endpoint/config, and absence of host credential mounts without reading secret values. Compare process configuration with runtime/API-access observation. Keep UNVERIFIED until runtime tracing proves no host credential fallback/access; otherwise INCOMPLETE. |
| ACT-HOST-OUTBOUND-NETWORK-ISOLATION | D3 feasibility; Phase B gate before metadata/jobs | Arch §9; network controls | Run approved/prohibited controls for host act and children; compare effective enforcement policy with process-attributed network observations. Any unobserved path, including IPv6, DNS or proxy, is INCOMPLETE and blocks further act execution. |
| GIT-HOOK-BYPASS-NO-VERIFY | Both; Phase B safe disposable Git fixture; S5 product acceptance later | Git §12; HK05, SN04 | In a newly created disposable repository, configure a local pre-commit hook that writes a marker and exits nonzero; verify a normal commit is rejected, then verify `git -c core.hooksPath=<fixture-hooks> commit --no-verify -m fixture` succeeds and leaves the hook marker unchanged. Compare exit codes, commit object and marker. This proves native bypass semantics only, not a CI See pass; retain UNVERIFIED if the fixture is not independently observed. Product acceptance must report limitations honestly; never claim hooks are tamper-proof. |
| ACT-EXECUTION-BLOCKED | D3 feasibility; Phase B authorization-gate regression before any act call | Arch D3 gate; preflight | Exercise the existing authorization-gate test/mocks with missing approval or isolation evidence. Oracle: zero act metadata/job and PASS-reuse invocations. Retain gate inputs, decisions, call counters and logs; any act call or reusable PASS is failure. |
| ACT-METADATA-VERSION-HELP | D3 feasibility; Phase B only after all gates | LCI §§3, 11; AC05 | Record reviewed act binary identity/hash and then its metadata/version output; compare reported version with executable identity. Retain sanitized output and command result. Unavailable or mismatched identity is INCOMPLETE. |
| ACT-FIXTURE-SEPARATION | D3 feasibility; Phase B only after all gates | LCI §§6–8; EV01, EV02, SN01 | Run only isolated pinned fixtures with unique temp state. Compare selected worktree/tree and observed event files with independently generated fixture Git objects. Any mutable checkout leakage, shared state, or unreviewed source is BLOCKED/INCOMPLETE. |
| ACT-INTERNAL-DOCKER-NETWORK | Both; Phase B network leg; S4 product test later | Arch §9; network controls; ACT fixture | Inspect actual network mode and routes, then run approved and prohibited controls from relevant job/service containers; correlate with independent daemon/host observations. Keep S4 network policy/enforcement OPEN until CI See configures and verifies it. |
| ACT-PRIVATE-INTERNAL-GIT-TRANSPORT | Both; Phase B only with specifically approved controlled transport; S4 later | LCI dependency policy; network controls | Observe reviewed private/internal Git fixture transfer; verify remote identity, commit/tree and materialized content independently. Test prohibited transport/source too. No approved source or unverifiable content means INCOMPLETE; product acquisition enforcement remains OPEN. |
| ACT-PUSH-CONTEXT-CHECKOUT | Both; Phase B workflow leg; S4/S5 later | LCI §§6–8; EV01, SN03 | Compare checkout ref, SHA and tree to independently resolved pushed non-deletion tip. Retain event and checkout observations with sanitized identity. Mismatch or guessed context is INCOMPLETE; generated product context remains OPEN. |
| ACT-PULL_REQUEST-CONTEXT-CHECKOUT | Both; Phase B workflow leg; S4/S5 later | LCI §§5–6; EV01, EV02 | Compare head/base refs and SHAs, repository identity and checkout tree with controlled fixture Git oracle. Missing required context must yield non-PASS. Live engine observation does not prove CI See payload generation; that acceptance remains OPEN. |
| ACT-WORKFLOW_DISPATCH-CONTEXT-CHECKOUT | Both; Phase B workflow leg; S4/S5 later | LCI §§5–6; EV02, W02 | Test a pinned dispatch workflow with known inputs and one missing required input; compare observed payload and checkout with fixture oracle. Missing/ambiguous input must be INCOMPLETE; product-generated payload acceptance remains OPEN. |
| ACT-LABELLED-CONTAINERS-INSPECTED | Both; Phase B resource inspection; S4 cleanup acceptance later | LCI §§7–8; AC03, ST02 | Inventory containers/networks before and after; inspect labels, image digest, mounts, privilege/capabilities, ports and network. Correlate runtime IDs with unique run label. Unknown ownership or unsafe properties blocks; CI See scoped creation/cleanup remains OPEN. |
| ACT-CACHED-CHECKOUT-RESOLUTION | Both; Phase B dependency leg; S4 later | LCI dependency policy; AC01, AC02 | Verify approved immutable action identity against Git object/tree and materialized bytes, including warm-cache reuse and changed-ref case. Compare with trusted origin/resolved commit and filesystem hashes. Cache name alone or unverified substitution is INCOMPLETE; product enforcement remains OPEN. |
| ACT-CREDENTIAL-ISOLATION-RUNTIME-AUDIT | Both; Phase B runtime audit; S4 later | LCI §§9–10; AC04, FP02 | Use dummy canaries only. Inspect sanitized env names, files, argv, logs, mounts and process access traces; test missing/empty values and cleanup. Oracle is canary containment and no secret-value disclosure. Any unexplained host credential access/leak is failure; product transport remains OPEN. |
| ACT-FAIL-CLOSED-PUSH | Both; Phase B controlled run; S4/S5 product tests later | LCI §§11, 14; Git §§7, 9, 11; SN03, ST01, HK03–04 | Cause a reviewed fixture job to fail; compare raw engine exit/job results with expected FAIL and non-PASS oracle. Capture run and ref identities. Phase B describes engine behavior only; CI See result normalization and Git blocking remain OPEN. |
| ACT-FAIL-CLOSED-PUSH-MISMATCHED-SHA | Both; Phase B identity case; S4/S5 later | Git §§7–10; SN03, FP01 | Pair a controlled pushed tip with a deliberately different expected SHA. Independently resolve both commit/tree objects and check observed checkout. Any mismatched identity treated as PASS is failure; product gate enforcement remains OPEN. |
| ACT-FAIL-CLOSED-PULL_REQUEST | Both; Phase B controlled run; S4/S5 later | LCI §§5–6, 11, 14; EV02, ST01 | Use reviewed PR fixture with required context withheld or mismatched. Compare result and payload to fixture oracle. Missing required context must remain non-PASS; CI See event handling and gate outcome remain OPEN. |
| ACT-FAIL-CLOSED-WORKFLOW_DISPATCH | Both; Phase B controlled run; S4/S5 later | LCI §§5, 11, 14; EV02, W02 | Omit a required dispatch input in pinned fixture. Inspect engine invocation and result; compare with expected INCOMPLETE/no reusable PASS. Product-side preflight and gate behavior remain OPEN. |
| ACT-FAIL-CLOSED-HISTORICAL-REPOSITORY-MISMATCH | Both; Phase B identity case; S4/S5 later | Git §§8–10; SN04, FP01 | Supply a fixture identity from another disposable repository and compare repository/worktree IDs, tree and fingerprint inputs. Mismatch must not be reused as PASS. Product repository isolation remains OPEN until implementation tests pass. |
| PASS-REUSE-STAGED-TREE-TO-COMMIT-BINDING | Both; Phase B derive feasibility oracle; S5 product test later | Git §§4, 6, 8; SN01–02, FP01 | Derive staged tree and committed tree independently; vary tree and fingerprint dimensions. Retain Git object identities and run fingerprints available. Only exact tree plus complete identity is eligible; CI See binding remains OPEN. |
| PASS-REUSE-PRECOMMIT-CI-VALIDATION-SUCCESS | Both; Phase B collect eligible run evidence; S5 later | Git §§4–6, 9; SN01–02, HK02 | Record a successful controlled validation’s staged tree, workflow/event/tool identities and outcome. Compare every identity component to later commit oracle. Feasibility output alone cannot establish reusable CI See PASS; S5 stays OPEN. |
| PASS-REUSE-COMMIT-TO-PUSH-EVENT-EQUIVALENCE | Both; Phase B derive equivalence oracle; S5 later | Git §§7–10; SN02–03, FP01 | Compare commit and push event fields from controlled Git refs and actual event observations. Record differing inputs and whether contracts define equivalence. Any unresolved semantic choice blocks PASS reuse; S5 stays OPEN. |
| PASS-REUSE-COMMIT-TO-PUSH-EQUIVALENT | Both; Phase B derive positive identity case; S5 later | Git §§6–9; SN02–03, FP01 | Construct identical tree and complete execution fingerprint, then compare independent Git/event oracle with observations. A positive feasibility case is not a product reuse PASS; S5 matching logic remains OPEN. |
| PASS-REUSE-CHANGED-DIMENSION-MISS | Both; Phase B derive negative identity cases; S5 later | Git §§8–10; FP01–02 | Change tree, workflow, event input, config, runner/image, act version, variable or secret digest one at a time. Verify independent fingerprint input diff and expected stale/miss outcome. Any reused PASS on changed identity fails; S5 remains OPEN. |
| PASS-REUSE-IDENTICAL-COMPLETE-IDENTITY | Both; Phase B derive complete identity case; S5 later | Git §§8–9; FP01–02, ST01 | Compare complete repository/worktree, tree, workflow/event, policy, engine/image and private variable/secret digest identity. Retain no plaintext secrets. Missing fingerprint evidence is INCOMPLETE; S5 reuse remains OPEN. |
| PASS-REUSE-COMMIT-TO-PUSH-EVENT-MISMATCH | Both; Phase B derive negative event case; S5 later | Git §§8–10; SN03, EV02, FP01 | Vary a meaningful event/ref/base/head input while keeping other fixture fields fixed. Compare actual payload and independent Git oracle. Mismatched event identity cannot reuse PASS; product enforcement remains OPEN. |

## Additional required scenarios

Carry these cases into the same matrix even where no dedicated UNVERIFIED result ID exists:

| Scenario | Fixture / source | Phase B evidence | Product acceptance still required |
| --- | --- | --- | --- |
| Dependency cold/warm cache, changed pin, tampered/stale entry, interrupted/unavailable acquisition, unapproved source | LCI dependency policy; AC01–02 | Source, resolved commit/digest/platform, integrity checks, daemon/network observations and actual selected artifact | S4 must enforce verified identity before execution and return INCOMPLETE on uncertainty. |
| Cancellation, timeout, child process and teardown | AC03; ST01–02 | Process/container tree before and after, timing, final non-PASS result, listener/resource inventory | S4/S5 bounded cancellation, non-PASS persistence and owned cleanup. |
| Concurrent runs, crashed writer, truncated/corrupt state | ST02; Git §13 | Isolated run IDs/directories, injected controlled interruption, no cross-associated outputs | S5 atomic state, locking and fail-closed gate behavior. |
| Unsupported-only/zero workflows and missing required inputs | W02; LCI §§4–5, 14 | Selected workflow/event inventory and explicit unsupported/incomplete reason | S4 never treats required unsupported or unexecuted validation as PASS. |
| Cleanup and unrelated-resource preservation | AC03; D3 cleanup policy | Before/after inventory; prove only run-labeled resources are eligible and unrelated resources remain | S4 cleanup ownership and failure handling. |


## Phase B command catalogue (not executed in Phase A)

These reviewed command shapes are for Phase B only. Fill values from the approved run manifest. Pin executable paths/hashes, workflow commit/tree, event JSON, runner image platform digest, approved/prohibited endpoints, and run ID first. Use a fresh disposable account and per-leg working directory. Never enable shell tracing; retain stdout/stderr and exit status privately. The act flags follow the [act usage guide](https://nektosact.com/usage/index.html) and must be checked against the pinned binary help before approval. If a flag differs, stop and revise this document through review; do not improvise.

### Windows-native leg (PowerShell)

~~~powershell
# Phase B only, inside the owner-selected disposable Windows 11 VM.
$ErrorActionPreference = 'Stop'
$PSNativeCommandUseErrorActionPreference = $false
$env:GOTOOLCHAIN = 'local'
$RunId = '<approved-unique-run-id>'
$Workflow = '<absolute-path-to-reviewed-workflow-at-pinned-tree>'
$EventFile = '<absolute-path-to-reviewed-event-json>'
$EventName = '<push|pull_request|workflow_dispatch>'
$RunnerImage = '<registry/repository@sha256:platform-manifest-digest>'
$ActExe = '<absolute-path-to-reviewed-act.exe>'
$ApprovedHost = '<approved-control-hostname>'
$ApprovedPort = <approved-tcp-port>
$ProhibitedHost = '<prohibited-control-hostname>'
$ProhibitedPort = <prohibited-tcp-port>

Get-ComputerInfo | Select-Object WindowsProductName, WindowsVersion, OsBuildNumber, OsArchitecture
wsl.exe --status
wsl.exe --list --verbose
Get-Command git.exe, docker.exe, $ActExe | Select-Object Name, Source
& git.exe --version
& docker.exe version
& go.exe version
if ((& go.exe version) -notmatch '^go version go1\.27\.2 ') { throw 'BLOCKED: local Go 1.27.2 unavailable' }
if ((& go.exe env GOTOOLCHAIN) -ne 'local') { throw 'BLOCKED: GOTOOLCHAIN is not local' }
& docker.exe context show
& docker.exe info --format '{{json .}}'  # inspect privately; redact before publication

# Transport controls only; enforcement config and outside-guest capture are separate evidence.
$positive = Test-NetConnection -ComputerName $ApprovedHost -Port $ApprovedPort -InformationLevel Detailed
if (-not $positive.TcpTestSucceeded) { throw 'BLOCKED: approved positive control failed' }
$negative = Test-NetConnection -ComputerName $ProhibitedHost -Port $ProhibitedPort -InformationLevel Detailed
if ($negative.TcpTestSucceeded) { throw 'FAIL: prohibited negative control succeeded' }

# Only after outer enforcement/capture and all preflight controls PASS:
& $ActExe --version
& $ActExe --help
# Exactly one reviewed fixture/event; explicit dummy token prevents implicit auth lookup.
& $ActExe -W $Workflow -e $EventFile -P "ubuntu-22.04=$RunnerImage" -s "GITHUB_TOKEN=CISEE-DUMMY-$RunId" $EventName
$ActExit = $LASTEXITCODE
# Record exit status, then inspect only resources tied to this run by daemon create-event evidence.
~~~

Expected preflight: identity commands succeed, Go is exactly 1.27.2/local, VM capabilities are independently evidenced, daemon and VM identities match the manifest, outer policy is active, and safety negative controls are denied. Each fixture's expected exit/result is declared in the manifest; mismatch is not PASS. Missing identity, capability, capture, enforcement, or policy evidence means do not invoke act.

### WSL2 leg (PowerShell launch plus Linux shell)

~~~powershell
# A distinct leg using a dedicated WSL distribution and no Windows-path worktree.
wsl.exe --status
wsl.exe --list --verbose
wsl.exe -d <approved-distribution> -- env GOTOOLCHAIN=local sh -lc 'uname -a; cat /etc/os-release; git --version; docker version; go version; go env GOTOOLCHAIN'
wsl.exe -d <approved-distribution> -- env GOTOOLCHAIN=local sh -lc 'test "$(go version | cut -d" " -f3)" = go1.27.2 && test "$(go env GOTOOLCHAIN)" = local'
~~~

Record the VM's approved Linux architecture separately. Any version mismatch or missing local tool blocks this leg; do not install it. Verify that the WSL distro is isolated from Windows drives, interop, host credentials and unrelated distros before any act invocation.

~~~sh
# In the approved WSL2 distribution; use only pinned Linux binaries and paths.
set -eu
export GOTOOLCHAIN=local
RUN_ID='<approved-unique-run-id>'
ACT='<absolute-path-to-reviewed-act-binary>'
WORKFLOW='<absolute-path-to-reviewed-workflow-at-pinned-tree>'
EVENT='<absolute-path-to-reviewed-event-json>'
EVENT_NAME='<push|pull_request|workflow_dispatch>'
IMAGE='<registry/repository@sha256:platform-manifest-digest>'

"$ACT" --version
"$ACT" --help
# Only after WSL, daemon, outer policy/capture and safety controls independently PASS:
set +e
"$ACT" -W "$WORKFLOW" -e "$EVENT" -P "ubuntu-22.04=$IMAGE" \
  -s "GITHUB_TOKEN=CISEE-DUMMY-$RUN_ID" "$EVENT_NAME"
ACT_EXIT=$?
set -e
printf 'act_exit=%s run_id=%s\n' "$ACT_EXIT" "$RUN_ID"
~~~

Do not use privileged mode, added capabilities, host namespaces, container options that relax the reviewed profile, unpinned images, gh auth token, real secret files, or host credential forwarding. Do not add action-offline-mode as a substitute for network enforcement; test acquisition and runtime egress separately.

### Network controls and Docker inspection

Run approved and prohibited endpoints through reviewed probe workflows in each boundary, using fixture actions that report only exit/result and run ID. A generic shell probe is a control for that shell only. In WSL2, require curl from the approved local tool inventory; use a three-second connect timeout and five-second total timeout against the approved /health endpoint (must succeed), then the prohibited endpoint (must fail). Repeat with explicit IPv4, IPv6, configured proxy and proxy bypass, and DNS cases. In PowerShell, use Test-NetConnection for the same endpoint/port pairs and assert TcpTestSucceeded is true for approved and false for prohibited. An absent IPv6 route counts only when enforcement configuration and independent observation establish the intended deny outcome.

~~~sh
curl --connect-timeout 3 --max-time 5 --fail https://$APPROVED_HOST/health
if curl --connect-timeout 3 --max-time 5 --fail https://$PROHIBITED_HOST/health; then
  printf 'FAIL: prohibited destination succeeded\\n' >&2
  exit 1
fi
~~~

The Windows-native and WSL2 act recipes cover the matrix's reviewed push, pull_request, and workflow_dispatch fixture invocations by selecting the matching pinned event name and event JSON. Use distinct event/ref fixtures for ACT-PUSH-CONTEXT-CHECKOUT, ACT-PULL_REQUEST-CONTEXT-CHECKOUT, and ACT-WORKFLOW_DISPATCH-CONTEXT-CHECKOUT; include the missing-input variant for dispatch. They also provide the raw engine observations for the ACT-FAIL-CLOSED-* rows. Runtime inspection and endpoint probes apply to ACT-INTERNAL-DOCKER-NETWORK, ACT-LABELLED-CONTAINERS-INSPECTED, ACT-CACHED-CHECKOUT-RESOLUTION, and ACT-CREDENTIAL-ISOLATION-RUNTIME-AUDIT. Use the environment gate plus network-control recipes for ACT-CREDENTIAL-ISOLATION-PREFLIGHT, ACT-HOST-OUTBOUND-NETWORK-ISOLATION, ACT-EXECUTION-BLOCKED, ACT-METADATA-VERSION-HELP, ACT-FIXTURE-SEPARATION, and ACT-PRIVATE-INTERNAL-GIT-TRANSPORT. The Git hook fixture directly covers GIT-HOOK-BYPASS-NO-VERIFY; Git tree/object and event-manifest comparisons from each case cover the PASS-REUSE-* rows.

After each fixture, inspect only IDs resolved from daemon create-event evidence tied to the run ID, such as docker ps with the unique run label where the pinned runtime actually applies that label, then docker inspect with a format limited to HostConfig, Mounts, and NetworkSettings.Networks. Record daemon/API config, mounts, ports, HostConfig, security options, image manifest, networks and namespace data separately from container probe output and outside-guest capture. These commands are inventory aids, not proof by themselves. If the selected daemon or VM platform lacks equivalent read-only policy and telemetry commands, record an environmental/architecture blocker; do not install an agent or change host policy.


### Git identity and event oracle commands

Use these against only the disposable fixture checkout. They derive the staged and committed trees from Git itself; compare the output and SHA-256 of each event/workflow/config fixture plus pinned act and runner image identities. Do not store plaintext secret values; record only approved dummy-canary identity or the canonical secret digest.

~~~powershell
$FixtureRepo = '<absolute-path-to-disposable-fixture-repo>'
$EventFile = '<absolute-path-to-pinned-event-json>'
$WorkflowFile = '<absolute-path-to-pinned-workflow>'
& git.exe -C $FixtureRepo status --porcelain=v2 --branch
$StagedTree = & git.exe -C $FixtureRepo write-tree
$CommittedTree = & git.exe -C $FixtureRepo rev-parse 'HEAD^{tree}'
$EventHash = (Get-FileHash -Algorithm SHA256 -LiteralPath $EventFile).Hash
$WorkflowHash = (Get-FileHash -Algorithm SHA256 -LiteralPath $WorkflowFile).Hash
[pscustomobject]@{ StagedTree=$StagedTree; CommittedTree=$CommittedTree; EventHash=$EventHash; WorkflowHash=$WorkflowHash }
~~~

~~~sh
FIXTURE_REPO='<absolute-path-to-disposable-fixture-repo>'
EVENT_JSON='<absolute-path-to-pinned-event-json>'
WORKFLOW_FILE='<absolute-path-to-pinned-workflow>'
git -C "$FIXTURE_REPO" status --porcelain=v2 --branch
git -C "$FIXTURE_REPO" write-tree
git -C "$FIXTURE_REPO" rev-parse 'HEAD^{tree}'
sha256sum "$EVENT_JSON" "$WORKFLOW_FILE"
~~~

Capture these outputs immediately before and after each controlled case. Exact tree plus complete equal fingerprint is the only candidate equivalence case; change one named dimension per negative case and require a miss. These are feasibility oracles only; CI See PASS binding remains OPEN until S5 implementation tests prove exact matching.

### Safe Git hook fixture (no production repository)

Run the following independently with Windows-native Git and WSL2 Git. Use a new empty fixture path for each leg.

~~~powershell
$Fixture = Join-Path $env:TEMP "ci-see-hook-$RunId"
$Hooks = Join-Path $Fixture 'hooks'
New-Item -ItemType Directory -Force -Path $Hooks | Out-Null
& git.exe init --quiet $Fixture
& git.exe -C $Fixture config user.name 'CI See fixture'
& git.exe -C $Fixture config user.email 'fixture@example.invalid'
[IO.File]::WriteAllText((Join-Path $Fixture 'input.txt'), "fixture`n")
$Hook = "#!/bin/sh`nprintf 'ran\\n' > '$Fixture/hook-ran'; exit 73`n"
[IO.File]::WriteAllText((Join-Path $Hooks 'pre-commit'), $Hook, [Text.UTF8Encoding]::new($false))
& git.exe -C $Fixture add input.txt
& git.exe -C $Fixture -c "core.hooksPath=$Hooks" commit -m should-fail
$NormalExit = $LASTEXITCODE
if ($NormalExit -eq 0) { throw 'FAIL: ordinary commit bypassed failing hook' }
if (& git.exe -C $Fixture rev-parse --verify HEAD 2>$null) { throw 'FAIL: ordinary commit created HEAD' }
if (-not (Test-Path (Join-Path $Fixture 'hook-ran'))) { throw 'FAIL: hook was not observed' }
Remove-Item -LiteralPath (Join-Path $Fixture 'hook-ran')
& git.exe -C $Fixture -c "core.hooksPath=$Hooks" commit --no-verify -m fixture
$BypassExit = $LASTEXITCODE
if ($BypassExit -ne 0) { throw 'FAIL: expected native no-verify commit to succeed' }
if (Test-Path (Join-Path $Fixture 'hook-ran')) { throw 'FAIL: hook ran under no-verify' }
& git.exe -C $Fixture cat-file -p HEAD
~~~

In WSL2, run the same fixture in another fresh directory:

~~~sh
set -eu
ROOT=$(mktemp -d)
REPO="$ROOT/repo"
HOOKS="$ROOT/hooks"
mkdir -p "$REPO" "$HOOKS"
cd "$REPO"
git init -q
git config user.name 'CI See fixture'
git config user.email 'fixture@example.invalid'
printf 'fixture\\n' > input.txt
printf '%s\\n' '#!/bin/sh' 'printf "ran\\n" > hook-ran' 'exit 73' > "$HOOKS/pre-commit"
chmod 700 "$HOOKS/pre-commit"
export FIXTURE_MARKER="$REPO/hook-ran"
git add input.txt
set +e
git -c core.hooksPath="$HOOKS" commit -m should-fail
NORMAL_EXIT=$?
set -e
test "$NORMAL_EXIT" -ne 0
test -f "$FIXTURE_MARKER"
if git rev-parse --verify HEAD >/dev/null 2>&1; then exit 1; fi
rm "$FIXTURE_MARKER"
git -c core.hooksPath="$HOOKS" commit --no-verify -m fixture
test -f "$REPO/.git/HEAD"
test ! -e "$FIXTURE_MARKER"
test "$(git rev-list --count HEAD)" -eq 1
git cat-file -p HEAD
sha256sum "$HOOKS/pre-commit"
~~~

Expected: normal commit fails with hook marker present and no HEAD; after marker removal, the no-verify commit succeeds, marker stays absent, and one fixture commit exists. Preserve transcript, exit statuses, hook bytes/hash, commit object and marker state. This proves Git native bypass semantics only; Phase A does not execute it and it cannot establish CI See enforcement. If the fixture cannot run safely in an empty disposable repository, retain this ID as UNVERIFIED. First run a normal commit and require nonzero exit with no commit object; then remove the marker and run git -c core.hooksPath=<fixture-hooks> commit --no-verify -m fixture, requiring exit 0, one new commit, and no hook marker. Preserve transcript, exit statuses, hook bytes/hash, git cat-file -p HEAD, and marker state. Run the fixture independently with Windows-native Git and WSL2 Git. This is Phase B feasibility only and cannot establish CI See enforcement. If it cannot run safely in a disposable repository, retain this evidence ID as UNVERIFIED.

## Phase A checks and unresolved decisions

Phase A checks are documentation-only: validate links/fragments, index reachability, exact coverage of the 26 IDs, uniqueness, fixture references, required matrix fields, classification/phase consistency, command-source links, and preservation of the D3 record’s historical 26 PASS / 26 UNVERIFIED status. Existing documentation validators are the applicable repository checks. The Phase B command catalogue is not executed in Phase A. No Go tests, real act, workflow execution, firewall modification, Docker configuration change, host-tool download, or unreviewed dependency execution is authorized in Phase A.

The following decisions remain open for independent review and owner resolution before Phase B:

1. VM provider/capability evidence and whether both Windows-native and WSL2 legs are available.
2. Exact approved and prohibited destinations per test stage, including DNS, proxies, registry/action acquisition, runtime egress and internal Git transport.
3. Exact outside-guest deny-by-default enforcement and capture mechanism for the selected VM; independent instrumentation for native act, WSL2, children/helpers, Docker daemon, and each container class. If unavailable, Phase B remains BLOCKED.
4. Trusted action source policy and per-platform OCI digest provenance; the historical runner artifact is not independently available for rehashing.
5. Credential-store access tracing, daemon socket/API reachability, named-pipe, host filesystem/WSL interop mount, privilege/capability/namespace, port-publishing, and residual risks.
6. Meaningful commit-to-push event equivalence dimensions and treatment of hook mutation/caller behavior where canonical requirements remain unresolved.
7. Abort authority, bounded-stop behavior, hypervisor suspend/restore/discard process, post-stop listener/resource oracle, and private evidence custody/sanitization reviewer/publication approval.

**Phase B readiness: NO** until independent security review is PASS, architecture decisions affecting the test boundary are resolved, and the owner explicitly approves the exact execution scope. Even after Phase B, S4/S5 product acceptance stays open until the implemented product passes its own tests. D3 and production status remain unchanged unless their separate acceptance gates are met.
