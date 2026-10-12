# Issue #17 — Phase B runtime command catalogue

**Status:** Proposed Phase B procedure only; no command here was run in Phase A.  
**Authority:** [D3 Phase A plan](d3-phase-a-verification-plan.md), [local CI specification](../specs/local-ci-execution.md), [Context policy](../context-policy.md). Check every option against the exact pinned act release and retain its release-tagged source/help with the run manifest. Upstream references: [act command flags](https://github.com/nektos/act/blob/4f411281417e88660bea1c1a1749aa71ae0bd60f/cmd/root.go), [configuration behavior](https://nektosact.com/usage/index.html#configuration-file), [socket-disable example](https://github.com/nektos/act/issues/2239).

This catalogue is not approval to execute. Phase B remains blocked until independent security PASS, architecture decisions resolved, and explicit owner approval names the VM, workflows, endpoints, and commands. If the pinned binary/runtime cannot honor any gate below, stop before workflow execution.

**Growth decision: KEEP.** This procedure is one operational unit: environment and snapshot gates, both launch legs, independent controls, inspection, and emergency stop/rollback must be read together to run safely. Splitting those dependencies further would make the commands easier to invoke without their governing gates.

## Required gates

1. Use a disposable Windows 11 VM with documented nested-virtualization and WSL2 support. Validate Windows-native and WSL2 legs independently. Require local Go 1.27.2 with `GOTOOLCHAIN=local`; missing tools or capabilities are environmental blockers. Do not install or download tools during restricted verification.
2. Use a dedicated disposable Docker daemon confined to the VM. The approved host `act` process alone may reach its endpoint. Never expose a daemon socket, Windows pipe, or daemon credential to job, service, action, or helper containers.
3. Establish outside-guest deny-by-default enforcement and independent capture before any act network attempt. Verify the effective policy with non-act controls. Prepare independent policy inspection and connection observations for host act, children/helpers, WSL2, daemon, and each container class.
4. Review and pin the guard workflow. Its only job step must be trusted inline probe code; no third-party action, service, checkout, or unreviewed command. The first act invocation may run only this guard workflow. It must test prohibited Docker API/socket and host-service paths and any explicitly approved positive endpoint. Independently inspect the resulting container HostConfig/mounts and correlate daemon events before permitting any other workflow. A failed or unavailable observation is NO-GO.
5. Create a fresh detached fixture snapshot from the approved commit/tree. Independently compare its complete file manifest and Git tree to the reviewed source. Record hashes before each invocation and after it. A mismatch, dirty state, changed workflow, or inability to compute the oracle blocks execution.
6. Run act from a new empty invocation directory with a newly created empty home, XDG config, Docker config, and run temp directory. Inventory these directories and the snapshot for `.actrc`, `.env`, `.secrets`, `.vars`, and `.input`. Require those ambient files absent. Pass explicit, hashed empty files to all four file flags. Clear the child environment and pass only an allowlist, including isolated HOME/USERPROFILE/XDG/APPDATA/LOCALAPPDATA, empty DOCKER_CONFIG/GH_CONFIG_DIR, GIT_CONFIG_NOSYSTEM=1 and a hashed empty GIT_CONFIG_GLOBAL. Inventory Git credential helpers, SSH agent/key sources, cloud credentials, and Docker credential stores without exposing values; they must be absent/inaccessible. Set GOTOOLCHAIN=local. Any inherited variable/config or credential source that cannot be excluded blocks invocation. If this isolation cannot be demonstrated, do not invoke act.

The reviewed upstream source sets `--network` to `host` by default; therefore pass the dedicated per-run network explicitly. It documents `--container-daemon-socket -` as disabling the socket bind mount. The release-tagged source and pinned binary help must confirm both before Phase B; otherwise NO-GO. The upstream configuration loader searches XDG, home, and the process invocation directory, so isolate all three. The exact act options below are required on every workflow invocation: `--directory "$Snapshot"`, `--network "$RunNetwork"` where the inspected per-run network is not `host`, and `--container-daemon-socket -`. Also pass explicit empty `--env-file`, `--secret-file`, `--var-file`, and `--input-file` paths plus the reviewed workflow/event/image. Verify these options exist in the pinned release; never improvise replacements. The flags select act behavior, while independent daemon inspection and process/container network observations prove effective behavior.

## Windows-native guard and fixture launch

The operator substitutes only values from the approved run manifest. `$InvocationDir` is empty and isolated; `$Snapshot` is the immutable reviewed checkout; `$Workflow` is a path inside it; `$RunNetwork` is a unique inspected non-host network; all `$Empty*` files are zero-byte files in the run temp directory.

```powershell
$ErrorActionPreference = 'Stop'
$RunId = '<approved-run-id>'
$ActExe = '<pinned-act.exe>'
$Snapshot = '<fresh-detached-snapshot>'
$InvocationDir = '<empty-run-owned-directory>'
$Workflow = '.github/workflows/security-guard.yml' # guard only on first invocation
$EventFile = '<pinned-event-json>'
$Image = '<registry/image@sha256:platform-manifest-digest>'
$RunNetwork = '<run-owned-bridge-network-name>'
$EmptyEnv = '<run-temp>/empty.env'
$EmptySecrets = '<run-temp>/empty.secrets'
$EmptyVars = '<run-temp>/empty.vars'
$EmptyInput = '<run-temp>/empty.input'
$DockerConfig = '<empty-run-owned-docker-config>'
$CleanHome = '<empty-run-owned-home>'
$XdgHome = '<empty-run-owned-xdg-config>'
$DaemonEndpoint = '<approved-disposable-VM-daemon-endpoint>'
$ActionCache = '<empty-run-owned-action-cache>'
$EmptyGitConfig = '<run-temp>/empty.gitconfig'
$GHConfig = '<empty-run-owned-gh-config>'

# Before invocation: verify exact commit/tree and full snapshot manifest; inspect
# HOME/XDG/invocation config locations and all ambient input names. Do not print values.
# Require all four explicit input files to be empty and hash-recorded.
$Arguments = @(
  '--directory', $Snapshot, '-W', $Workflow, '-e', $EventFile, '--action-cache-path', $ActionCache,
  '-P', "ubuntu-22.04=$Image", '--network', $RunNetwork,
  '--container-daemon-socket', '-',
  '--env-file', $EmptyEnv, '--secret-file', $EmptySecrets,
  '--var-file', $EmptyVars, '--input-file', $EmptyInput,
  '-s', "GITHUB_TOKEN=CISEE-DUMMY-$RunId", 'push'
)
if ($RunNetwork -eq 'host') { throw 'NO-GO: host network is forbidden' }
$psi = [Diagnostics.ProcessStartInfo]::new($ActExe)
$psi.WorkingDirectory = $InvocationDir
$psi.UseShellExecute = $false
$psi.RedirectStandardOutput = $true
$psi.RedirectStandardError = $true
$psi.Environment.Clear()
$Allow = @{
  PATH = '<pinned-system-and-tool-paths>'; SystemRoot = $env:SystemRoot
  TEMP = '<run-owned-temp>'; TMP = '<run-owned-temp>'
  HOME = $CleanHome; USERPROFILE = $CleanHome; XDG_CONFIG_HOME = $XdgHome
  APPDATA = '<empty-run-owned-appdata>'; LOCALAPPDATA = '<empty-run-owned-localappdata>'
  DOCKER_CONFIG = $DockerConfig; DOCKER_HOST = $DaemonEndpoint; GH_CONFIG_DIR = $GHConfig
  GIT_CONFIG_NOSYSTEM = '1'; GIT_CONFIG_GLOBAL = $EmptyGitConfig; GIT_TERMINAL_PROMPT = '0'
  GOTOOLCHAIN = 'local'
}
foreach ($entry in $Allow.GetEnumerator()) { $psi.Environment[$entry.Key] = $entry.Value }
foreach ($arg in $Arguments) { [void]$psi.ArgumentList.Add($arg) }
# First invocation only: reviewed guard workflow, never a target/untrusted workflow.
$proc = [Diagnostics.Process]::Start($psi)
$outTask = $proc.StandardOutput.ReadToEndAsync()
$errTask = $proc.StandardError.ReadToEndAsync()
$proc.WaitForExit()
$stdout = $outTask.GetAwaiter().GetResult()
$stderr = $errTask.GetAwaiter().GetResult()
$ActExit = $proc.ExitCode
# Store raw output privately; inspect only daemon-event IDs owned by $RunId.
```

Before a later reviewed fixture, independently approve the guard result and inspect actual daemon HostConfig, mounts, network mode/routes, ports, privilege/capabilities, security options, image digest, and namespace data. Confirm no socket/pipe or credential mount and no host service is reachable. Preserve the configured policy and independent observations separately. Any mismatch, hidden inherited config, unknown mount, host mode, or reachable prohibited endpoint leaves isolation INCOMPLETE and prohibits the next invocation.

For each later case, change only the approved workflow/event from the manifest; retain the same explicit safety options, clean environment, `--directory` snapshot, empty files, isolated Git/credential configuration, and manifest-bound run ID. Do not use `--privileged`, container options that relax the profile, host namespaces, implicit action cache substitutions, `gh auth`, or `--no-skip-checkout`. Do not treat `--dryrun`, offline mode, image pinning, or an outer egress filter as proof of isolation.

## WSL2 leg

Use a dedicated approved distribution and Linux-native snapshot path, binaries, empty config/input files, and daemon endpoint. First record WSL/kernel/distribution identity and confirm WSL2 capabilities. Do not infer any result from Windows-native verification.

```sh
set -eu
RUN_ID='<approved-run-id>'
ACT='<pinned-act-linux-binary>'
SNAPSHOT='<fresh-detached-snapshot>'
INVOCATION_DIR='<empty-run-owned-directory>'
WORKFLOW='.github/workflows/security-guard.yml' # guard only on first invocation
EVENT='<pinned-event-json>'
IMAGE='<registry/image@sha256:platform-manifest-digest>'
RUN_NETWORK='<run-owned-bridge-network-name>'
ACTION_CACHE='<empty-run-owned-action-cache>'
EMPTY_GIT_CONFIG='<run-temp>/empty.gitconfig'
GH_CONFIG_DIR='<empty-run-owned-gh-config>'
EMPTY_ENV='<run-temp>/empty.env'
EMPTY_SECRETS='<run-temp>/empty.secrets'
EMPTY_VARS='<run-temp>/empty.vars'
EMPTY_INPUT='<run-temp>/empty.input'
CLEAN_HOME='<empty-run-owned-home>'
XDG_HOME='<empty-run-owned-xdg-config>'
DOCKER_CONFIG='<empty-run-owned-docker-config>'
DOCKER_HOST='<approved-disposable-VM-daemon-endpoint>'
test "$RUN_NETWORK" != host
# Before launch, verify the full snapshot/tree manifest and absent ambient files.
# env -i prevents inherited secrets/config; use only manifest-approved values.
cd "$INVOCATION_DIR"
env -i PATH='<pinned-tool-path>' HOME="$CLEAN_HOME" XDG_CONFIG_HOME="$XDG_HOME" \
  DOCKER_CONFIG="$DOCKER_CONFIG" DOCKER_HOST="$DOCKER_HOST" GOTOOLCHAIN=local \
  "$ACT" --directory "$SNAPSHOT" -W "$WORKFLOW" -e "$EVENT" \
  -P "ubuntu-22.04=$IMAGE" --network "$RUN_NETWORK" \
  --container-daemon-socket - \
  --env-file "$EMPTY_ENV" --secret-file "$EMPTY_SECRETS" \
  --var-file "$EMPTY_VARS" --input-file "$EMPTY_INPUT" \
  -s "GITHUB_TOKEN=CISEE-DUMMY-$RUN_ID" push
```

For the later fixture invocation, select only an individually reviewed workflow/event from the manifest. Retain the same required options and environment clearing. Store output privately, capture exit status, and inspect only run-owned daemon resources. If the exact pinned version rejects an option or any environment/config source remains uncontrolled, stop and classify the leg BLOCKED.

## Reproducible endpoint probes

For a reviewed guard container with pinned curl, use controlled endpoints from the approved manifest. Run the commands inside the guard's only trusted inline step; the approved endpoint must succeed and the prohibited endpoint must fail. Repeat with the approved IPv4 and IPv6 addresses, DNS, configured proxy and proxy bypass, and approved/prohibited host-service destinations. Missing protocol support or an unobservable denial is INCOMPLETE.

```sh
set -eu
curl --connect-timeout 3 --max-time 5 --fail "$APPROVED_URL" >/dev/null
if curl --connect-timeout 3 --max-time 5 --fail "$PROHIBITED_URL" >/dev/null; then
  echo 'FAIL: prohibited destination succeeded' >&2
  exit 1
fi
curl --ipv4 --connect-timeout 3 --max-time 5 --fail "$APPROVED_IPV4_URL" >/dev/null
if curl --ipv4 --connect-timeout 3 --max-time 5 --fail "$PROHIBITED_IPV4_URL" >/dev/null; then exit 1; fi
curl --ipv6 --connect-timeout 3 --max-time 5 --fail "$APPROVED_IPV6_URL" >/dev/null
if curl --ipv6 --connect-timeout 3 --max-time 5 --fail "$PROHIBITED_IPV6_URL" >/dev/null; then exit 1; fi
if test -S /var/run/docker.sock && curl --max-time 3 --unix-socket /var/run/docker.sock http://localhost/version; then
  echo 'FAIL: Docker API reachable in guard container' >&2
  exit 1
fi
```

The control values, tool version, policy, and independent sink/capture are manifest-bound. A separate PowerShell `Test-NetConnection` or shell `curl` invocation is only a transport control for that process. To claim host-act, helper, or daemon isolation, generate the approved and prohibited attempts through that exact component's reviewed operation and attribute them to its PID/process tree or daemon; if no safe controlled operation exists, mark that boundary INCOMPLETE. Never substitute a container or operator-shell result for another boundary.

## Independent network controls and inspection

Every applicable boundary needs one approved positive control and one prohibited negative control generated by that boundary. Exercise IPv4 and IPv6, DNS, configured proxy and bypass, loopback, bridge/host-gateway, daemon endpoint, image pulls, and action acquisition as separate cases. Record intended policy, effective enforcement configuration, and independent observation as separate evidence. A generic shell test is only evidence about that shell.

| Boundary | Positive control | Negative control and independent oracle |
| --- | --- | --- |
| act host and each child/helper class | Reach an explicitly approved destination for that phase. | Prohibited destination must fail; collect process-attributed socket/connection trace and effective per-process policy. |
| Docker daemon | Pull only the approved pinned image during the acquisition phase. | Unapproved registry/DNS/proxy must fail; compare daemon policy/logs with independent boundary capture and verified digest. |
| Job and service containers | Reach only an approved test endpoint/peer allowed by the fixture. | Prohibited egress, host service, loopback/bridge peer and daemon API probes must fail; compare probe output with daemon HostConfig/routes and independent observations. |
| WSL2 and acquisition helpers | Approved controls only where destination policy permits. | Prohibited IPv4/IPv6/DNS/proxy route must fail; attribute to WSL process/helper and reconcile outside-guest capture. |

A missing positive control, failed prohibited control, unobserved route, or unverifiable policy yields **INCOMPLETE**, even if the other checks pass. The external VM boundary does not prove loopback or daemon API isolation.

After the guard/fixture, use run-owned daemon create-event IDs to select resources. Inspect only those IDs with the pinned runtime's read-only equivalent of `docker inspect`; retain HostConfig, Mounts, NetworkSettings, image digest, ports, security options, and labels. Correlate process tree, independent packet/connection observations, and VM listeners. If the platform cannot expose equivalent independent evidence, record an environmental/architecture blocker; do not install an agent or change host policy.

## Stop, rollback, and evidence custody

Stop before acquisition if approval, VM restore point, enforcement, capture, identity, runtime profile, or resource ownership evidence is absent. During a leg, stop on a prohibited connection succeeding, policy drift, unapproved content, credential/helper access, timeout, leaked canary, unexpected mount/privilege/network/port, missing telemetry, or failed cancellation. Mark the leg BLOCKED/INCOMPLETE; do not retry through another path.

Request bounded graceful cancellation, then terminate only verified run descendants after checking process identity/PID reuse. Remove only resources independently tied to the unique run ID. Compare process/container/network/volume/temp/listener inventories before and after. If ownership is uncertain, do not delete; quarantine the VM. If a child escaped, credentials may have leaked, policy cannot be restored, evidence is inconclusive, or observations disagree, disconnect/suspend through the owner-controlled hypervisor and restore a known snapshot or discard the VM. Never clean unrelated or physical-host resources.

Seal raw evidence privately with hashes and custody metadata. Publish only after independent sanitization review. Failed cleanup, an unaccounted resource, lost custody, or missing post-stop listener checks leaves isolation INCOMPLETE.
