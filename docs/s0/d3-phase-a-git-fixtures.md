# Issue #17 — Phase B Git identity and hook fixtures

**Status:** Proposed Phase B feasibility procedure only; not executed in Phase A.  
**Authority:** [D3 Phase A plan](d3-phase-a-verification-plan.md), [Git gating specification](../specs/git-gating.md), [fixture catalog](fixtures.md).

**Growth decision: KEEP (306 lines).** The Windows-native and WSL2 sections describe paired
legs of the same Git identity and hook-bypass fixtures. They share the fixture purpose,
safety requirements, evidence, and acceptance criteria; keeping them together supports
direct platform comparison and keeps each procedure beside its shared controls.
Splitting them would fragment that comparison and duplicate shared requirements.

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

Run each leg in its disposable VM with no untrusted process running concurrently from directory creation through cleanup. If cleanup cannot revalidate the path, directory identity, and ownership record, preserve the directory for quarantine; never guess that it is safe to remove.

```sh
#!/bin/sh
printf 'ran\n' > hook-ran
exit 73
```

The marker is relative to the hook's repository working directory, so no native Windows path is embedded in a POSIX shell script.

### Windows-native Git

```powershell
function Get-CanonicalFixtureDirectory([string]$Path) {
    $Full = [IO.Path]::GetFullPath($Path)
    $DriveRoot = [IO.Path]::GetPathRoot($Full)
    if (-not $DriveRoot -or $Full.StartsWith('\\')) { throw 'BLOCKED: fixture paths must be on a local drive' }
    if ([IO.DriveInfo]::new($DriveRoot).DriveType -ne [IO.DriveType]::Fixed) {
        throw 'BLOCKED: fixture paths must use a fixed local volume'
    }
    if ($Full.Length -gt $DriveRoot.Length) { $Full = $Full.TrimEnd([char[]]@('\')) }
    $Current = $DriveRoot
    foreach ($Part in $Full.Substring($DriveRoot.Length).Split([char[]]@('\'))) {
        if (-not $Part) { continue }
        $Current = Join-Path $Current $Part
        $Item = Get-Item -LiteralPath $Current -Force -ErrorAction Stop
        if (($Item.Attributes -band [IO.FileAttributes]::ReparsePoint) -ne 0) {
            throw "BLOCKED: reparse point in fixture path: $Current"
        }
    }
    $Resolved = (Resolve-Path -LiteralPath $Full -ErrorAction Stop).ProviderPath
    if (-not [string]::Equals([IO.Path]::GetFullPath($Resolved), $Full, [StringComparison]::OrdinalIgnoreCase)) {
        throw 'BLOCKED: fixture path did not resolve canonically'
    }
    return $Full
}

$ErrorActionPreference = 'Stop'
$RunGuid = [guid]::Empty
if (-not [guid]::TryParseExact([string]$RunId, 'D', [ref]$RunGuid)) {
    throw 'BLOCKED: RunId must be a GUID in D format'
}
$RunId = $RunGuid.ToString('N')
$TempRoot = Get-CanonicalFixtureDirectory $env:TEMP
$Fixture = Join-Path $TempRoot "ci-see-hook-$RunId"
$FixtureCreated = $false
$FixtureSucceeded = $false
$OwnerToken = [guid]::NewGuid().ToString('N')
$OwnerMarker = Join-Path $Fixture '.ci-see-fixture-owner'

if (-not ('CiSee.NativeMethods' -as [type])) {
    Add-Type -TypeDefinition @'
using System;
using System.ComponentModel;
using System.Runtime.InteropServices;
using System.Runtime.InteropServices.ComTypes;
namespace CiSee {
    public static class NativeMethods {
        [StructLayout(LayoutKind.Sequential)]
        public struct FileInformation {
            public uint Attributes;
            public FILETIME CreationTime, LastAccessTime, LastWriteTime;
            public uint VolumeSerialNumber, SizeHigh, SizeLow, LinkCount, IndexHigh, IndexLow;
        }
        [DllImport("kernel32.dll", CharSet = CharSet.Unicode, SetLastError = true)]
        public static extern bool CreateDirectory(string path, IntPtr securityAttributes);
        [DllImport("kernel32.dll", CharSet = CharSet.Unicode, SetLastError = true)]
        private static extern IntPtr CreateFile(string path, uint access, uint share, IntPtr security,
                                                uint creation, uint flags, IntPtr template);
        [DllImport("kernel32.dll", SetLastError = true)]
        private static extern bool GetFileInformationByHandle(IntPtr handle, out FileInformation info);
        [DllImport("kernel32.dll")]
        private static extern bool CloseHandle(IntPtr handle);
        public static string GetDirectoryIdentity(string path) {
            IntPtr handle = CreateFile(path, 0x80, 0x7, IntPtr.Zero, 3, 0x02000000 | 0x00200000, IntPtr.Zero);
            if (handle == new IntPtr(-1)) throw new Win32Exception(Marshal.GetLastWin32Error());
            try {
                FileInformation info;
                if (!GetFileInformationByHandle(handle, out info)) throw new Win32Exception(Marshal.GetLastWin32Error());
                if ((info.Attributes & 0x10) == 0 || (info.Attributes & 0x400) != 0)
                    throw new InvalidOperationException("path is not a plain directory");
                return String.Format("{0:x8}:{1:x8}{2:x8}", info.VolumeSerialNumber, info.IndexHigh, info.IndexLow);
            } finally { CloseHandle(handle); }
        }
    }
}
'@
}
if (-not [CiSee.NativeMethods]::CreateDirectory($Fixture, [IntPtr]::Zero)) {
    $CreateError = [Runtime.InteropServices.Marshal]::GetLastWin32Error()
    throw "BLOCKED: exclusive fixture directory creation failed (Win32 $CreateError)"
}
$FixtureCreated = $true
try {
$FixtureIdentity = [CiSee.NativeMethods]::GetDirectoryIdentity($Fixture)
$Fixture = Get-CanonicalFixtureDirectory $Fixture
if ([CiSee.NativeMethods]::GetDirectoryIdentity($Fixture) -cne $FixtureIdentity) {
    throw 'BLOCKED: fixture directory identity changed after exclusive creation'
}
if (-not [string]::Equals([IO.Path]::GetDirectoryName($Fixture), $TempRoot, [StringComparison]::OrdinalIgnoreCase)) {
    throw 'BLOCKED: fixture is not a direct child of the canonical temp root'
}
$OwnerStream = [IO.File]::Open($OwnerMarker, [IO.FileMode]::CreateNew, [IO.FileAccess]::Write, [IO.FileShare]::None)
try {
    $OwnerBytes = [Text.Encoding]::ASCII.GetBytes("$RunId`n$OwnerToken`n$FixtureIdentity`n")
    $OwnerStream.Write($OwnerBytes, 0, $OwnerBytes.Length)
    $OwnerStream.Flush($true)
}
finally { $OwnerStream.Dispose() }

$Hooks = Join-Path $Fixture 'hooks'
if (-not [CiSee.NativeMethods]::CreateDirectory($Hooks, [IntPtr]::Zero)) {
    $CreateError = [Runtime.InteropServices.Marshal]::GetLastWin32Error()
    throw "BLOCKED: exclusive hooks directory creation failed (Win32 $CreateError)"
}
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
# Save the private transcript and hashes before marking the fixture successful.
$FixtureSucceeded = $true
}
finally {
    if ($FixtureCreated -and -not $FixtureSucceeded) {
        Write-Warning "BLOCKED: preserving failed fixture for quarantine: $Fixture"
    }
    elseif ($FixtureCreated) {
        try {
            $CheckedFixture = Get-CanonicalFixtureDirectory $Fixture
            if (-not [string]::Equals($CheckedFixture, $Fixture, [StringComparison]::OrdinalIgnoreCase) -or
                -not [string]::Equals([IO.Path]::GetDirectoryName($CheckedFixture), $TempRoot, [StringComparison]::OrdinalIgnoreCase)) {
                throw 'fixture path/containment changed'
            }
            if ([CiSee.NativeMethods]::GetDirectoryIdentity($CheckedFixture) -cne $FixtureIdentity) {
                throw 'fixture directory identity changed'
            }
            $MarkerItem = Get-Item -LiteralPath $OwnerMarker -Force -ErrorAction Stop
            if (($MarkerItem.Attributes -band [IO.FileAttributes]::ReparsePoint) -ne 0 -or
                [IO.File]::ReadAllText($OwnerMarker) -cne "$RunId`n$OwnerToken`n$FixtureIdentity`n") {
                throw 'fixture ownership marker is missing, changed, or redirected'
            }
            $ReparseItems = @(Get-ChildItem -LiteralPath $Fixture -Force -Recurse -ErrorAction Stop |
                Where-Object { ($_.Attributes -band [IO.FileAttributes]::ReparsePoint) -ne 0 })
            if ($ReparseItems.Count -ne 0) { throw 'fixture contains a reparse point' }
            if ([CiSee.NativeMethods]::GetDirectoryIdentity($CheckedFixture) -cne $FixtureIdentity) {
                throw 'fixture directory identity changed before deletion'
            }
            Remove-Item -LiteralPath $CheckedFixture -Recurse -Force -ErrorAction Stop
        }
        catch {
            throw "BLOCKED: fixture ownership/cleanup could not be proven; preserve for quarantine at '$Fixture'. $($_.Exception.Message)"
        }
    }
}
```

### WSL2 Git

```sh
set -eu
ROOT_CREATED=0
ROOT=
TMP_ROOT=
ROOT_ID=
OWNER_TOKEN=
OWNER_MARKER=
cleanup_fixture() {
  status=$?
  trap - EXIT HUP INT TERM
  if [ "$ROOT_CREATED" -eq 1 ]; then
    if [ "$status" -ne 0 ]; then
      printf 'BLOCKED: preserving failed fixture for quarantine: %s\n' "$ROOT" >&2
      exit "$status"
    fi
    if [ -z "$ROOT" ] || [ ! -d "$ROOT" ] || [ -L "$ROOT" ] ||
       [ "$(realpath -e -- "$ROOT" 2>/dev/null || true)" != "$ROOT" ] ||
       [ "$(dirname -- "$ROOT")" != "$TMP_ROOT" ] ||
       [ "$(stat -c '%d:%i' -- "$ROOT" 2>/dev/null || true)" != "$ROOT_ID" ] ||
       [ ! -f "$OWNER_MARKER" ] || [ -L "$OWNER_MARKER" ] ||
       [ "$(cat -- "$OWNER_MARKER" 2>/dev/null || true)" != "$(printf '%s\n%s' "$OWNER_TOKEN" "$ROOT_ID")" ]; then
      printf 'BLOCKED: fixture ownership could not be proven; preserve for quarantine: %s\n' "$ROOT" >&2
      exit 1
    fi
    reparse=$(find "$ROOT" -xdev -type l -print -quit) || {
      printf 'BLOCKED: fixture tree could not be inspected; preserve for quarantine: %s\n' "$ROOT" >&2
      exit 1
    }
    if [ -n "$reparse" ]; then
      printf 'BLOCKED: fixture contains a symbolic link; preserve for quarantine: %s\n' "$ROOT" >&2
      exit 1
    fi
    rm -rf --one-file-system -- "$ROOT" || {
      printf 'BLOCKED: verified fixture cleanup failed; inspect/quarantine: %s\n' "$ROOT" >&2
      exit 1
    }
  fi
  exit "$status"
}
trap cleanup_fixture EXIT
trap 'exit 129' HUP
trap 'exit 130' INT
trap 'exit 143' TERM

TMP_ROOT=$(cd -- "${TMPDIR:-/tmp}" && pwd -P)
ROOT=$(mktemp -d "$TMP_ROOT/ci-see-hook.XXXXXXXX")
ROOT_CREATED=1
ROOT_ID=$(stat -c '%d:%i' -- "$ROOT")
ROOT=$(cd -- "$ROOT" && pwd -P)
if [ "$(dirname -- "$ROOT")" != "$TMP_ROOT" ] ||
   [ "$(realpath -e -- "$ROOT")" != "$ROOT" ] ||
   [ "$(stat -c '%d:%i' -- "$ROOT")" != "$ROOT_ID" ]; then
  printf 'BLOCKED: run-owned directory changed during setup; preserve for quarantine: %s\n' "$ROOT" >&2
  exit 1
fi
RUN_ID=${ROOT##*/}
case "$RUN_ID" in ci-see-hook.????????) ;; *) exit 1 ;; esac
RUN_ID_SUFFIX=${RUN_ID#ci-see-hook.}
case "$RUN_ID_SUFFIX" in ''|*[!A-Za-z0-9]*) exit 1 ;; esac
ROOT_ID=$(stat -c '%d:%i' -- "$ROOT")
OWNER_TOKEN=$(cat /proc/sys/kernel/random/uuid)
OWNER_MARKER="$ROOT/.ci-see-fixture-owner"
(umask 077; set -C; printf '%s\n%s\n' "$OWNER_TOKEN" "$ROOT_ID" > "$OWNER_MARKER")
chmod 600 "$OWNER_MARKER"
REPO="$ROOT/repo"; HOOKS="$ROOT/hooks"
mkdir "$REPO" "$HOOKS"
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
# Preserve private transcript and hashes before successful exit; EXIT cleanup
# revalidates canonical containment, directory identity, owner marker, and links.
```

Expected evidence: the ordinary commit exits nonzero, the hook marker exists, and no commit object is created; after removing the marker, `--no-verify` succeeds, creates exactly one commit, and leaves the marker absent. Preserve command transcript, exit statuses, hook bytes/hash, commit object, and marker states privately. If the fixture fails for a path/shell problem or cannot prove those observations, classify it INCOMPLETE and keep `GIT-HOOK-BYPASS-NO-VERIFY` UNVERIFIED. This controlled Phase B fixture is feasibility evidence only; product hook behavior and enforcement limitations remain open for S5.
