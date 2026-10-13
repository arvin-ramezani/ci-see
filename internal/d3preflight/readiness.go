package d3preflight

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	ReadinessSchema     = "ci-see.d3.environment-readiness.v1"
	MaxReadinessProfile = 64 << 10
	maxReadinessPins    = 128
)

// ReadinessInput validates an environment profile against a separately validated
// signed manifest. A successful result is configuration validation only.
type ReadinessInput struct {
	Manifest     Input
	ProfileBytes []byte
}

type ReadinessProfile struct {
	Schema       string                `json:"schema"`
	Host         ReadinessHost         `json:"host"`
	VM           ReadinessVM           `json:"vm"`
	Legs         ReadinessLegs         `json:"legs"`
	Network      ReadinessNetwork      `json:"network"`
	Docker       ReadinessDocker       `json:"docker"`
	Dependencies ReadinessDependencies `json:"dependencies"`
	Lifecycle    ReadinessLifecycle    `json:"lifecycle"`
}

type ReadinessHost struct {
	ID                string `json:"id"`
	WindowsVersion    string `json:"windows_version"`
	HyperVAvailable   *bool  `json:"hyper_v_available"`
	HyperVEvidenceRef string `json:"hyper_v_evidence_ref"`
}

type ReadinessVM struct {
	ID                      string `json:"id"`
	Disposable              *bool  `json:"disposable"`
	NestedVirtualization    *bool  `json:"nested_virtualization"`
	WSL2Available           *bool  `json:"wsl2_available"`
	SnapshotID              string `json:"snapshot_id"`
	PrerequisiteEvidenceRef string `json:"prerequisite_evidence_ref"`
}

type ReadinessLegs struct {
	WindowsNativeEvidenceRef string `json:"windows_native_evidence_ref"`
	WSL2EvidenceRef          string `json:"wsl2_evidence_ref"`
}

type ReadinessNetwork struct {
	BoundaryID             string `json:"boundary_id"`
	DenyByDefault          *bool  `json:"deny_by_default"`
	EnforcedOutsideGuest   *bool  `json:"enforced_outside_guest"`
	IndependentCapture     *bool  `json:"independent_capture"`
	EnforcementEvidenceRef string `json:"enforcement_evidence_ref"`
	CaptureEvidenceRef     string `json:"capture_evidence_ref"`
}

type ReadinessDocker struct {
	DaemonID                string `json:"daemon_id"`
	Disposable              *bool  `json:"disposable"`
	ConfinedToVM            *bool  `json:"confined_to_vm"`
	CredentialsIsolated     *bool  `json:"credentials_isolated"`
	HostSocketExposed       *bool  `json:"host_socket_exposed"`
	HostMountsAllowed       *bool  `json:"host_mounts_allowed"`
	PrivilegedContainers    *bool  `json:"privileged_containers"`
	RunOwnedProcessBoundary *bool  `json:"run_owned_process_boundary"`
	BoundaryEvidenceRef     string `json:"boundary_evidence_ref"`
}

type ReadinessDependencies struct {
	AcquisitionBoundaryID  string               `json:"acquisition_boundary_id"`
	AcquisitionEvidenceRef string               `json:"acquisition_evidence_ref"`
	Tools                  []ReadinessToolPin   `json:"tools"`
	Images                 []ReadinessImagePin  `json:"images"`
	Actions                []ReadinessActionPin `json:"actions"`
}

type ReadinessToolPin struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	SHA256  string `json:"sha256"`
}

type ReadinessImagePin struct {
	Name      string `json:"name"`
	Reference string `json:"reference"`
	Platform  string `json:"platform"`
	Digest    string `json:"digest"`
}

type ReadinessActionPin struct {
	Name          string `json:"name"`
	Repository    string `json:"repository"`
	Commit        string `json:"commit"`
	TreeSHA       string `json:"tree_sha"`
	ContentSHA256 string `json:"content_sha256"`
}

type ReadinessLifecycle struct {
	TimeoutSeconds    int64  `json:"timeout_seconds"`
	ResourceOwnerID   string `json:"resource_owner_id"`
	CleanupPolicy     string `json:"cleanup_policy"`
	RollbackSnapshot  string `json:"rollback_snapshot_id"`
	EvidenceCustodyID string `json:"evidence_custody_id"`
	EvidenceRef       string `json:"evidence_ref"`
}

// ValidateReadiness requires the existing manifest signature, scope, and clock
// checks before accepting a profile. It does not authorize or launch execution.
func ValidateReadiness(input ReadinessInput) Result {
	if len(input.ProfileBytes) == 0 {
		return result(StatusIncomplete, "readiness_profile_missing")
	}
	if len(input.ProfileBytes) > MaxReadinessProfile {
		return result(StatusBlocked, "readiness_profile_too_large")
	}
	manifestResult := Validate(input.Manifest)
	if manifestResult.Status != StatusConfigValid {
		return manifestResult
	}
	var manifest Manifest
	if !decodeStrict(input.Manifest.ManifestBytes, &manifest, "manifest") {
		return result(StatusBlocked, "manifest_invalid")
	}
	var profile ReadinessProfile
	if !decodeReadinessStrict(input.ProfileBytes, &profile) {
		return result(StatusBlocked, "readiness_profile_invalid")
	}
	digest := sha256.Sum256(input.ProfileBytes)
	if hex.EncodeToString(digest[:]) != manifest.Environment.EnvironmentProfileSHA {
		return result(StatusBlocked, "readiness_profile_digest_mismatch")
	}
	if status, code := validateReadinessProfile(profile, manifest); status != StatusConfigValid {
		return result(status, code)
	}
	return result(StatusConfigValid, "readiness_profile_validated_offline")
}

func validateReadinessProfile(profile ReadinessProfile, manifest Manifest) (Status, string) {
	if profile.Schema != ReadinessSchema {
		return StatusBlocked, "readiness_profile_schema_invalid"
	}
	if missing(profile.Host.ID, profile.Host.WindowsVersion, profile.Host.HyperVEvidenceRef,
		profile.VM.ID, profile.VM.SnapshotID, profile.VM.PrerequisiteEvidenceRef,
		profile.Legs.WindowsNativeEvidenceRef, profile.Legs.WSL2EvidenceRef,
		profile.Network.BoundaryID, profile.Network.EnforcementEvidenceRef, profile.Network.CaptureEvidenceRef,
		profile.Docker.DaemonID, profile.Docker.BoundaryEvidenceRef,
		profile.Dependencies.AcquisitionBoundaryID, profile.Dependencies.AcquisitionEvidenceRef,
		profile.Lifecycle.ResourceOwnerID, profile.Lifecycle.RollbackSnapshot, profile.Lifecycle.EvidenceCustodyID, profile.Lifecycle.EvidenceRef) {
		return StatusIncomplete, "readiness_selection_or_evidence_missing"
	}
	if profile.Host.ID != manifest.Environment.HostID || profile.Host.WindowsVersion != "Windows 11" || profile.VM.ID != manifest.Environment.VMID ||
		profile.VM.SnapshotID != manifest.Environment.SnapshotID || profile.Network.BoundaryID == "host" ||
		profile.Docker.DaemonID == "host" || profile.Lifecycle.RollbackSnapshot != profile.VM.SnapshotID {
		return StatusBlocked, "readiness_scope_or_boundary_conflict"
	}
	if profile.Legs.WindowsNativeEvidenceRef == profile.Legs.WSL2EvidenceRef {
		return StatusBlocked, "readiness_leg_evidence_must_be_distinct"
	}
	if profile.Dependencies.AcquisitionBoundaryID == profile.Network.BoundaryID {
		return StatusBlocked, "acquisition_boundary_must_be_separate"
	}
	if !validReadinessIDs(profile.Host.ID, profile.Host.HyperVEvidenceRef, profile.VM.ID, profile.VM.SnapshotID,
		profile.VM.PrerequisiteEvidenceRef, profile.Legs.WindowsNativeEvidenceRef, profile.Legs.WSL2EvidenceRef,
		profile.Network.BoundaryID, profile.Network.EnforcementEvidenceRef, profile.Network.CaptureEvidenceRef,
		profile.Docker.DaemonID, profile.Docker.BoundaryEvidenceRef, profile.Dependencies.AcquisitionBoundaryID,
		profile.Dependencies.AcquisitionEvidenceRef, profile.Lifecycle.ResourceOwnerID, profile.Lifecycle.RollbackSnapshot,
		profile.Lifecycle.EvidenceCustodyID, profile.Lifecycle.EvidenceRef) {
		return StatusBlocked, "readiness_identifier_invalid"
	}
	if status, code := requireTrue("host_hyper_v_unavailable", profile.Host.HyperVAvailable); status != StatusConfigValid {
		return status, code
	}
	if status, code := requireTrue("disposable_vm_required", profile.VM.Disposable); status != StatusConfigValid {
		return status, code
	}
	if status, code := requireTrue("nested_virtualization_required", profile.VM.NestedVirtualization); status != StatusConfigValid {
		return status, code
	}
	if status, code := requireTrue("wsl2_required", profile.VM.WSL2Available); status != StatusConfigValid {
		return status, code
	}
	if status, code := requireTrue("network_deny_by_default_required", profile.Network.DenyByDefault); status != StatusConfigValid {
		return status, code
	}
	if status, code := requireTrue("outside_guest_enforcement_required", profile.Network.EnforcedOutsideGuest); status != StatusConfigValid {
		return status, code
	}
	if status, code := requireTrue("independent_network_capture_required", profile.Network.IndependentCapture); status != StatusConfigValid {
		return status, code
	}
	for _, check := range []struct {
		code string
		flag *bool
	}{
		{"disposable_docker_daemon_required", profile.Docker.Disposable},
		{"docker_must_be_confined_to_vm", profile.Docker.ConfinedToVM},
		{"docker_credentials_must_be_isolated", profile.Docker.CredentialsIsolated},
		{"docker_host_socket_must_be_disabled", inverse(profile.Docker.HostSocketExposed)},
		{"docker_host_mounts_must_be_disabled", inverse(profile.Docker.HostMountsAllowed)},
		{"privileged_containers_must_be_disabled", inverse(profile.Docker.PrivilegedContainers)},
		{"run_owned_process_boundary_required", profile.Docker.RunOwnedProcessBoundary},
	} {
		if status, code := requireTrue(check.code, check.flag); status != StatusConfigValid {
			return status, code
		}
	}
	if profile.Lifecycle.TimeoutSeconds < 1 || profile.Lifecycle.TimeoutSeconds > int64((8*time.Hour)/time.Second) ||
		profile.Lifecycle.CleanupPolicy != "run-owned-only" {
		return StatusBlocked, "readiness_lifecycle_policy_invalid"
	}
	if len(profile.Dependencies.Tools) == 0 || len(profile.Dependencies.Images) == 0 || len(profile.Dependencies.Actions) == 0 {
		return StatusIncomplete, "readiness_dependency_pins_missing"
	}
	if len(profile.Dependencies.Tools)+len(profile.Dependencies.Images)+len(profile.Dependencies.Actions) > maxReadinessPins {
		return StatusBlocked, "readiness_dependency_pins_too_many"
	}
	if !readinessPinsMatchManifest(profile.Dependencies, manifest.Artifacts) {
		return StatusBlocked, "readiness_dependency_binding_mismatch"
	}
	for _, pin := range profile.Dependencies.Tools {
		if missing(pin.Name, pin.Version, pin.SHA256) {
			return StatusIncomplete, "readiness_dependency_pins_missing"
		}
	}
	for _, pin := range profile.Dependencies.Images {
		if missing(pin.Name, pin.Reference, pin.Platform, pin.Digest) {
			return StatusIncomplete, "readiness_dependency_pins_missing"
		}
	}
	for _, pin := range profile.Dependencies.Actions {
		if missing(pin.Name, pin.Repository, pin.Commit, pin.TreeSHA, pin.ContentSHA256) {
			return StatusIncomplete, "readiness_dependency_pins_missing"
		}
	}
	return StatusConfigValid, "readiness_profile_validated_offline"
}

func validReadinessIDs(values ...string) bool {
	for _, value := range values {
		if !validID(value) {
			return false
		}
	}
	return true
}

func readinessPinsMatchManifest(pins ReadinessDependencies, artifacts []Artifact) bool {
	tools := make(map[string]ReadinessToolPin, len(pins.Tools))
	for _, pin := range pins.Tools {
		if _, duplicate := tools[pin.Name]; duplicate {
			return false
		}
		tools[pin.Name] = pin
	}
	images := make(map[string]ReadinessImagePin, len(pins.Images))
	for _, pin := range pins.Images {
		if _, duplicate := images[pin.Name]; duplicate {
			return false
		}
		images[pin.Name] = pin
	}
	actions := make(map[string]ReadinessActionPin, len(pins.Actions))
	for _, pin := range pins.Actions {
		if _, duplicate := actions[pin.Name]; duplicate {
			return false
		}
		actions[pin.Name] = pin
	}
	for _, artifact := range artifacts {
		switch artifact.Kind {
		case "tool":
			want := ReadinessToolPin{Name: artifact.Name, Version: artifact.Revision, SHA256: artifact.ContentSHA256}
			if tools[artifact.Name] != want {
				return false
			}
			delete(tools, artifact.Name)
		case "image":
			want := ReadinessImagePin{Name: artifact.Name, Reference: artifact.Source + "@" + artifact.OCIDigest, Platform: artifact.Platform, Digest: artifact.OCIDigest}
			if images[artifact.Name] != want {
				return false
			}
			delete(images, artifact.Name)
		case "action":
			want := ReadinessActionPin{Name: artifact.Name, Repository: artifact.Source, Commit: artifact.Revision, TreeSHA: artifact.TreeSHA, ContentSHA256: artifact.ContentSHA256}
			if actions[artifact.Name] != want {
				return false
			}
			delete(actions, artifact.Name)
		default:
			return false
		}
	}
	return len(tools) == 0 && len(images) == 0 && len(actions) == 0
}

func requireTrue(code string, value *bool) (Status, string) {
	if value == nil {
		return StatusIncomplete, "readiness_selection_or_evidence_missing"
	}
	if !*value {
		return StatusBlocked, code
	}
	return StatusConfigValid, ""
}

func inverse(value *bool) *bool {
	if value == nil {
		return nil
	}
	result := !*value
	return &result
}

func missing(values ...string) bool {
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			return true
		}
	}
	return false
}

func decodeReadinessStrict(data []byte, profile *ReadinessProfile) bool {
	if len(data) == 0 || !utf8.Valid(data) || bytes.HasPrefix(data, []byte{0xef, 0xbb, 0xbf}) {
		return false
	}
	scan := json.NewDecoder(bytes.NewReader(data))
	token, err := scan.Token()
	delim, ok := token.(json.Delim)
	if err != nil || !ok || delim != '{' || scanReadinessObject(scan, 1, "profile") != nil {
		return false
	}
	if _, err := scan.Token(); err != io.EOF {
		return false
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(profile); err != nil {
		return false
	}
	var trailing any
	return decoder.Decode(&trailing) == io.EOF
}

var readinessObjectFields = map[string]map[string]string{
	"profile":      {"schema": "", "host": "host", "vm": "vm", "legs": "legs", "network": "network", "docker": "docker", "dependencies": "dependencies", "lifecycle": "lifecycle"},
	"host":         {"id": "", "windows_version": "", "hyper_v_available": "", "hyper_v_evidence_ref": ""},
	"vm":           {"id": "", "disposable": "", "nested_virtualization": "", "wsl2_available": "", "snapshot_id": "", "prerequisite_evidence_ref": ""},
	"legs":         {"windows_native_evidence_ref": "", "wsl2_evidence_ref": ""},
	"network":      {"boundary_id": "", "deny_by_default": "", "enforced_outside_guest": "", "independent_capture": "", "enforcement_evidence_ref": "", "capture_evidence_ref": ""},
	"docker":       {"daemon_id": "", "disposable": "", "confined_to_vm": "", "credentials_isolated": "", "host_socket_exposed": "", "host_mounts_allowed": "", "privileged_containers": "", "run_owned_process_boundary": "", "boundary_evidence_ref": ""},
	"dependencies": {"acquisition_boundary_id": "", "acquisition_evidence_ref": "", "tools": "tool_pin[]", "images": "image_pin[]", "actions": "action_pin[]"},
	"tool_pin":     {"name": "", "version": "", "sha256": ""},
	"image_pin":    {"name": "", "reference": "", "platform": "", "digest": ""},
	"action_pin":   {"name": "", "repository": "", "commit": "", "tree_sha": "", "content_sha256": ""},
	"lifecycle":    {"timeout_seconds": "", "resource_owner_id": "", "cleanup_policy": "", "rollback_snapshot_id": "", "evidence_custody_id": "", "evidence_ref": ""},
}

func scanReadinessObject(decoder *json.Decoder, depth int, shape string) error {
	if depth > 32 {
		return io.ErrUnexpectedEOF
	}
	fields, ok := readinessObjectFields[shape]
	if !ok {
		return io.ErrUnexpectedEOF
	}
	seen := make(map[string]struct{})
	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			return err
		}
		key, ok := token.(string)
		if !ok {
			return io.ErrUnexpectedEOF
		}
		if _, duplicate := seen[key]; duplicate {
			return io.ErrUnexpectedEOF
		}
		childShape, known := fields[key]
		if !known {
			return io.ErrUnexpectedEOF
		}
		seen[key] = struct{}{}
		value, err := decoder.Token()
		if err != nil {
			return err
		}
		if err := scanReadinessValue(decoder, value, depth+1, childShape); err != nil {
			return err
		}
	}
	end, err := decoder.Token()
	if err != nil || end != json.Delim('}') {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func scanReadinessValue(decoder *json.Decoder, value any, depth int, shape string) error {
	if depth > 32 {
		return io.ErrUnexpectedEOF
	}
	if strings.HasSuffix(shape, "[]") {
		array, ok := value.(json.Delim)
		if !ok || array != '[' {
			return io.ErrUnexpectedEOF
		}
		for decoder.More() {
			element, err := decoder.Token()
			if err != nil {
				return err
			}
			if err := scanReadinessValue(decoder, element, depth+1, strings.TrimSuffix(shape, "[]")); err != nil {
				return err
			}
		}
		end, err := decoder.Token()
		if err != nil || end != json.Delim(']') {
			return io.ErrUnexpectedEOF
		}
		return nil
	}
	if shape != "" {
		object, ok := value.(json.Delim)
		if !ok || object != '{' {
			return io.ErrUnexpectedEOF
		}
		return scanReadinessObject(decoder, depth+1, shape)
	}
	if _, ok := value.(json.Delim); ok {
		return io.ErrUnexpectedEOF
	}
	return nil
}
