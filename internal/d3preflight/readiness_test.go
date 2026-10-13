package d3preflight

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestValidateReadinessAcceptsBoundProfileAsConfigurationOnly(t *testing.T) {
	input := validReadinessInput(t, readinessProfileBytes(t, validReadinessProfile()))
	if result := ValidateReadiness(input); result.Status != StatusConfigValid || result.Code != "readiness_profile_validated_offline" {
		t.Fatalf("result = %+v", result)
	}
}

func TestValidateReadinessHashesOriginalProfileBytes(t *testing.T) {
	profileBytes := readinessProfileBytes(t, validReadinessProfile())
	input := validReadinessInput(t, profileBytes)
	input.ProfileBytes = append(append([]byte(nil), profileBytes...), '\n')
	if result := ValidateReadiness(input); result.Status != StatusBlocked || result.Code != "readiness_profile_digest_mismatch" {
		t.Fatalf("result = %+v", result)
	}
}

func TestValidateReadinessRequiresOwnerSelectionsAndEvidence(t *testing.T) {
	profile := validReadinessProfile()
	profile.Network.BoundaryID = ""
	input := validReadinessInput(t, readinessProfileBytes(t, profile))
	if result := ValidateReadiness(input); result.Status != StatusIncomplete || result.Code != "readiness_selection_or_evidence_missing" {
		t.Fatalf("result = %+v", result)
	}
}

func TestValidateReadinessRejectsMalformedAndAmbiguousJSON(t *testing.T) {
	base := readinessProfileBytes(t, validReadinessProfile())
	tests := []struct {
		name string
		edit func([]byte) []byte
	}{
		{"duplicate key", func(data []byte) []byte {
			return bytes.Replace(data, []byte(`"schema":`), []byte(`"schema":"ci-see.d3.environment-readiness.v1","schema":`), 1)
		}},
		{"nested duplicate key", func(data []byte) []byte {
			return bytes.Replace(data, []byte(`"id":"host-01"`), []byte(`"id":"other","id":"host-01"`), 1)
		}},
		{"unknown field", func(data []byte) []byte {
			return append(append([]byte(nil), data[:len(data)-1]...), []byte(`,"unknown":true}`)...)
		}},
		{"nested unknown field", func(data []byte) []byte {
			return bytes.Replace(data, []byte(`"windows_version":"Windows 11"`), []byte(`"windows_version":"Windows 11","unexpected":true`), 1)
		}},
		{"malformed JSON", func(data []byte) []byte { return append(append([]byte(nil), data...), '{') }},
		{"trailing JSON value", func(data []byte) []byte { return append(append([]byte(nil), data...), []byte(` {}`)...) }},
		{"invalid UTF-8", func(data []byte) []byte { return append(append([]byte(nil), data...), 0xff) }},
		{"UTF-8 BOM", func(data []byte) []byte { return append([]byte{0xef, 0xbb, 0xbf}, data...) }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := validReadinessInput(t, test.edit(base))
			if result := ValidateReadiness(input); result.Status != StatusBlocked || result.Code != "readiness_profile_invalid" {
				t.Fatalf("result = %+v", result)
			}
		})
	}
}

func TestValidateReadinessRejectsScopeAndUnsafeSettings(t *testing.T) {
	tests := []struct {
		name string
		edit func(*ReadinessProfile)
		code string
	}{
		{"host scope mismatch", func(profile *ReadinessProfile) { profile.Host.ID = "different-host" }, "readiness_scope_or_boundary_conflict"},
		{"outside guest enforcement disabled", func(profile *ReadinessProfile) { setReadinessBool(&profile.Network.EnforcedOutsideGuest, false) }, "outside_guest_enforcement_required"},
		{"docker host socket enabled", func(profile *ReadinessProfile) { setReadinessBool(&profile.Docker.HostSocketExposed, true) }, "docker_host_socket_must_be_disabled"},
		{"manifest artifact pin mismatch", func(profile *ReadinessProfile) { profile.Dependencies.Tools[0].SHA256 = strings.Repeat("0", 64) }, "readiness_dependency_binding_mismatch"},
		{"rollback snapshot mismatch", func(profile *ReadinessProfile) { profile.Lifecycle.RollbackSnapshot = "other-snapshot" }, "readiness_scope_or_boundary_conflict"},
		{"same leg evidence", func(profile *ReadinessProfile) { profile.Legs.WSL2EvidenceRef = profile.Legs.WindowsNativeEvidenceRef }, "readiness_leg_evidence_must_be_distinct"},
		{"shared acquisition boundary", func(profile *ReadinessProfile) {
			profile.Dependencies.AcquisitionBoundaryID = profile.Network.BoundaryID
		}, "acquisition_boundary_must_be_separate"},
		{"unbounded timeout", func(profile *ReadinessProfile) { profile.Lifecycle.TimeoutSeconds = 0 }, "readiness_lifecycle_policy_invalid"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			profile := validReadinessProfile()
			test.edit(&profile)
			input := validReadinessInput(t, readinessProfileBytes(t, profile))
			if result := ValidateReadiness(input); result.Status != StatusBlocked || result.Code != test.code {
				t.Fatalf("result = %+v, want BLOCKED/%s", result, test.code)
			}
		})
	}
}

func TestValidateReadinessRequiresValidSignedManifest(t *testing.T) {
	input := validReadinessInput(t, readinessProfileBytes(t, validReadinessProfile()))
	input.Manifest.TrustedKeys = nil
	if result := ValidateReadiness(input); result.Status != StatusBlocked || result.Code != "trusted_key_unavailable" {
		t.Fatalf("result = %+v", result)
	}
}

func TestValidateReadinessBoundsInput(t *testing.T) {
	input := validReadinessInput(t, bytes.Repeat([]byte{'x'}, MaxReadinessProfile+1))
	if result := ValidateReadiness(input); result.Status != StatusBlocked || result.Code != "readiness_profile_too_large" {
		t.Fatalf("result = %+v", result)
	}
}

func validReadinessInput(t *testing.T, profileBytes []byte) ReadinessInput {
	t.Helper()
	now := time.Date(2026, 10, 13, 9, 0, 0, 0, time.UTC)
	manifest := fixtureManifest(now.Add(-time.Hour), now.Add(time.Hour))
	digest := sha256.Sum256(profileBytes)
	manifest.Environment.EnvironmentProfileSHA = hex.EncodeToString(digest[:])
	input, _ := signManifest(t, manifest, now)
	return ReadinessInput{Manifest: input, ProfileBytes: profileBytes}
}

func validReadinessProfile() ReadinessProfile {
	return ReadinessProfile{
		Schema:  ReadinessSchema,
		Host:    ReadinessHost{ID: "host-01", WindowsVersion: "Windows 11", HyperVAvailable: readinessBool(true), HyperVEvidenceRef: "host-hyperv-evidence"},
		VM:      ReadinessVM{ID: "vm-01", Disposable: readinessBool(true), NestedVirtualization: readinessBool(true), WSL2Available: readinessBool(true), SnapshotID: "snapshot-01", PrerequisiteEvidenceRef: "vm-prereq-evidence"},
		Legs:    ReadinessLegs{WindowsNativeEvidenceRef: "windows-native-evidence", WSL2EvidenceRef: "wsl2-evidence"},
		Network: ReadinessNetwork{BoundaryID: "outside-boundary-01", DenyByDefault: readinessBool(true), EnforcedOutsideGuest: readinessBool(true), IndependentCapture: readinessBool(true), EnforcementEvidenceRef: "network-enforcement-evidence", CaptureEvidenceRef: "network-capture-evidence"},
		Docker:  ReadinessDocker{DaemonID: "disposable-daemon-01", Disposable: readinessBool(true), ConfinedToVM: readinessBool(true), CredentialsIsolated: readinessBool(true), HostSocketExposed: readinessBool(false), HostMountsAllowed: readinessBool(false), PrivilegedContainers: readinessBool(false), RunOwnedProcessBoundary: readinessBool(true), BoundaryEvidenceRef: "docker-boundary-evidence"},
		Dependencies: ReadinessDependencies{
			AcquisitionBoundaryID: "acquisition-boundary-01", AcquisitionEvidenceRef: "acquisition-evidence",
			Tools:   []ReadinessToolPin{{Name: "act", Version: "v0.2.81", SHA256: strings.Repeat("3", 64)}},
			Images:  []ReadinessImagePin{{Name: "runner", Reference: "registry.example.test/runner@sha256:" + strings.Repeat("f", 64), Platform: "linux/amd64", Digest: "sha256:" + strings.Repeat("f", 64)}},
			Actions: []ReadinessActionPin{{Name: "actions-checkout", Repository: "https://github.com/actions/checkout", Commit: strings.Repeat("4", 40), TreeSHA: strings.Repeat("5", 40), ContentSHA256: strings.Repeat("6", 64)}},
		},
		Lifecycle: ReadinessLifecycle{TimeoutSeconds: 300, ResourceOwnerID: "run-owner-01", CleanupPolicy: "run-owned-only", RollbackSnapshot: "snapshot-01", EvidenceCustodyID: "private-evidence-01", EvidenceRef: "lifecycle-evidence"},
	}
}

func readinessProfileBytes(t *testing.T, profile ReadinessProfile) []byte {
	t.Helper()
	data, err := json.Marshal(profile)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func readinessBool(value bool) *bool { return &value }

func setReadinessBool(target **bool, value bool) { *target = readinessBool(value) }
