package d3preflight

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

const fixtureRunID = "123e4567-e89b-42d3-a456-426614174000"

func TestValidateReportsOfflineValidationOnly(t *testing.T) {
	input, _ := validInput(t, time.Date(2026, 10, 13, 9, 0, 0, 0, time.UTC))
	result := Validate(input)
	if result.Status != StatusConfigValid || result.Code != "manifest_validated_offline" {
		t.Fatalf("result = %+v", result)
	}
}

func TestValidateFailsClosed(t *testing.T) {
	now := time.Date(2026, 10, 13, 9, 0, 0, 0, time.UTC)
	tests := []struct {
		name string
		edit func(*Input)
		want Status
		code string
	}{
		{"missing approval", func(input *Input) { input.ApprovalBytes = nil }, StatusIncomplete, "approval_input_missing"},
		{"unpinned key", func(input *Input) { input.TrustedKeys = nil }, StatusBlocked, "trusted_key_unavailable"},
		{"untrusted clock", func(input *Input) { input.Clock.Trusted = false }, StatusBlocked, "trusted_utc_clock_unavailable"},
		{"non UTC clock", func(input *Input) {
			input.Clock.Now = time.Date(2026, 10, 13, 9, 0, 0, 0, time.FixedZone("UTC-like", 0))
		}, StatusBlocked, "trusted_utc_clock_unavailable"},
		{"scope mismatch", func(input *Input) { input.Expected.HeadSHA = strings.Repeat("f", 40) }, StatusBlocked, "scope_mismatch"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input, _ := validInput(t, now)
			test.edit(&input)
			result := Validate(input)
			if result.Status != test.want || result.Code != test.code {
				t.Fatalf("result = %+v, want %s/%s", result, test.want, test.code)
			}
		})
	}
}

func TestValidateRejectsSignatureFromDifferentPinnedKey(t *testing.T) {
	now := time.Date(2026, 10, 13, 9, 0, 0, 0, time.UTC)
	input, _ := validInput(t, now)
	wrongPrivateKey := ed25519.NewKeyFromSeed(bytes.Repeat([]byte{0x43}, ed25519.SeedSize))
	input.TrustedKeys = map[string]ed25519.PublicKey{
		"fixture-key": append(ed25519.PublicKey(nil), wrongPrivateKey.Public().(ed25519.PublicKey)...),
	}
	if result := Validate(input); result.Status != StatusBlocked || result.Code != "signature_invalid" {
		t.Fatalf("result = %+v", result)
	}
}

func TestValidateRequiresExactSignedBytes(t *testing.T) {
	now := time.Date(2026, 10, 13, 9, 0, 0, 0, time.UTC)
	input, privateKey := validInput(t, now)
	mutated := append([]byte(" "), input.ManifestBytes...)
	input.ManifestBytes = mutated
	input.ApprovalBytes = approvalFor(t, mutated, "fixture-key", privateKey, false)
	if result := Validate(input); result.Code != "signature_invalid" || result.Status != StatusBlocked {
		t.Fatalf("whitespace change with recomputed digest result = %+v", result)
	}

	input, _ = validInput(t, now)
	input.ManifestBytes = append([]byte(" "), input.ManifestBytes...)
	if result := Validate(input); result.Code != "manifest_digest_mismatch" || result.Status != StatusBlocked {
		t.Fatalf("whitespace change result = %+v", result)
	}
}

func TestValidateRejectsMalformedAndAmbiguousJSON(t *testing.T) {
	now := time.Date(2026, 10, 13, 9, 0, 0, 0, time.UTC)
	tests := []struct {
		name string
		edit func([]byte) []byte
	}{
		{"duplicate key", func(data []byte) []byte { return append([]byte(`{"schema":"wrong",`), data[1:]...) }},
		{"case-mismatched key", func(data []byte) []byte { return bytes.Replace(data, []byte(`"schema"`), []byte(`"Schema"`), 1) }},
		{"unknown field", func(data []byte) []byte {
			return append(append([]byte(nil), data[:len(data)-1]...), []byte(`,"unknown":true}`)...)
		}},
		{"trailing value", func(data []byte) []byte { return append(append([]byte(nil), data...), []byte(` {}`)...) }},
		{"UTF-8 BOM", func(data []byte) []byte { return append([]byte{0xef, 0xbb, 0xbf}, data...) }},
		{"invalid UTF-8", func(data []byte) []byte { return append(append([]byte(nil), data...), 0xff) }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input, _ := validInput(t, now)
			input.ManifestBytes = test.edit(input.ManifestBytes)
			if result := Validate(input); result.Status != StatusBlocked || result.Code != "manifest_invalid" {
				t.Fatalf("result = %+v", result)
			}
		})
	}

	input, _ := validInput(t, now)
	input.ApprovalBytes = []byte(`{"schema":"ci-see.d3.approval.v1","schema":"ci-see.d3.approval.v1"}`)
	if result := Validate(input); result.Status != StatusBlocked || result.Code != "approval_invalid" {
		t.Fatalf("duplicate approval keys result = %+v", result)
	}
}

func TestValidateEnforcesTimeWindow(t *testing.T) {
	now := time.Date(2026, 10, 13, 9, 0, 0, 0, time.UTC)
	tests := []struct {
		name    string
		issued  time.Time
		expires time.Time
		code    string
	}{
		{"expired at boundary", now.Add(-time.Hour), now, "approval_expired_or_not_yet_valid"},
		{"not yet issued", now.Add(time.Second), now.Add(time.Hour), "approval_expired_or_not_yet_valid"},
		{"validity exceeds limit", now.Add(-time.Second), now.Add(MaxValidity), "approval_validity_too_long"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			manifest := fixtureManifest(test.issued, test.expires)
			input, _ := signManifest(t, manifest, now)
			result := Validate(input)
			if result.Status != StatusBlocked || result.Code != test.code {
				t.Fatalf("result = %+v, want code %q", result, test.code)
			}
		})
	}
}

func TestValidateAcceptsExactMaximumWindowWithoutSkew(t *testing.T) {
	now := time.Date(2026, 10, 13, 9, 0, 0, 0, time.UTC)
	manifest := fixtureManifest(now.Add(-4*time.Hour), now.Add(4*time.Hour))
	input, _ := signManifest(t, manifest, now)
	if result := Validate(input); result.Status != StatusConfigValid {
		t.Fatalf("exactly eight-hour window result = %+v", result)
	}
}

func TestValidateRejectsInvalidManifestSemantics(t *testing.T) {
	now := time.Date(2026, 10, 13, 9, 0, 0, 0, time.UTC)
	tests := []struct {
		name string
		edit func(*Manifest)
	}{
		{"wildcard destination", func(manifest *Manifest) {
			manifest.Destinations = []Destination{{Stage: "execution", Host: "*.example.com", Port: 443, Protocol: "tcp", Purpose: "registry"}}
		}},
		{"unbounded timeout", func(manifest *Manifest) { manifest.Commands[0].TimeoutSeconds = 0 }},
		{"credential policy", func(manifest *Manifest) { manifest.CredentialPolicy = "host-default" }},
		{"cleanup label mismatch", func(manifest *Manifest) { manifest.Cleanup.RunLabel = "other-run" }},
		{"unknown execution leg", func(manifest *Manifest) { manifest.Environment.ExecutionLeg = "host" }},
		{"bad action pin", func(manifest *Manifest) { manifest.Artifacts[1].Revision = "branch-name" }},
		{"empty artifact source", func(manifest *Manifest) { manifest.Artifacts[0].Source = "" }},
		{"source userinfo", func(manifest *Manifest) { manifest.Artifacts[0].Source = "https://user:secret@example.test/tool" }},
		{"image digest mismatch", func(manifest *Manifest) { manifest.Artifacts[2].ContentSHA256 = strings.Repeat("0", 64) }},
		{"command references unknown artifact", func(manifest *Manifest) { manifest.Commands[0].ArtifactID = "missing" }},
		{"command executable does not match pin", func(manifest *Manifest) { manifest.Commands[0].Argv[0] = "other-tool" }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			manifest := fixtureManifest(now.Add(-time.Hour), now.Add(time.Hour))
			test.edit(&manifest)
			input, _ := signManifest(t, manifest, now)
			if result := Validate(input); result.Status != StatusBlocked || result.Code != "manifest_invalid" {
				t.Fatalf("result = %+v", result)
			}
		})
	}
}

func TestAssessDeadlineIsDeterministic(t *testing.T) {
	now := time.Date(2026, 10, 13, 9, 0, 0, 0, time.UTC)
	tests := []struct {
		name     string
		deadline time.Time
		want     DeadlineDisposition
	}{
		{"before deadline", now.Add(time.Second), DeadlineWithinBudget},
		{"at deadline", now, DeadlineExpired},
		{"after deadline", now.Add(-time.Second), DeadlineExpired},
		{"deadline beyond policy bound", now.Add(maxTimeout + time.Second), DeadlineUnknown},
		{"missing deadline", time.Time{}, DeadlineUnknown},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := AssessDeadline(now, test.deadline); got != test.want {
				t.Fatalf("disposition = %q, want %q", got, test.want)
			}
		})
	}
}

type fakeCleanupStore struct {
	complete  bool
	resources []ResourceObservation
}

func (store fakeCleanupStore) snapshot() (bool, []ResourceObservation) {
	return store.complete, store.resources
}

func TestAssessCleanupUsesOnlyRunOwnedObservations(t *testing.T) {
	tests := []struct {
		name  string
		store fakeCleanupStore
		want  CleanupDisposition
	}{
		{"complete and absent", fakeCleanupStore{true, []ResourceObservation{{ID: "container-1", OwnerRunID: fixtureRunID}}}, CleanupEvidenceComplete},
		{"owned resource remains", fakeCleanupStore{true, []ResourceObservation{{ID: "container-1", OwnerRunID: fixtureRunID, Present: true}}}, CleanupEvidenceUnknown},
		{"foreign resource", fakeCleanupStore{true, []ResourceObservation{{ID: "shared-container", OwnerRunID: "another-run"}}}, CleanupEvidenceBlocked},
		{"incomplete inventory", fakeCleanupStore{false, nil}, CleanupEvidenceUnknown},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			complete, resources := test.store.snapshot()
			if got := AssessCleanup(fixtureRunID, complete, resources); got != test.want {
				t.Fatalf("cleanup evidence = %q, want %q", got, test.want)
			}
		})
	}
}

func validInput(t *testing.T, now time.Time) (Input, ed25519.PrivateKey) {
	t.Helper()
	return signManifest(t, fixtureManifest(now.Add(-time.Hour), now.Add(time.Hour)), now)
}

func signManifest(t *testing.T, manifest Manifest, now time.Time) (Input, ed25519.PrivateKey) {
	t.Helper()
	manifestBytes, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	privateKey := ed25519.NewKeyFromSeed(bytes.Repeat([]byte{0x42}, ed25519.SeedSize))
	approvalBytes := approvalFor(t, manifestBytes, "fixture-key", privateKey, true)
	return Input{
		ManifestBytes: manifestBytes,
		ApprovalBytes: approvalBytes,
		TrustedKeys: map[string]ed25519.PublicKey{
			"fixture-key": append(ed25519.PublicKey(nil), privateKey.Public().(ed25519.PublicKey)...),
		},
		Expected: Scope{
			Repository: "arvin-ramezani/ci-see", PullRequest: 18, HeadSHA: strings.Repeat("a", 40), TreeSHA: strings.Repeat("b", 40),
			RunID: fixtureRunID, HostID: "host-01", VMID: "vm-01", SnapshotID: "snapshot-01", GuestID: "guest-01", ExecutionLeg: "windows-native",
		},
		Clock: Clock{Now: now.UTC(), Trusted: true},
	}, privateKey
}

func approvalFor(t *testing.T, manifestBytes []byte, keyID string, privateKey ed25519.PrivateKey, sign bool) []byte {
	t.Helper()
	digest := sha256.Sum256(manifestBytes)
	signed := append([]byte(signatureDomain), manifestBytes...)
	signature := make([]byte, ed25519.SignatureSize)
	if sign {
		signature = ed25519.Sign(privateKey, signed)
	}
	approvalBytes, err := json.Marshal(Approval{
		Schema: ApprovalSchema, KeyID: keyID, ManifestSHA256: hex.EncodeToString(digest[:]), Signature: base64.StdEncoding.EncodeToString(signature),
	})
	if err != nil {
		t.Fatal(err)
	}
	return approvalBytes
}

func fixtureManifest(issued, expires time.Time) Manifest {
	imageHash := strings.Repeat("f", 64)
	return Manifest{
		Schema: ManifestSchema,
		Identity: Identity{
			Repository: "arvin-ramezani/ci-see", PullRequest: 18, HeadSHA: strings.Repeat("a", 40), TreeSHA: strings.Repeat("b", 40),
			RunID: fixtureRunID, Nonce: strings.Repeat("c", 64), IssuedAt: issued.UTC().Format(time.RFC3339), ExpiresAt: expires.UTC().Format(time.RFC3339),
		},
		Environment: Environment{
			HostID: "host-01", VMID: "vm-01", SnapshotID: "snapshot-01", GuestID: "guest-01", ExecutionLeg: "windows-native",
			WorkspaceSHA256: strings.Repeat("d", 64), EnvironmentProfileSHA: strings.Repeat("e", 64),
		},
		Workflow: Workflow{Path: ".github/workflows/ci.yml", ContentSHA256: strings.Repeat("1", 64), EventSHA256: strings.Repeat("2", 64)},
		Artifacts: []Artifact{
			{ID: "act-binary", Kind: "tool", Name: "act", Source: "https://github.com/nektos/act", Revision: "v0.2.81", ContentSHA256: strings.Repeat("3", 64)},
			{ID: "checkout-action", Kind: "action", Name: "actions-checkout", Source: "https://github.com/actions/checkout", Revision: strings.Repeat("4", 40), TreeSHA: strings.Repeat("5", 40), ContentSHA256: strings.Repeat("6", 64)},
			{ID: "runner-image", Kind: "image", Name: "runner", Source: "registry.example.test/runner", ContentSHA256: imageHash, OCIDigest: "sha256:" + imageHash, Platform: "linux/amd64"},
		},
		Commands:         []Command{{Stage: "execution", ArtifactID: "act-binary", Argv: []string{"act", "--directory", "snapshot"}, TimeoutSeconds: 300}},
		Destinations:     []Destination{},
		Cleanup:          Cleanup{RunLabel: "ci-see-d3-" + fixtureRunID, RollbackSnapshot: "snapshot-01", EvidenceCustodyID: "evidence-01"},
		CredentialPolicy: "dummy-only",
	}
}

