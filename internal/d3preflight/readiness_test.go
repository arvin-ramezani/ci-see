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
	for _, edit := range []struct {
		name string
		edit func(*ReadinessProfile)
	}{
		{"network boundary", func(profile *ReadinessProfile) { profile.Network.BoundaryID = "" }},
		{"execution leg", func(profile *ReadinessProfile) { profile.ExecutionLeg = "" }},
		{"Windows guest ID", func(profile *ReadinessProfile) { profile.Legs.WindowsNativeGuestID = "" }},
		{"WSL2 guest ID", func(profile *ReadinessProfile) { profile.Legs.WSL2GuestID = "" }},
	} {
		t.Run(edit.name, func(t *testing.T) {
			profile := validReadinessProfile()
			edit.edit(&profile)
			input := validReadinessInput(t, readinessProfileBytes(t, profile))
			if result := ValidateReadiness(input); result.Status != StatusIncomplete || result.Code != "readiness_selection_or_evidence_missing" {
				t.Fatalf("result = %+v", result)
			}
		})
	}
}

func TestValidateReadinessBindsLifecycleToSignedManifest(t *testing.T) {
	tests := []struct {
		name string
		edit func(*testing.T, *ReadinessProfile, *Manifest)
		code string
	}{
		{"resource owner differs from signed run label", func(_ *testing.T, profile *ReadinessProfile, _ *Manifest) {
			profile.Lifecycle.ResourceOwnerID = "another-run"
		}, "readiness_resource_owner_mismatch"},
		{"signed run label differs from profile owner", func(t *testing.T, _ *ReadinessProfile, manifest *Manifest) {
			manifest.Identity.RunID = "323e4567-e89b-42d3-a456-426614174000"
			manifest.Cleanup.RunLabel = "ci-see-d3-" + manifest.Identity.RunID
			manifest.Environment.ContainerNetwork = manifest.Cleanup.RunLabel
			if !replaceOptionValue(manifest.Commands[0].Argv, "--network", manifest.Cleanup.RunLabel) {
				t.Fatal("network option not found")
			}
			oldToken := "GITHUB_TOKEN=CISEE-DUMMY-" + fixtureRunID
			for index, arg := range manifest.Commands[0].Argv {
				if arg == oldToken {
					manifest.Commands[0].Argv[index] = "GITHUB_TOKEN=CISEE-DUMMY-" + manifest.Identity.RunID
					return
				}
			}
			t.Fatal("dummy token not found")
		}, "readiness_resource_owner_mismatch"},
		{"evidence custody differs from signed cleanup", func(_ *testing.T, profile *ReadinessProfile, _ *Manifest) {
			profile.Lifecycle.EvidenceCustodyID = "another-custody"
		}, "readiness_evidence_custody_mismatch"},
		{"signed evidence custody differs from profile", func(_ *testing.T, _ *ReadinessProfile, manifest *Manifest) {
			manifest.Cleanup.EvidenceCustodyID = "other-evidence-01"
		}, "readiness_evidence_custody_mismatch"},
		{"profile total differs from signed command budget", func(_ *testing.T, profile *ReadinessProfile, _ *Manifest) {
			profile.Lifecycle.TimeoutSeconds++
		}, "readiness_lifecycle_timeout_mismatch"},
		{"signed command budget differs from profile total", func(_ *testing.T, _ *ReadinessProfile, manifest *Manifest) {
			manifest.Commands[0].TimeoutSeconds++
		}, "readiness_lifecycle_timeout_mismatch"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			profile := validReadinessProfile()
			manifest := fixtureManifest(time.Date(2026, 10, 13, 8, 0, 0, 0, time.UTC), time.Date(2026, 10, 13, 10, 0, 0, 0, time.UTC))
			test.edit(t, &profile, &manifest)
			profileBytes := readinessProfileBytes(t, profile)
			digest := sha256.Sum256(profileBytes)
			manifest.Environment.EnvironmentProfileSHA = hex.EncodeToString(digest[:])
			input, _ := signManifest(t, manifest, time.Date(2026, 10, 13, 9, 0, 0, 0, time.UTC))
			input.Expected.RunID = manifest.Identity.RunID
			input.Expected.GuestID = manifest.Environment.GuestID
			input.Expected.ExecutionLeg = manifest.Environment.ExecutionLeg
			result := ValidateReadiness(ReadinessInput{Manifest: input, ProfileBytes: profileBytes})
			if result.Status != StatusBlocked || result.Code != test.code {
				t.Fatalf("result = %+v, want BLOCKED/%s", result, test.code)
			}
		})
	}
}

func TestValidateReadinessAcceptsSumOfSignedSerialCommandTimeouts(t *testing.T) {
	profile := validReadinessProfile()
	profile.Lifecycle.TimeoutSeconds = 360
	profile.Dependencies.Tools = append(profile.Dependencies.Tools, ReadinessToolPin{
		Name: "fetcher", Version: "v1", SHA256: strings.Repeat("7", 64),
	})
	manifest := fixtureManifest(time.Date(2026, 10, 13, 8, 0, 0, 0, time.UTC), time.Date(2026, 10, 13, 10, 0, 0, 0, time.UTC))
	manifest.Artifacts = append(manifest.Artifacts, Artifact{
		ID: "fetch-tool", Kind: "tool", Name: "fetcher", Source: "https://example.test/fetcher", Revision: "v1", ContentSHA256: strings.Repeat("7", 64),
	})
	manifest.Commands = append(manifest.Commands, Command{Stage: "acquisition", ArtifactID: "fetch-tool", Argv: []string{"fetcher"}, TimeoutSeconds: 60})
	profileBytes := readinessProfileBytes(t, profile)
	digest := sha256.Sum256(profileBytes)
	manifest.Environment.EnvironmentProfileSHA = hex.EncodeToString(digest[:])
	input, _ := signManifest(t, manifest, time.Date(2026, 10, 13, 9, 0, 0, 0, time.UTC))
	if result := ValidateReadiness(ReadinessInput{Manifest: input, ProfileBytes: profileBytes}); result.Status != StatusConfigValid {
		t.Fatalf("summed timeout budget result = %+v", result)
	}
}

func TestValidateReadinessBindsSelectedGuestAndLeg(t *testing.T) {
	tests := []struct {
		name string
		edit func(*ReadinessProfile, *Manifest)
		code string
	}{
		{"execution leg differs from signed environment", func(profile *ReadinessProfile, _ *Manifest) {
			profile.ExecutionLeg = "wsl2"
		}, "readiness_execution_leg_mismatch"},
		{"selected Windows guest differs from signed environment", func(profile *ReadinessProfile, _ *Manifest) {
			profile.Legs.WindowsNativeGuestID = "other-guest"
		}, "readiness_guest_scope_mismatch"},
		{"selected WSL2 guest differs from signed environment", func(profile *ReadinessProfile, manifest *Manifest) {
			manifest.Environment.ExecutionLeg = "wsl2"
			manifest.Environment.GuestID = profile.Legs.WSL2GuestID
			profile.Legs.WSL2GuestID = "other-guest"
			profile.ExecutionLeg = "wsl2"
		}, "readiness_guest_scope_mismatch"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			profile := validReadinessProfile()
			manifest := fixtureManifest(time.Date(2026, 10, 13, 8, 0, 0, 0, time.UTC), time.Date(2026, 10, 13, 10, 0, 0, 0, time.UTC))
			test.edit(&profile, &manifest)
			profileBytes := readinessProfileBytes(t, profile)
			digest := sha256.Sum256(profileBytes)
			manifest.Environment.EnvironmentProfileSHA = hex.EncodeToString(digest[:])
			input, _ := signManifest(t, manifest, time.Date(2026, 10, 13, 9, 0, 0, 0, time.UTC))
			input.Expected.GuestID = manifest.Environment.GuestID
			input.Expected.ExecutionLeg = manifest.Environment.ExecutionLeg
			result := ValidateReadiness(ReadinessInput{Manifest: input, ProfileBytes: profileBytes})
			if result.Status != StatusBlocked || result.Code != test.code {
				t.Fatalf("result = %+v, want BLOCKED/%s", result, test.code)
			}
		})
	}

	t.Run("valid WSL2 selected leg", func(t *testing.T) {
		profile := validReadinessProfile()
		profile.ExecutionLeg = "wsl2"
		manifest := fixtureManifest(time.Date(2026, 10, 13, 8, 0, 0, 0, time.UTC), time.Date(2026, 10, 13, 10, 0, 0, 0, time.UTC))
		manifest.Environment.ExecutionLeg = "wsl2"
		manifest.Environment.GuestID = profile.Legs.WSL2GuestID
		profileBytes := readinessProfileBytes(t, profile)
		digest := sha256.Sum256(profileBytes)
		manifest.Environment.EnvironmentProfileSHA = hex.EncodeToString(digest[:])
		input, _ := signManifest(t, manifest, time.Date(2026, 10, 13, 9, 0, 0, 0, time.UTC))
		input.Expected.GuestID = manifest.Environment.GuestID
		input.Expected.ExecutionLeg = manifest.Environment.ExecutionLeg
		if result := ValidateReadiness(ReadinessInput{Manifest: input, ProfileBytes: profileBytes}); result.Status != StatusConfigValid {
			t.Fatalf("WSL2 result = %+v", result)
		}
	})
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
		Schema:       ReadinessSchema,
		Host:         ReadinessHost{ID: "host-01", WindowsVersion: "Windows 11", HyperVAvailable: readinessBool(true), HyperVEvidenceRef: "host-hyperv-evidence"},
		VM:           ReadinessVM{ID: "vm-01", Disposable: readinessBool(true), NestedVirtualization: readinessBool(true), WSL2Available: readinessBool(true), SnapshotID: "snapshot-01", PrerequisiteEvidenceRef: "vm-prereq-evidence"},
		ExecutionLeg: "windows-native",
		Legs: ReadinessLegs{
			WindowsNativeGuestID: "guest-01", WindowsNativeEvidenceRef: "windows-native-evidence",
			WSL2GuestID: "wsl-guest-01", WSL2EvidenceRef: "wsl2-evidence",
		},
		Network: ReadinessNetwork{BoundaryID: "outside-boundary-01", DenyByDefault: readinessBool(true), EnforcedOutsideGuest: readinessBool(true), IndependentCapture: readinessBool(true), EnforcementEvidenceRef: "network-enforcement-evidence", CaptureEvidenceRef: "network-capture-evidence"},
		Docker:  ReadinessDocker{DaemonID: "disposable-daemon-01", Disposable: readinessBool(true), ConfinedToVM: readinessBool(true), CredentialsIsolated: readinessBool(true), HostSocketExposed: readinessBool(false), HostMountsAllowed: readinessBool(false), PrivilegedContainers: readinessBool(false), RunOwnedProcessBoundary: readinessBool(true), BoundaryEvidenceRef: "docker-boundary-evidence"},
		Dependencies: ReadinessDependencies{
			AcquisitionBoundaryID: "acquisition-boundary-01", AcquisitionEvidenceRef: "acquisition-evidence",
			Tools:   []ReadinessToolPin{{Name: "act", Version: "v0.2.81", SHA256: strings.Repeat("3", 64)}},
			Images:  []ReadinessImagePin{{Name: "runner", Reference: "registry.example.test/runner@sha256:" + strings.Repeat("f", 64), Platform: "linux/amd64", Digest: "sha256:" + strings.Repeat("f", 64)}},
			Actions: []ReadinessActionPin{{Name: "actions-checkout", Repository: "https://github.com/actions/checkout", Commit: strings.Repeat("4", 40), TreeSHA: strings.Repeat("5", 40), ContentSHA256: strings.Repeat("6", 64)}},
		},
		Lifecycle: ReadinessLifecycle{TimeoutSeconds: 300, ResourceOwnerID: "ci-see-d3-123e4567-e89b-42d3-a456-426614174000", CleanupPolicy: "run-owned-only", RollbackSnapshot: "snapshot-01", EvidenceCustodyID: "evidence-01", EvidenceRef: "lifecycle-evidence"},
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
