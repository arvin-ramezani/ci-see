package main

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/arvin-ramezani/ci-see/internal/d3preflight"
)

func TestRunRequiresAllInputsAndScope(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run(nil, &stdout, &stderr, nil, d3preflight.Clock{}); code != exitIncomplete {
		t.Fatalf("exit = %d, want %d", code, exitIncomplete)
	}
	if !strings.Contains(stderr.String(), "INCOMPLETE") || strings.Contains(stderr.String(), "READY") || strings.Contains(stderr.String(), "PASS") {
		t.Fatalf("unexpected result text: %q", stderr.String())
	}
}

func TestRunWithNoOperationalTrustKeyBlocks(t *testing.T) {
	args, _, _ := writeFixtureFiles(t)
	var stdout, stderr bytes.Buffer
	if code := runDefault(args, &stdout, &stderr); code != exitBlocked {
		t.Fatalf("exit = %d, want %d", code, exitBlocked)
	}
	if !strings.Contains(stderr.String(), "BLOCKED code=trusted_key_unavailable") || !strings.Contains(stderr.String(), "execution_authorized=false approval_redeemed=false") {
		t.Fatalf("unexpected blocked result: %q", stderr.String())
	}
	if strings.Contains(stderr.String()+stdout.String(), "READY") || strings.Contains(stderr.String()+stdout.String(), "PASS") {
		t.Fatalf("tool emitted an authorization-like status: stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
}

func TestRunPositiveResultStillDoesNotAuthorizeOrRedeem(t *testing.T) {
	args, key, now := writeFixtureFiles(t)
	var stdout, stderr bytes.Buffer
	keys := map[string]ed25519.PublicKey{"fixture-key": key}
	if code := run(args, &stdout, &stderr, keys, d3preflight.Clock{Now: now, Trusted: true}); code != exitOK {
		t.Fatalf("exit = %d, stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	got := stdout.String()
	for _, expected := range []string{"SIGNATURE_VALID", "CONFIG_VALID", "execution_authorized=false", "approval_redeemed=false"} {
		if !strings.Contains(got, expected) {
			t.Errorf("stdout %q missing %q", got, expected)
		}
	}
	if strings.Contains(got, "READY") || strings.Contains(got, "PASS") || stderr.Len() != 0 {
		t.Fatalf("unexpected positive result: stdout=%q stderr=%q", got, stderr.String())
	}
}

func TestRunReadinessProfileValidatesWithoutAuthorizing(t *testing.T) {
	args, _, now := writeFixtureFiles(t)
	profileBytes := writeReadinessFixture(t, args)
	profilePath := filepath.Join(t.TempDir(), "readiness.json")
	if err := os.WriteFile(profilePath, profileBytes, 0600); err != nil {
		t.Fatal(err)
	}
	args = append(args, "--readiness-profile", profilePath)
	privateKey := ed25519.NewKeyFromSeed(bytes.Repeat([]byte{0x42}, ed25519.SeedSize))
	keys := map[string]ed25519.PublicKey{"fixture-key": privateKey.Public().(ed25519.PublicKey)}
	var stdout, stderr bytes.Buffer
	if code := run(args, &stdout, &stderr, keys, d3preflight.Clock{Now: now, Trusted: true}); code != exitOK {
		t.Fatalf("exit = %d, stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	got := stdout.String()
	for _, expected := range []string{"READINESS_PROFILE_VALIDATED", "CONFIG_VALID", "execution_authorized=false", "approval_redeemed=false"} {
		if !strings.Contains(got, expected) {
			t.Errorf("stdout %q missing %q", got, expected)
		}
	}
	if strings.Contains(got, "READY") || strings.Contains(got, "PASS") {
		t.Fatalf("readiness result implied execution or acceptance: %q", got)
	}
}

func TestRunRejectsDuplicateReadinessProfileFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run(nil, &stdout, &stderr, nil, d3preflight.Clock{}); code != exitIncomplete {
		t.Fatalf("baseline exit = %d", code)
	}
	if code := run([]string{"--readiness-profile", "one", "--readiness-profile", "two"}, &stdout, &stderr, nil, d3preflight.Clock{}); code != exitUsage {
		t.Fatalf("duplicate profile flag exit = %d", code)
	}
	if !strings.Contains(stderr.String(), "duplicate_argument") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func writeReadinessFixture(t *testing.T, args []string) []byte {
	t.Helper()
	profile := d3preflight.ReadinessProfile{
		Schema: d3preflight.ReadinessSchema,
		Host:   d3preflight.ReadinessHost{ID: "host-01", WindowsVersion: "Windows 11", HyperVAvailable: boolPointer(true), HyperVEvidenceRef: "host-hyperv-evidence"},
		VM:     d3preflight.ReadinessVM{ID: "vm-01", Disposable: boolPointer(true), NestedVirtualization: boolPointer(true), WSL2Available: boolPointer(true), SnapshotID: "snapshot-01", PrerequisiteEvidenceRef: "vm-prereq-evidence"},
		Legs: d3preflight.ReadinessLegs{
			WindowsNativeGuestID: "guest-01", WindowsNativeEvidenceRef: "windows-native-evidence",
			WSL2GuestID: "wsl-guest-01", WSL2EvidenceRef: "wsl2-evidence",
		},
		Network: d3preflight.ReadinessNetwork{BoundaryID: "outside-boundary-01", DenyByDefault: boolPointer(true), EnforcedOutsideGuest: boolPointer(true), IndependentCapture: boolPointer(true), EnforcementEvidenceRef: "network-enforcement-evidence", CaptureEvidenceRef: "network-capture-evidence"},
		Docker:  d3preflight.ReadinessDocker{DaemonID: "disposable-daemon-01", Disposable: boolPointer(true), ConfinedToVM: boolPointer(true), CredentialsIsolated: boolPointer(true), HostSocketExposed: boolPointer(false), HostMountsAllowed: boolPointer(false), PrivilegedContainers: boolPointer(false), RunOwnedProcessBoundary: boolPointer(true), BoundaryEvidenceRef: "docker-boundary-evidence"},
		Dependencies: d3preflight.ReadinessDependencies{
			AcquisitionBoundaryID: "acquisition-boundary-01", AcquisitionEvidenceRef: "acquisition-evidence",
			Tools:   []d3preflight.ReadinessToolPin{{Name: "act", Version: "v0.2.81", SHA256: strings.Repeat("3", 64)}},
			Images:  []d3preflight.ReadinessImagePin{{Name: "runner", Reference: "registry.example.test/runner@sha256:" + strings.Repeat("f", 64), Platform: "linux/amd64", Digest: "sha256:" + strings.Repeat("f", 64)}},
			Actions: []d3preflight.ReadinessActionPin{{Name: "actions-checkout", Repository: "https://github.com/actions/checkout", Commit: strings.Repeat("4", 40), TreeSHA: strings.Repeat("5", 40), ContentSHA256: strings.Repeat("6", 64)}},
		},
		Lifecycle: d3preflight.ReadinessLifecycle{CleanupPolicy: "run-owned-only", RollbackSnapshot: "snapshot-01", EvidenceRef: "lifecycle-evidence"},
	}
	manifestPath := valueAfter(args, "--manifest")
	approvalPath := valueAfter(args, "--approval")
	manifestBytes, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	var manifest d3preflight.Manifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		t.Fatal(err)
	}
	profile.ExecutionLeg = manifest.Environment.ExecutionLeg
	if profile.ExecutionLeg == "wsl2" {
		profile.Legs.WSL2GuestID = manifest.Environment.GuestID
	} else {
		profile.Legs.WindowsNativeGuestID = manifest.Environment.GuestID
	}
	profile.Lifecycle.ResourceOwnerID = manifest.Cleanup.RunLabel
	profile.Lifecycle.EvidenceCustodyID = manifest.Cleanup.EvidenceCustodyID
	for _, command := range manifest.Commands {
		profile.Lifecycle.TimeoutSeconds += command.TimeoutSeconds
	}
	profileBytes, err := json.Marshal(profile)
	if err != nil {
		t.Fatal(err)
	}
	profileDigest := sha256.Sum256(profileBytes)
	manifest.Environment.EnvironmentProfileSHA = hex.EncodeToString(profileDigest[:])
	manifestBytes, err = json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(manifestPath, manifestBytes, 0600); err != nil {
		t.Fatal(err)
	}
	privateKey := ed25519.NewKeyFromSeed(bytes.Repeat([]byte{0x42}, ed25519.SeedSize))
	manifestDigest := sha256.Sum256(manifestBytes)
	signed := append([]byte("CI-SEE-D3-MANIFEST-V1\x00"), manifestBytes...)
	approvalBytes, err := json.Marshal(d3preflight.Approval{
		Schema: d3preflight.ApprovalSchema, KeyID: "fixture-key", ManifestSHA256: hex.EncodeToString(manifestDigest[:]),
		Signature: base64.StdEncoding.EncodeToString(ed25519.Sign(privateKey, signed)),
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(approvalPath, approvalBytes, 0600); err != nil {
		t.Fatal(err)
	}
	return profileBytes
}

func valueAfter(args []string, name string) string {
	for index := 0; index+1 < len(args); index++ {
		if args[index] == name {
			return args[index+1]
		}
	}
	return ""
}

func boolPointer(value bool) *bool { return &value }

func TestRunWithPinnedFixtureKeyStillRequiresTrustedClock(t *testing.T) {
	args, key, _ := writeFixtureFiles(t)
	var stdout, stderr bytes.Buffer
	if code := run(args, &stdout, &stderr, map[string]ed25519.PublicKey{"fixture-key": key}, d3preflight.Clock{}); code != exitBlocked {
		t.Fatalf("exit = %d, want %d", code, exitBlocked)
	}
	if !strings.Contains(stderr.String(), "BLOCKED code=trusted_utc_clock_unavailable") || stdout.Len() != 0 {
		t.Fatalf("stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
}

func TestRunInputFailureIsIncomplete(t *testing.T) {
	args, _, _ := writeFixtureFiles(t)
	args[1] = "missing-manifest.json"
	var stdout, stderr bytes.Buffer
	if code := run(args, &stdout, &stderr, nil, d3preflight.Clock{}); code != exitIncomplete {
		t.Fatalf("exit = %d, want %d", code, exitIncomplete)
	}
	if !strings.Contains(stderr.String(), "INCOMPLETE code=manifest_unavailable") || stdout.Len() != 0 {
		t.Fatalf("stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
}

func TestRunRejectsExtraAndMalformedArguments(t *testing.T) {
	for _, args := range [][]string{{"--unknown"}, {"extra-positional"}} {
		var stdout, stderr bytes.Buffer
		if code := run(args, &stdout, &stderr, nil, d3preflight.Clock{}); code != exitUsage || !strings.Contains(stderr.String(), "INCOMPLETE") {
			t.Errorf("args=%q exit=%d stderr=%q", args, code, stderr.String())
		}
	}
	for _, duplicate := range [][]string{{"--pr", "17"}, {"-pr", "17"}, {"--pr=18", "-pr=19"}} {
		args, _, _ := writeFixtureFiles(t)
		args = append(args, duplicate...)
		var stdout, stderr bytes.Buffer
		if code := run(args, &stdout, &stderr, nil, d3preflight.Clock{}); code != exitUsage || !strings.Contains(stderr.String(), "duplicate_argument") {
			t.Errorf("duplicate args=%q exit=%d stderr=%q", duplicate, code, stderr.String())
		}
	}
}

func TestReadBoundedRejectsOversizeAndMissingFiles(t *testing.T) {
	path := filepath.Join(t.TempDir(), "input.json")
	if err := os.WriteFile(path, []byte("1234"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, ok := readBounded(path, 3); ok {
		t.Fatal("oversize input accepted")
	}
	if _, ok := readBounded(filepath.Join(t.TempDir(), "missing"), 4); ok {
		t.Fatal("missing input accepted")
	}
}

func writeFixtureFiles(t *testing.T) ([]string, ed25519.PublicKey, time.Time) {
	t.Helper()
	now := time.Date(2026, 10, 13, 9, 0, 0, 0, time.UTC)
	runID := "123e4567-e89b-42d3-a456-426614174000"
	hashImage := strings.Repeat("f", 64)
	manifest := d3preflight.Manifest{
		Schema: d3preflight.ManifestSchema,
		Identity: d3preflight.Identity{
			Repository: "arvin-ramezani/ci-see", PullRequest: 18, HeadSHA: strings.Repeat("a", 40), TreeSHA: strings.Repeat("b", 40),
			RunID: runID, Nonce: strings.Repeat("c", 64), IssuedAt: now.Add(-time.Minute).Format(time.RFC3339), ExpiresAt: now.Add(time.Hour).Format(time.RFC3339),
		},
		Environment: d3preflight.Environment{
			HostID: "host-01", VMID: "vm-01", SnapshotID: "snapshot-01", GuestID: "guest-01", ExecutionLeg: "windows-native",
			ContainerNetwork: "ci-see-d3-" + runID,
			WorkspaceSHA256:  strings.Repeat("d", 64), EnvironmentProfileSHA: strings.Repeat("e", 64),
		},
		Workflow: d3preflight.Workflow{Path: ".github/workflows/ci.yml", EventName: "push", ContentSHA256: strings.Repeat("1", 64), EventSHA256: strings.Repeat("2", 64)},
		Artifacts: []d3preflight.Artifact{
			{ID: "act-binary", Kind: "tool", Name: "act", Source: "https://github.com/nektos/act", Revision: "v0.2.81", ContentSHA256: strings.Repeat("3", 64)},
			{ID: "checkout-action", Kind: "action", Name: "actions-checkout", Source: "https://github.com/actions/checkout", Revision: strings.Repeat("4", 40), TreeSHA: strings.Repeat("5", 40), ContentSHA256: strings.Repeat("6", 64)},
			{ID: "runner-image", Kind: "image", Name: "runner", Source: "registry.example.test/runner", ContentSHA256: hashImage, OCIDigest: "sha256:" + hashImage, Platform: "linux/amd64", ActPlatform: "ubuntu-22.04"},
		},
		Commands: []d3preflight.Command{{
			Stage: "execution", ArtifactID: "act-binary", ImageArtifactID: "runner-image", ActionArtifactIDs: []string{"checkout-action"},
			Argv: []string{
				"act", "--directory", "snapshot", "-W", ".github/workflows/ci.yml", "-e", "event.json", "--action-cache-path", "action-cache",
				"-P", "ubuntu-22.04=registry.example.test/runner@sha256:" + hashImage, "--network", "ci-see-d3-" + runID,
				"--container-daemon-socket", "-", "--pull=false", "--rebuild=false", "--action-offline-mode",
				"--env-file", "empty.env", "--secret-file", "empty.secrets", "--var-file", "empty.vars", "--input-file", "empty.input",
				"-s", "GITHUB_TOKEN=CISEE-DUMMY-" + runID, "push",
			},
			TimeoutSeconds: 300,
		}},
		Destinations:     []d3preflight.Destination{},
		Cleanup:          d3preflight.Cleanup{RunLabel: "ci-see-d3-" + runID, RollbackSnapshot: "snapshot-01", EvidenceCustodyID: "evidence-01"},
		CredentialPolicy: "dummy-only",
	}
	manifestBytes, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	privateKey := ed25519.NewKeyFromSeed(bytes.Repeat([]byte{0x42}, ed25519.SeedSize))
	publicKey := append(ed25519.PublicKey(nil), privateKey.Public().(ed25519.PublicKey)...)
	digest := sha256.Sum256(manifestBytes)
	signed := append([]byte("CI-SEE-D3-MANIFEST-V1\x00"), manifestBytes...)
	approvalBytes, err := json.Marshal(d3preflight.Approval{
		Schema: d3preflight.ApprovalSchema, KeyID: "fixture-key", ManifestSHA256: hex.EncodeToString(digest[:]),
		Signature: base64.StdEncoding.EncodeToString(ed25519.Sign(privateKey, signed)),
	})
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	manifestPath := filepath.Join(root, "manifest.json")
	approvalPath := filepath.Join(root, "approval.json")
	if err := os.WriteFile(manifestPath, manifestBytes, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(approvalPath, approvalBytes, 0600); err != nil {
		t.Fatal(err)
	}
	args := []string{
		"--manifest", manifestPath, "--approval", approvalPath, "--repository", "arvin-ramezani/ci-see", "--pr", "18",
		"--head-sha", strings.Repeat("a", 40), "--tree-sha", strings.Repeat("b", 40), "--run-id", runID,
		"--host-id", "host-01", "--vm-id", "vm-01", "--snapshot-id", "snapshot-01", "--guest-id", "guest-01", "--leg", "windows-native",
	}
	return args, publicKey, now
}
