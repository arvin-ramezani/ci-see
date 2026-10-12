# Issue #17 — Phase B runtime command catalogue

**Status:** Proposed Phase B procedure only; no command here was run in Phase A.  
**Authority:** [D3 Phase A plan](d3-phase-a-verification-plan.md), [local CI specification](../specs/local-ci-execution.md), [Context policy](../context-policy.md). Check every option against the exact pinned act release and retain its release-tagged source/help with the run manifest. Upstream references: [act command flags](https://github.com/nektos/act/blob/4f411281417e88660bea1c1a1749aa71ae0bd60f/cmd/root.go), [configuration behavior](https://nektosact.com/usage/index.html#configuration-file), [socket-disable example](https://github.com/nektos/act/issues/2239).

This catalogue is not approval to execute. Phase B remains blocked until independent security PASS, architecture decisions resolved, and explicit owner approval names the VM, workflows, endpoints, and commands. If the pinned binary/runtime cannot honor any gate below, stop before workflow execution.

**Growth decision: KEEP (>300 lines).** This remains one gated operational procedure: both platform launch examples depend on the shared environment/snapshot gates and the later independent controls, inspection, and rollback steps. Neither launch block is a safe independent retrieval unit; separating it would detach executable commands from the gates they require. The cohesive exception is retained under the repository context policy.

## Required gates

1. Use a disposable Windows 11 VM with documented nested-virtualization and WSL2 support. Validate Windows-native and WSL2 legs independently. Require local Go 1.27.2 with `GOTOOLCHAIN=local`; missing tools or capabilities are environmental blockers. Do not install or download tools during restricted verification.
2. Use a dedicated disposable Docker daemon confined to the VM. The approved host `act` process alone may reach its endpoint. Never expose a daemon socket, Windows pipe, or daemon credential to job, service, action, or helper containers.
3. Establish outside-guest deny-by-default enforcement and independent capture before any act network attempt. Verify the effective policy with non-act controls. Prepare independent policy inspection and connection observations for host act, children/helpers, WSL2, daemon, and each container class.
4. Review and pin the guard workflow. Its only job step must be trusted inline probe code; no third-party action, service, checkout, or unreviewed command. The first act invocation may run only this guard workflow. It must test prohibited Docker API/socket and host-service paths and any explicitly approved positive endpoint. Independently inspect the resulting container HostConfig/mounts and correlate daemon events before permitting any other workflow. A failed or unavailable observation is NO-GO.
5. Create a fresh detached fixture snapshot from the approved commit/tree. Independently compare its complete file manifest and Git tree to the reviewed source. Record hashes before each invocation and after it. A mismatch, dirty state, changed workflow, or inability to compute the oracle blocks execution.
6. Run act from a new empty invocation directory with a newly created empty home, XDG config, Docker config, and run temp directory. Inventory these directories and the snapshot for `.actrc`, `.env`, `.secrets`, `.vars`, and `.input`. Require those ambient files absent. Pass explicit, hashed empty files to all four file flags. Clear the child environment and pass only an allowlist, including isolated HOME/USERPROFILE/XDG/APPDATA/LOCALAPPDATA, empty DOCKER_CONFIG/GH_CONFIG_DIR, GIT_CONFIG_NOSYSTEM=1, a hashed empty GIT_CONFIG_GLOBAL, GIT_TERMINAL_PROMPT=0, and the same command-scope empty credential.helper override in both legs. Inventory Git credential-helper and HTTP extraheader sources by key/origin only; reject file-based matches. Exclude GIT_ASKPASS, SSH_AUTH_SOCK, SSH keys/agents, cloud credential variables, and credential stores from each allowlist; they must be absent/inaccessible. Set GOTOOLCHAIN=local. Require finite per-leg deadlines, captured stdout/stderr/exit status, bounded normal-stop then verified-tree forced-stop handling, and post-stop run-owned resource/listener checks. Any inherited variable/config or credential source that cannot be excluded blocks invocation. If this isolation or recovery cannot be demonstrated, do not invoke act.

The reviewed pinned source defaults --network to host, --pull to true, and --rebuild to true. It documents --container-daemon-socket - as disabling the socket bind mount. Every workflow invocation below must pass --directory "$Snapshot", the inspected non-host --network "$RunNetwork", --container-daemon-socket -, --pull=false, --rebuild=false, and --action-offline-mode, plus explicit empty --env-file, --secret-file, --var-file, and --input-file paths. Confirm every option against the exact pinned binary's help and immutable source before Phase B; any mismatch is NO-GO. Flags select act behavior, while independent daemon inspection and process/container network observations prove effective behavior.

### Dependency acquisition boundary

Acquire dependencies only in a distinct, owner-approved stage with its own allowlisted destinations, deny-by-default enforcement, positive/negative controls, and independent capture. The execution leg must consume a read-only or immutable handoff. Pin and verify the platform-specific OCI manifest digest and each action's immutable repository commit/content hash before launch; record the complete action-cache inventory. The pinned source says --pull and --rebuild default to true, so --pull=false --rebuild=false is mandatory even when a digest is supplied. --action-offline-mode is an additional cache behavior, not an isolation boundary: preflight every referenced action and image, and block before invocation if any content is missing or mismatched. Do not permit an act fallback pull, rebuild, action fetch, or alternate registry. A blocked attempt or a successful local-cache run does not prove egress isolation; retain enforcement configuration and independent observations separately.

## Windows-native guard and fixture launch

The operator substitutes values only from the approved run manifest. $InvocationDir is empty and isolated; $Snapshot is the immutable reviewed checkout; $Workflow is inside it; $RunNetwork is a unique inspected non-host network; and all $Empty* files are zero-byte files in the run temp directory. The manifest must approve finite limits; the examples use a 300-second act deadline, a 15-second normal-termination grace period, and a 15-second forced-stop settlement period. Require a pinned PowerShell version and independently verified process-tree observer before launch. If the observer cannot attribute descendants to this run, do not start act. Both platforms must use the same Git policy: disable system config, use a hash-verified empty global config, disable prompts, and reset credential.helper at command scope. Inventory repository-local helper and http.*.extraheader key origins without values; any file-based match blocks the leg.

```powershell
$ErrorActionPreference = 'Stop'
$RunId = '<approved-run-id>'
$ActExe = '<pinned-act.exe>'
$ExpectedActPath = [IO.Path]::GetFullPath($ActExe)
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
$ActionCache = '<read-only-verified-run-action-cache>'
$EmptyGitConfig = '<run-temp>/empty.gitconfig'
$GHConfig = '<empty-run-owned-gh-config>'
$PrivateEvidence = '<private-run-evidence-directory>'
$TimeoutSeconds = 300
$GraceSeconds = 15
$ForcedSettlementSeconds = 15

# Before invocation: verify exact commit/tree, full snapshot manifest, dependency
# digests/cache entries, isolated config locations, and absent ambient input files.
# Hash-record the empty files and private evidence directory; never print secrets.
$Arguments = @(
  '--directory', $Snapshot, '-W', $Workflow, '-e', $EventFile, '--action-cache-path', $ActionCache,
  '-P', "ubuntu-22.04=$Image", '--network', $RunNetwork,
  '--container-daemon-socket', '-',
  '--pull=false', '--rebuild=false', '--action-offline-mode',
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
  GIT_CONFIG_COUNT = '1'; GIT_CONFIG_KEY_0 = 'credential.helper'; GIT_CONFIG_VALUE_0 = ''
  GOTOOLCHAIN = 'local'
}
foreach ($entry in $Allow.GetEnumerator()) { $psi.Environment[$entry.Key] = $entry.Value }

# Run a read-only preflight with the same isolated HOME and no-system/empty-global policy as the act child; omit the command-scope helper reset here so file-based helper sources remain visible.
# Preserve and restore this PowerShell process environment; write names/origins only.
$GitEnvNames = @('HOME','USERPROFILE','XDG_CONFIG_HOME','APPDATA','LOCALAPPDATA',
  'GIT_CONFIG_NOSYSTEM','GIT_CONFIG_GLOBAL','GIT_TERMINAL_PROMPT')
$PriorGitEnv = @{}
foreach ($name in $GitEnvNames) {
  $PriorGitEnv[$name] = [Environment]::GetEnvironmentVariable($name, 'Process')
  [Environment]::SetEnvironmentVariable($name, $Allow[$name], 'Process')
}
try {
  $PriorConfigEnv = @{}
  foreach ($entry in Get-ChildItem Env: | Where-Object { $_.Name -match '^GIT_CONFIG_(COUNT|KEY_[0-9]+|VALUE_[0-9]+|PARAMETERS)$' }) {
    $PriorConfigEnv[$entry.Name] = $entry.Value
    [Environment]::SetEnvironmentVariable($entry.Name, $null, 'Process')
  }
  $GitKeys = & git.exe -C $Snapshot config --show-origin --name-only --get-regexp '^(credential\.helper|http\..*extraheader$)' 2>$null
  $GitConfigExit = $LASTEXITCODE
  $GitKeys | Set-Content -LiteralPath "$PrivateEvidence\git-config-keys.txt"
} finally {
  foreach ($name in $PriorConfigEnv.Keys) {
    [Environment]::SetEnvironmentVariable($name, $PriorConfigEnv[$name], 'Process')
  }
  foreach ($name in $GitEnvNames) {
    [Environment]::SetEnvironmentVariable($name, $PriorGitEnv[$name], 'Process')
  }
}
if ($GitConfigExit -eq 0) { throw 'NO-GO: credential helper or HTTP extraheader source is present' }
if ($GitConfigExit -ne 1) { throw 'NO-GO: Git credential-source inspection failed' }
foreach ($arg in $Arguments) { [void]$psi.ArgumentList.Add($arg) }

# The guard workflow is the only first invocation. Capture exact process identity,
# start time and descendants with the independent observer before waiting.
$proc = [Diagnostics.Process]::Start($psi)
$StartedAt = $proc.StartTime.ToUniversalTime()
$ProcessId = $proc.Id
$outTask = $proc.StandardOutput.ReadToEndAsync()
$errTask = $proc.StandardError.ReadToEndAsync()
$TimedOut = -not $proc.WaitForExit($TimeoutSeconds * 1000)
$StopRequestExit = $null
$ForcedStopUsed = $false
if ($TimedOut) {
  # Normal termination request to this run's verified process tree, then a bounded
  # grace period. taskkill without /F is the normal request; record its exit status.
  $current = Get-Process -Id $ProcessId -ErrorAction SilentlyContinue
  if ($current -and $current.StartTime.ToUniversalTime() -eq $StartedAt -and
      $current.Path -and [IO.Path]::GetFullPath($current.Path) -eq $ExpectedActPath) {
    & "$env:SystemRoot\System32\taskkill.exe" /PID $ProcessId /T
    $StopRequestExit = $LASTEXITCODE
    if (-not $proc.WaitForExit($GraceSeconds * 1000)) {
      # Recheck PID/start-time/image against the observer's run manifest before force.
      $current = Get-Process -Id $ProcessId -ErrorAction SilentlyContinue
      if (-not $current -or $current.StartTime.ToUniversalTime() -ne $StartedAt -or
          -not $current.Path -or [IO.Path]::GetFullPath($current.Path) -ne $ExpectedActPath) {
        throw 'UNRESOLVED: process identity changed; quarantine/restore the VM'
      }
      & "$env:SystemRoot\System32\taskkill.exe" /PID $ProcessId /T /F
      $ForcedStopUsed = $true
      if (-not $proc.WaitForExit($ForcedSettlementSeconds * 1000)) {
        throw 'UNRESOLVED: act did not settle; quarantine/restore the VM'
      }
    }
  } else {
    throw 'UNRESOLVED: act process identity unavailable; quarantine/restore the VM'
  }
}
if (-not $proc.WaitForExit($ForcedSettlementSeconds * 1000)) {
  throw 'UNRESOLVED: process exit not observed; quarantine/restore the VM'
}
if (-not $outTask.Wait($ForcedSettlementSeconds * 1000) -or
    -not $errTask.Wait($ForcedSettlementSeconds * 1000)) {
  throw 'UNRESOLVED: output capture did not close; quarantine/restore the VM'
}
$stdout = $outTask.GetAwaiter().GetResult()
$stderr = $errTask.GetAwaiter().GetResult()
$ActExit = $proc.ExitCode
$RunStatus = if ($TimedOut) { 'TIMEOUT' } elseif ($ActExit -eq 0) { 'EXIT_0' } else { 'NONZERO' }
[IO.File]::WriteAllText("$PrivateEvidence\act.stdout", $stdout)
[IO.File]::WriteAllText("$PrivateEvidence\act.stderr", $stderr)
[pscustomobject]@{
  runId=$RunId; processId=$ProcessId; startedAt=$StartedAt
  status=$RunStatus; actExit=$ActExit; stopRequestExit=$StopRequestExit
  forcedStop=$ForcedStopUsed; timeoutSeconds=$TimeoutSeconds
} | ConvertTo-Json | Set-Content -LiteralPath "$PrivateEvidence\act-result.json"
# Continue to bounded run-owned Docker cleanup and post-stop inventory even on
# TIMEOUT/NONZERO. Preserve evidence first; never remove resources by broad name.
```

A timeout or nonzero exit is recorded as TIMEOUT or NONZERO, never PASS. The Windows stop request and forced-stop path must be rehearsed with a harmless owned process before the approved run; if normal stop, descendant attribution, bounded wait, exit capture, or post-stop inventory is not demonstrated, classify the environment BLOCKED and do not invoke act. The run-owned cleanup and VM quarantine/restore procedure below applies after every status.

## WSL2 leg

Use a dedicated approved distribution and Linux-native snapshot path, binaries, empty config/input files, and daemon endpoint. First record WSL/kernel/distribution identity and confirm WSL2 capabilities. Do not infer any result from Windows-native verification. Require a locally available GNU coreutils timeout; verify its version and rehearse its process-group TERM/KILL behavior with a harmless owned process before any act invocation. Missing or unverified supervisor behavior is an environmental blocker.

```sh
set -u
PATH='<pinned-tool-path>'
export PATH
TIMEOUT_BIN='<hash-pinned-coreutils-timeout>'
TIMEOUT_SHA256_FILE='<private-run-evidence>/timeout.sha256'
test -x "$TIMEOUT_BIN" || { echo 'BLOCKED: pinned timeout missing' >&2; exit 2; }
RUN_ID='<approved-run-id>'
ACT='<pinned-act-linux-binary>'
SNAPSHOT='<fresh-detached-snapshot>'
INVOCATION_DIR='<empty-run-owned-directory>'
WORKFLOW='.github/workflows/security-guard.yml' # guard only on first invocation
EVENT='<pinned-event-json>'
IMAGE='<registry/image@sha256:platform-manifest-digest>'
RUN_NETWORK='<run-owned-bridge-network-name>'
ACTION_CACHE='<read-only-verified-run-action-cache>'
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
PRIVATE_EVIDENCE='<private-run-evidence-directory>'
sha256sum --check "$TIMEOUT_SHA256_FILE" > "$PRIVATE_EVIDENCE/timeout-sha256-check.txt" 2>&1 || { echo 'BLOCKED: timeout hash mismatch' >&2; exit 2; }
"$TIMEOUT_BIN" --version > "$PRIVATE_EVIDENCE/timeout-version.txt" 2>&1 || { echo 'BLOCKED: timeout version check failed' >&2; exit 2; }
TIMEOUT_SECONDS=300
GRACE_SECONDS=15
test "$RUN_NETWORK" != host || { echo 'NO-GO: host network' >&2; exit 2; }
cd "$INVOCATION_DIR"
# Files must be empty/hash-recorded. Inventory act config/default inputs and verify
# action-cache + image digests. Capture Git key names/origins only; never values.
set +e
env -i PATH='<pinned-tool-path>' HOME="$CLEAN_HOME" XDG_CONFIG_HOME="$XDG_HOME" LC_ALL=C \
  GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL="$EMPTY_GIT_CONFIG" GIT_TERMINAL_PROMPT=0 \
  git -C "$SNAPSHOT" config --show-origin --name-only --get-regexp '^(credential\.helper|http\..*extraheader$)' \
  > "$PRIVATE_EVIDENCE/git-config-keys.txt" 2> "$PRIVATE_EVIDENCE/git-config-inspection.stderr"
GIT_CONFIG_EXIT=$?
set -e
if [ "$GIT_CONFIG_EXIT" -eq 0 ]; then
  echo 'NO-GO: credential or HTTP extraheader config source is present' >&2
  exit 2
fi
if [ "$GIT_CONFIG_EXIT" -ne 1 ]; then
  echo 'NO-GO: Git credential-source inspection failed' >&2
  exit 2
fi
set -e

# No --foreground: GNU timeout signals the act process group, then escalates after
# GRACE_SECONDS. Pin locale so timeout diagnostics are parseable; capture status.
LC_ALL=C
INTERRUPTED=0
trap 'INTERRUPTED=1' INT TERM
set +e
"$TIMEOUT_BIN" --verbose --signal=TERM --kill-after="$GRACE_SECONDS"s "$TIMEOUT_SECONDS"s \
  env -i PATH='<pinned-tool-path>' HOME="$CLEAN_HOME" XDG_CONFIG_HOME="$XDG_HOME" \
    DOCKER_CONFIG="$DOCKER_CONFIG" DOCKER_HOST="$DOCKER_HOST" GH_CONFIG_DIR="$GH_CONFIG_DIR" \
    GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL="$EMPTY_GIT_CONFIG" GIT_TERMINAL_PROMPT=0 \
    GIT_CONFIG_COUNT=1 GIT_CONFIG_KEY_0=credential.helper GIT_CONFIG_VALUE_0= \
    GOTOOLCHAIN=local \
    "$ACT" --directory "$SNAPSHOT" -W "$WORKFLOW" -e "$EVENT" \
    -P "ubuntu-22.04=$IMAGE" --network "$RUN_NETWORK" \
    --container-daemon-socket - --pull=false --rebuild=false --action-offline-mode \
    --action-cache-path "$ACTION_CACHE" \
    --env-file "$EMPTY_ENV" --secret-file "$EMPTY_SECRETS" \
    --var-file "$EMPTY_VARS" --input-file "$EMPTY_INPUT" \
    -s "GITHUB_TOKEN=CISEE-DUMMY-$RUN_ID" push \
    > "$PRIVATE_EVIDENCE/act.stdout" 2> "$PRIVATE_EVIDENCE/act.stderr"
ACT_EXIT=$?
set -e
case "$ACT_EXIT" in
  0) RUN_STATUS=EXIT_0 ;;
  124|137)
    if grep -Fq 'sending signal TERM' "$PRIVATE_EVIDENCE/act.stderr" || grep -Fq 'sending signal KILL' "$PRIVATE_EVIDENCE/act.stderr"; then RUN_STATUS=TIMEOUT
    else RUN_STATUS=AMBIGUOUS_EXIT; fi ;;
  *) RUN_STATUS=NONZERO ;;
esac
if [ "$INTERRUPTED" -ne 0 ]; then RUN_STATUS=INTERRUPTED; fi
printf 'run_id=%s\nstatus=%s\nexit=%s\ntimeout_seconds=%s\ngrace_seconds=%s\n' \
  "$RUN_ID" "$RUN_STATUS" "$ACT_EXIT" "$TIMEOUT_SECONDS" "$GRACE_SECONDS" \
  > "$PRIVATE_EVIDENCE/act-result.txt"
# Continue to bounded run-owned Docker cleanup and post-stop inventory for every
# status. TIMEOUT/NONZERO/AMBIGUOUS_EXIT is never PASS; missing process/resource observations block.
```

The timeout wrapper's finite deadline, process-group TERM, and KILL-after grace are mandatory. Its exact local version and harmless-process rehearsal are evidence prerequisites, not an act test. If the wrapper itself is unavailable, its exit is ambiguous, output/status cannot be retained, or any descendant/resource remains unaccounted for, mark the leg BLOCKED/INCOMPLETE and quarantine/restore the VM. For every later fixture, keep the same environment, finite limits, options, verified cache, and run-owned recovery path; only the reviewed workflow/event may change. Do not proceed to another fixture unless RUN_STATUS is EXIT_0 and recovery/inventory is complete; EXIT_0 is only a raw feasibility observation, not product acceptance.

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
