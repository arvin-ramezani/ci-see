# Issue #17 — Phase B Git identity and hook fixtures

**Status:** Proposed Phase B feasibility procedure only; not executed in Phase A.  
**Authority:** [D3 Phase A plan](d3-phase-a-verification-plan.md), [Git gating specification](../specs/git-gating.md), [fixture catalog](fixtures.md).

Use a fresh disposable repository and isolated evidence directory for each Windows-native and WSL2 leg. Do not use the product repository. This fixture demonstrates native Git behavior only; it cannot prove CI See enforcement or make hooks tamper-proof.

## Tree and event identity oracle

Record source repository identity, immutable commit/tree, workflow/event hashes, runner image platform digest, act binary hash, and run fingerprint inputs without storing plaintext credentials. Capture before and after each controlled case.

```powershell
$FixtureRepo = '<absolute-path-to-disposable-fixture-repo>'
$EventFile = '<pinned-event-json>'
$WorkflowFile = '<pinned-workflow>'
& git.exe -C $FixtureRepo status --porcelain=v2 --branch
$StagedTree = & git.exe -C $FixtureRepo write-tree
$CommittedTree = & git.exe -C $FixtureRepo rev-parse 'HEAD^{tree}'
$EventHash = (Get-FileHash -Algorithm SHA256 -LiteralPath $EventFile).Hash
$WorkflowHash = (Get-FileHash -Algorithm SHA256 -LiteralPath $WorkflowFile).Hash
[pscustomobject]@{ StagedTree=$StagedTree; CommittedTree=$CommittedTree; EventHash=$EventHash; WorkflowHash=$WorkflowHash }
```

```sh
FIXTURE_REPO='<absolute-path-to-disposable-fixture-repo>'
EVENT_JSON='<pinned-event-json>'
WORKFLOW_FILE='<pinned-workflow>'
git -C "$FIXTURE_REPO" status --porcelain=v2 --branch
git -C "$FIXTURE_REPO" write-tree
git -C "$FIXTURE_REPO" rev-parse 'HEAD^{tree}'
sha256sum "$EVENT_JSON" "$WORKFLOW_FILE"
```

Derive the expected event/ref/commit/tree independently from Git objects and the approved manifest. The exact tree plus complete execution fingerprint is the only candidate equivalence case. Change one named input at a time for mismatch cases. Record all differing fields and outcomes. Feasibility observations never establish CI See PASS reuse; S5 matching remains OPEN until product acceptance passes.

## Safe disposable `--no-verify` fixture

Run independently with Windows-native Git and WSL2 Git. Create a new empty fixture repo and hooks directory for each leg. Configure identity locally, create one staged input file, and install this hook:

```sh
#!/bin/sh
printf 'ran\n' > hook-ran
exit 73
```

The marker is relative to the hook's repository working directory, so no native Windows path is embedded in a POSIX shell script.

### Windows-native Git

```powershell
$Fixture = Join-Path $env:TEMP "ci-see-hook-$RunId"
$Hooks = Join-Path $Fixture 'hooks'
New-Item -ItemType Directory -Force -Path $Hooks | Out-Null
& git.exe init --quiet $Fixture
& git.exe -C $Fixture config user.name 'CI See fixture'
& git.exe -C $Fixture config user.email 'fixture@example.invalid'
[IO.File]::WriteAllText((Join-Path $Fixture 'input.txt'), "fixture`n")
$Hook = "#!/bin/sh`nprintf 'ran\n' > hook-ran`nexit 73`n"
[IO.File]::WriteAllText((Join-Path $Hooks 'pre-commit'), $Hook, [Text.UTF8Encoding]::new($false))
& git.exe -C $Fixture add input.txt
& git.exe -C $Fixture -c "core.hooksPath=$Hooks" commit -m should-fail
$NormalExit = $LASTEXITCODE
if ($NormalExit -eq 0) { throw 'FAIL: ordinary commit bypassed failing hook' }
if (& git.exe -C $Fixture rev-parse --verify HEAD 2>$null) { throw 'FAIL: ordinary commit created HEAD' }
if (-not (Test-Path (Join-Path $Fixture 'hook-ran'))) { throw 'FAIL: hook marker missing' }
Remove-Item -LiteralPath (Join-Path $Fixture 'hook-ran')
& git.exe -C $Fixture -c "core.hooksPath=$Hooks" commit --no-verify -m fixture
$BypassExit = $LASTEXITCODE
if ($BypassExit -ne 0) { throw 'FAIL: expected no-verify commit to succeed' }
if (Test-Path (Join-Path $Fixture 'hook-ran')) { throw 'FAIL: hook ran under no-verify' }
& git.exe -C $Fixture cat-file -p HEAD
if ((& git.exe -C $Fixture rev-list --count HEAD) -ne '1') { throw 'FAIL: expected one fixture commit' }
# Save private transcript and hashes before ownership-checked cleanup.
$Expected = Join-Path $env:TEMP "ci-see-hook-$RunId"
if ([IO.Path]::GetFullPath($Fixture) -ne [IO.Path]::GetFullPath($Expected)) { throw 'BLOCKED: ownership mismatch' }
Remove-Item -LiteralPath $Fixture -Recurse -Force
```

### WSL2 Git

```sh
set -eu
ROOT=$(mktemp -d "${TMPDIR:-/tmp}/ci-see-hook.XXXXXX")
REPO="$ROOT/repo"; HOOKS="$ROOT/hooks"
mkdir -p "$REPO" "$HOOKS"
cd "$REPO"
git init -q
git config user.name 'CI See fixture'
git config user.email 'fixture@example.invalid'
printf 'fixture\n' > input.txt
printf '%s\n' '#!/bin/sh' 'printf "ran\n" > hook-ran' 'exit 73' > "$HOOKS/pre-commit"
chmod 700 "$HOOKS/pre-commit"
git add input.txt
set +e
git -c core.hooksPath="$HOOKS" commit -m should-fail
NORMAL_EXIT=$?
set -e
test "$NORMAL_EXIT" -ne 0
test -f "$REPO/hook-ran"
if git rev-parse --verify HEAD >/dev/null 2>&1; then exit 1; fi
rm "$REPO/hook-ran"
git -c core.hooksPath="$HOOKS" commit --no-verify -m fixture
test -f "$REPO/.git/HEAD"
test ! -e "$REPO/hook-ran"
test "$(git rev-list --count HEAD)" -eq 1
git cat-file -p HEAD
sha256sum "$HOOKS/pre-commit"
case "$ROOT" in "${TMPDIR:-/tmp}"/ci-see-hook.*) rm -rf -- "$ROOT" ;; *) exit 1 ;; esac
```

Expected evidence: the ordinary commit exits nonzero, the hook marker exists, and no commit object is created; after removing the marker, `--no-verify` succeeds, creates exactly one commit, and leaves the marker absent. Preserve command transcript, exit statuses, hook bytes/hash, commit object, and marker states privately. If the fixture fails for a path/shell problem or cannot prove those observations, classify it INCOMPLETE and keep `GIT-HOOK-BYPASS-NO-VERIFY` UNVERIFIED. This controlled Phase B fixture is feasibility evidence only; product hook behavior and enforcement limitations remain open for S5.
