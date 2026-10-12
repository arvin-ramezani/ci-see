// Package d3preflight validates D3 run manifests without redeeming or executing them.
package d3preflight

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/netip"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	ManifestSchema = "ci-see.d3.run-manifest.v1"
	ApprovalSchema = "ci-see.d3.approval.v1"
	MaxManifest    = 1 << 20
	MaxApproval    = 16 << 10
	MaxValidity    = 8 * time.Hour
	maxTimeout     = 8 * time.Hour
)

const signatureDomain = "CI-SEE-D3-MANIFEST-V1\x00"

type Status string

const (
	StatusConfigValid Status = "CONFIG_VALID"
	StatusBlocked     Status = "BLOCKED"
	StatusIncomplete  Status = "INCOMPLETE"
)

// Result contains a safe machine-readable classification and never includes input bytes.
type Result struct {
	Status Status
	Code   string
}

// Clock is supplied by the caller. A clock is not trusted merely because it returns a time.
type Clock struct {
	Now     time.Time
	Trusted bool
}

// Scope is supplied independently of the manifest and must match every bound identity.
type Scope struct {
	Repository   string
	PullRequest  int
	HeadSHA      string
	TreeSHA      string
	RunID        string
	HostID       string
	VMID         string
	SnapshotID   string
	GuestID      string
	ExecutionLeg string
}

type Input struct {
	ManifestBytes []byte
	ApprovalBytes []byte
	TrustedKeys   map[string]ed25519.PublicKey
	Expected      Scope
	Clock         Clock
}

type Manifest struct {
	Schema           string        `json:"schema"`
	Identity         Identity      `json:"identity"`
	Environment      Environment   `json:"environment"`
	Workflow         Workflow      `json:"workflow"`
	Artifacts        []Artifact    `json:"artifacts"`
	Commands         []Command     `json:"commands"`
	Destinations     []Destination `json:"destinations"`
	Cleanup          Cleanup       `json:"cleanup"`
	CredentialPolicy string        `json:"credential_policy"`
}

type Identity struct {
	Repository  string `json:"repository"`
	PullRequest int    `json:"pull_request"`
	HeadSHA     string `json:"head_sha"`
	TreeSHA     string `json:"tree_sha"`
	RunID       string `json:"run_id"`
	Nonce       string `json:"nonce"`
	IssuedAt    string `json:"issued_at"`
	ExpiresAt   string `json:"expires_at"`
}

type Environment struct {
	HostID                string `json:"host_id"`
	VMID                  string `json:"vm_id"`
	SnapshotID            string `json:"snapshot_id"`
	GuestID               string `json:"guest_id"`
	ExecutionLeg          string `json:"execution_leg"`
	WorkspaceSHA256       string `json:"workspace_sha256"`
	EnvironmentProfileSHA string `json:"environment_profile_sha256"`
}

type Workflow struct {
	Path          string `json:"path"`
	ContentSHA256 string `json:"content_sha256"`
	EventSHA256   string `json:"event_sha256"`
}

type Artifact struct {
	ID            string `json:"id"`
	Kind          string `json:"kind"`
	Name          string `json:"name"`
	Source        string `json:"source"`
	Revision      string `json:"revision"`
	TreeSHA       string `json:"tree_sha,omitempty"`
	ContentSHA256 string `json:"content_sha256"`
	OCIDigest     string `json:"oci_digest,omitempty"`
	Platform      string `json:"platform,omitempty"`
}

type Command struct {
	Stage          string   `json:"stage"`
	ArtifactID     string   `json:"artifact_id"`
	Argv           []string `json:"argv"`
	TimeoutSeconds int64    `json:"timeout_seconds"`
}

type Destination struct {
	Stage    string `json:"stage"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Purpose  string `json:"purpose"`
}

type Cleanup struct {
	RunLabel          string `json:"run_label"`
	RollbackSnapshot  string `json:"rollback_snapshot_id"`
	EvidenceCustodyID string `json:"evidence_custody_id"`
}

type Approval struct {
	Schema         string `json:"schema"`
	KeyID          string `json:"key_id"`
	ManifestSHA256 string `json:"manifest_sha256"`
	Signature      string `json:"signature"`
}

// Validate performs offline syntax, signature, scope, and expiration checks only.
// It never claims runtime authorization and has no redemption interface.
func Validate(input Input) Result {
	if len(input.ManifestBytes) == 0 || len(input.ApprovalBytes) == 0 {
		return result(StatusIncomplete, "approval_input_missing")
	}
	if len(input.ManifestBytes) > MaxManifest || len(input.ApprovalBytes) > MaxApproval {
		return result(StatusBlocked, "input_too_large")
	}

	var manifest Manifest
	if !decodeStrict(input.ManifestBytes, &manifest, "manifest") || !validManifest(manifest) {
		return result(StatusBlocked, "manifest_invalid")
	}
	var approval Approval
	if !decodeStrict(input.ApprovalBytes, &approval, "approval") || !validApproval(approval) {
		return result(StatusBlocked, "approval_invalid")
	}

	digest := sha256.Sum256(input.ManifestBytes)
	digestHex := hex.EncodeToString(digest[:])
	if approval.ManifestSHA256 != digestHex {
		return result(StatusBlocked, "manifest_digest_mismatch")
	}

	key, ok := input.TrustedKeys[approval.KeyID]
	if !ok || len(key) != ed25519.PublicKeySize {
		return result(StatusBlocked, "trusted_key_unavailable")
	}
	signature, err := base64.StdEncoding.Strict().DecodeString(approval.Signature)
	if err != nil || len(signature) != ed25519.SignatureSize {
		return result(StatusBlocked, "signature_invalid")
	}
	signed := make([]byte, 0, len(signatureDomain)+len(input.ManifestBytes))
	signed = append(signed, signatureDomain...)
	signed = append(signed, input.ManifestBytes...)
	if !ed25519.Verify(key, signed, signature) {
		return result(StatusBlocked, "signature_invalid")
	}

	if !validScope(input.Expected) {
		return result(StatusIncomplete, "expected_scope_missing")
	}
	if !scopeMatches(manifest, input.Expected) {
		return result(StatusBlocked, "scope_mismatch")
	}
	if !input.Clock.Trusted || input.Clock.Now.IsZero() || input.Clock.Now.Location() != time.UTC {
		return result(StatusBlocked, "trusted_utc_clock_unavailable")
	}

	issued, _ := time.Parse(time.RFC3339, manifest.Identity.IssuedAt)
	expires, _ := time.Parse(time.RFC3339, manifest.Identity.ExpiresAt)
	now := input.Clock.Now
	if issued.After(now) || !expires.After(now) {
		return result(StatusBlocked, "approval_expired_or_not_yet_valid")
	}
	if expires.Sub(issued) > MaxValidity {
		return result(StatusBlocked, "approval_validity_too_long")
	}
	return result(StatusConfigValid, "manifest_validated_offline")
}

func result(status Status, code string) Result { return Result{Status: status, Code: code} }

var jsonObjectFields = map[string]map[string]string{
	"manifest": {
		"schema": "", "identity": "identity", "environment": "environment", "workflow": "workflow",
		"artifacts": "artifact[]", "commands": "command[]", "destinations": "destination[]", "cleanup": "cleanup", "credential_policy": "",
	},
	"approval": {"schema": "", "key_id": "", "manifest_sha256": "", "signature": ""},
	"identity": {
		"repository": "", "pull_request": "", "head_sha": "", "tree_sha": "", "run_id": "", "nonce": "", "issued_at": "", "expires_at": "",
	},
	"environment": {
		"host_id": "", "vm_id": "", "snapshot_id": "", "guest_id": "", "execution_leg": "", "workspace_sha256": "", "environment_profile_sha256": "",
	},
	"workflow": {"path": "", "content_sha256": "", "event_sha256": ""},
	"artifact": {
		"id": "", "kind": "", "name": "", "source": "", "revision": "", "tree_sha": "", "content_sha256": "", "oci_digest": "", "platform": "",
	},
	"command":     {"stage": "", "artifact_id": "", "argv": "string[]", "timeout_seconds": ""},
	"destination": {"stage": "", "host": "", "port": "", "protocol": "", "purpose": ""},
	"cleanup":     {"run_label": "", "rollback_snapshot_id": "", "evidence_custody_id": ""},
}

func decodeStrict(data []byte, dst any, shape string) bool {
	if len(data) == 0 || !utf8.Valid(data) || bytes.HasPrefix(data, []byte{0xef, 0xbb, 0xbf}) {
		return false
	}
	scan := json.NewDecoder(bytes.NewReader(data))
	token, err := scan.Token()
	delim, ok := token.(json.Delim)
	if err != nil || !ok || delim != '{' || scanObject(scan, 1, shape) != nil {
		return false
	}
	if _, err := scan.Token(); err != io.EOF {
		return false
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return false
	}
	var extra any
	return decoder.Decode(&extra) == io.EOF
}

func scanObject(decoder *json.Decoder, depth int, shape string) error {
	if depth > 64 {
		return io.ErrUnexpectedEOF
	}
	fields, knownShape := jsonObjectFields[shape]
	if !knownShape {
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
		childShape, knownField := fields[key]
		if !knownField {
			return io.ErrUnexpectedEOF
		}
		seen[key] = struct{}{}
		value, err := decoder.Token()
		if err != nil {
			return err
		}
		if err := scanValue(decoder, value, depth+1, childShape); err != nil {
			return err
		}
	}
	end, err := decoder.Token()
	if err != nil || end != json.Delim('}') {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func scanArray(decoder *json.Decoder, depth int, elementShape string) error {
	if depth > 64 {
		return io.ErrUnexpectedEOF
	}
	for decoder.More() {
		value, err := decoder.Token()
		if err != nil {
			return err
		}
		if err := scanValue(decoder, value, depth+1, elementShape); err != nil {
			return err
		}
	}
	end, err := decoder.Token()
	if err != nil || end != json.Delim(']') {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func scanValue(decoder *json.Decoder, value any, depth int, shape string) error {
	if depth > 64 {
		return io.ErrUnexpectedEOF
	}
	delim, ok := value.(json.Delim)
	if !ok {
		return nil
	}
	switch delim {
	case '{':
		return scanObject(decoder, depth+1, shape)
	case '[':
		if !strings.HasSuffix(shape, "[]") {
			return io.ErrUnexpectedEOF
		}
		return scanArray(decoder, depth+1, strings.TrimSuffix(shape, "[]"))
	default:
		return io.ErrUnexpectedEOF
	}
}

func validApproval(approval Approval) bool {
	if approval.Schema != ApprovalSchema || !validID(approval.KeyID) || !lowerHex(approval.ManifestSHA256, 64) {
		return false
	}
	signature, err := base64.StdEncoding.Strict().DecodeString(approval.Signature)
	return err == nil && len(signature) == ed25519.SignatureSize && base64.StdEncoding.EncodeToString(signature) == approval.Signature
}

func validManifest(manifest Manifest) bool {
	if manifest.Schema != ManifestSchema || !validIdentity(manifest.Identity) || !validEnvironment(manifest.Environment) || !validWorkflow(manifest.Workflow) {
		return false
	}
	if manifest.Artifacts == nil || manifest.Commands == nil || manifest.Destinations == nil || len(manifest.Artifacts) == 0 || len(manifest.Commands) == 0 {
		return false
	}
	if manifest.CredentialPolicy != "dummy-only" || !validID(manifest.Cleanup.RunLabel) ||
		manifest.Cleanup.RunLabel != "ci-see-d3-"+manifest.Identity.RunID ||
		manifest.Cleanup.RollbackSnapshot != manifest.Environment.SnapshotID || !validID(manifest.Cleanup.EvidenceCustodyID) {
		return false
	}
	artifacts := make(map[string]Artifact, len(manifest.Artifacts))
	for _, artifact := range manifest.Artifacts {
		if !validArtifact(artifact) {
			return false
		}
		if _, exists := artifacts[artifact.ID]; exists {
			return false
		}
		artifacts[artifact.ID] = artifact
	}
	for _, command := range manifest.Commands {
		if !validStage(command.Stage) || command.TimeoutSeconds < 1 || command.TimeoutSeconds > int64(maxTimeout/time.Second) || len(command.Argv) == 0 || len(command.Argv) > 128 {
			return false
		}
		artifact, exists := artifacts[command.ArtifactID]
		if !exists || artifact.Kind != "tool" || command.Argv[0] != artifact.Name {
			return false
		}
		for _, arg := range command.Argv {
			if len(arg) > 8192 || !validText(arg) {
				return false
			}
		}
	}
	for _, destination := range manifest.Destinations {
		if !validStage(destination.Stage) || !validHost(destination.Host) || destination.Port < 1 || destination.Port > 65535 ||
			(destination.Protocol != "tcp" && destination.Protocol != "udp") || !validID(destination.Purpose) {
			return false
		}
	}
	return true
}

func validIdentity(identity Identity) bool {
	if !validRepository(identity.Repository) || identity.PullRequest < 1 || !lowerHex(identity.HeadSHA, 40) || !lowerHex(identity.TreeSHA, 40) ||
		!validUUIDv4(identity.RunID) || !lowerHex(identity.Nonce, 64) {
		return false
	}
	issued, err := time.Parse(time.RFC3339, identity.IssuedAt)
	if err != nil || issued.Format(time.RFC3339) != identity.IssuedAt || !strings.HasSuffix(identity.IssuedAt, "Z") {
		return false
	}
	expires, err := time.Parse(time.RFC3339, identity.ExpiresAt)
	return err == nil && expires.Format(time.RFC3339) == identity.ExpiresAt && strings.HasSuffix(identity.ExpiresAt, "Z") && expires.After(issued)
}

func validEnvironment(environment Environment) bool {
	return validID(environment.HostID) && validID(environment.VMID) && validID(environment.SnapshotID) && validID(environment.GuestID) &&
		(environment.ExecutionLeg == "windows-native" || environment.ExecutionLeg == "wsl2") &&
		lowerHex(environment.WorkspaceSHA256, 64) && lowerHex(environment.EnvironmentProfileSHA, 64)
}

func validWorkflow(workflow Workflow) bool {
	if !lowerHex(workflow.ContentSHA256, 64) || !lowerHex(workflow.EventSHA256, 64) || workflow.Path == "" || strings.Contains(workflow.Path, "\\") || strings.HasPrefix(workflow.Path, "/") {
		return false
	}
	for _, segment := range strings.Split(workflow.Path, "/") {
		if segment == "" || segment == "." || segment == ".." || !validText(segment) {
			return false
		}
	}
	return true
}

func validArtifact(artifact Artifact) bool {
	if !validID(artifact.ID) || !validID(artifact.Name) || !validSource(artifact.Source) || !lowerHex(artifact.ContentSHA256, 64) {
		return false
	}
	switch artifact.Kind {
	case "tool":
		return validID(artifact.Revision) && artifact.TreeSHA == "" && artifact.OCIDigest == "" && artifact.Platform == ""
	case "action":
		return lowerHex(artifact.Revision, 40) && lowerHex(artifact.TreeSHA, 40) && artifact.OCIDigest == "" && artifact.Platform == ""
	case "image":
		if len(artifact.OCIDigest) != len("sha256:")+64 || !strings.HasPrefix(artifact.OCIDigest, "sha256:") || !lowerHex(strings.TrimPrefix(artifact.OCIDigest, "sha256:"), 64) || !validPlatform(artifact.Platform) {
			return false
		}
		return artifact.ContentSHA256 == strings.TrimPrefix(artifact.OCIDigest, "sha256:") && artifact.Revision == "" && artifact.TreeSHA == ""
	default:
		return false
	}
}

func validSource(source string) bool {
	return source != "" && len(source) <= 2048 && strings.TrimSpace(source) == source &&
		!strings.ContainsAny(source, " \t\r\n\\@") && validText(source)
}

func validStage(stage string) bool {
	switch stage {
	case "acquisition", "metadata", "execution", "teardown":
		return true
	default:
		return false
	}
}

func validHost(host string) bool {
	if host == "" || len(host) > 253 || strings.ContainsAny(host, "*/\\%") || !validText(host) {
		return false
	}
	if address, err := netip.ParseAddr(host); err == nil {
		return address.Zone() == "" && address.String() == host
	}
	if strings.ToLower(host) != host {
		return false
	}
	for _, label := range strings.Split(host, ".") {
		if len(label) == 0 || len(label) > 63 || label[0] == '-' || label[len(label)-1] == '-' {
			return false
		}
		for _, ch := range label {
			if !(ch >= 'a' && ch <= 'z') && !(ch >= '0' && ch <= '9') && ch != '-' {
				return false
			}
		}
	}
	return true
}

func validPlatform(platform string) bool {
	parts := strings.Split(platform, "/")
	if len(parts) < 2 || len(parts) > 3 {
		return false
	}
	for _, part := range parts {
		if part == "" || strings.ToLower(part) != part {
			return false
		}
		for _, ch := range part {
			if !(ch >= 'a' && ch <= 'z') && !(ch >= '0' && ch <= '9') && ch != '.' && ch != '_' && ch != '-' {
				return false
			}
		}
	}
	return true
}

func validScope(scope Scope) bool {
	return validRepository(scope.Repository) && scope.PullRequest > 0 && lowerHex(scope.HeadSHA, 40) && lowerHex(scope.TreeSHA, 40) &&
		validUUIDv4(scope.RunID) && validID(scope.HostID) && validID(scope.VMID) && validID(scope.SnapshotID) && validID(scope.GuestID) &&
		(scope.ExecutionLeg == "windows-native" || scope.ExecutionLeg == "wsl2")
}

func scopeMatches(manifest Manifest, scope Scope) bool {
	identity, environment := manifest.Identity, manifest.Environment
	return identity.Repository == scope.Repository && identity.PullRequest == scope.PullRequest && identity.HeadSHA == scope.HeadSHA && identity.TreeSHA == scope.TreeSHA &&
		identity.RunID == scope.RunID && environment.HostID == scope.HostID && environment.VMID == scope.VMID && environment.SnapshotID == scope.SnapshotID &&
		environment.GuestID == scope.GuestID && environment.ExecutionLeg == scope.ExecutionLeg
}

func validRepository(repository string) bool {
	parts := strings.Split(repository, "/")
	return len(parts) == 2 && validID(parts[0]) && validID(parts[1])
}

func validUUIDv4(value string) bool {
	if len(value) != 36 || value[8] != '-' || value[13] != '-' || value[18] != '-' || value[23] != '-' || value[14] != '4' {
		return false
	}
	if value[19] != '8' && value[19] != '9' && value[19] != 'a' && value[19] != 'b' {
		return false
	}
	for index, ch := range value {
		if index == 8 || index == 13 || index == 18 || index == 23 {
			continue
		}
		if !(ch >= '0' && ch <= '9') && !(ch >= 'a' && ch <= 'f') {
			return false
		}
	}
	return true
}

func lowerHex(value string, length int) bool {
	if len(value) != length {
		return false
	}
	for _, ch := range value {
		if !(ch >= '0' && ch <= '9') && !(ch >= 'a' && ch <= 'f') {
			return false
		}
	}
	return true
}

func validID(value string) bool {
	if value == "" || len(value) > 128 || strings.TrimSpace(value) != value {
		return false
	}
	for _, ch := range value {
		if !(ch >= 'a' && ch <= 'z') && !(ch >= 'A' && ch <= 'Z') && !(ch >= '0' && ch <= '9') && ch != '.' && ch != '_' && ch != ':' && ch != '-' {
			return false
		}
	}
	return true
}

func validText(value string) bool {
	if !utf8.ValidString(value) {
		return false
	}
	for _, ch := range value {
		if ch < 0x20 || ch == 0x7f {
			return false
		}
	}
	return true
}

type DeadlineDisposition string

const (
	DeadlineWithinBudget DeadlineDisposition = "WITHIN_BUDGET"
	DeadlineExpired      DeadlineDisposition = "TIMEOUT_BLOCKED"
	DeadlineUnknown      DeadlineDisposition = "INCOMPLETE"
)

// AssessDeadline is a pure timestamp comparison; it neither waits nor cancels processes.
func AssessDeadline(now, deadline time.Time) DeadlineDisposition {
	if now.IsZero() || deadline.IsZero() || deadline.After(now.Add(maxTimeout)) {
		return DeadlineUnknown
	}
	if !now.Before(deadline) {
		return DeadlineExpired
	}
	return DeadlineWithinBudget
}

type ResourceObservation struct {
	ID         string
	OwnerRunID string
	Present    bool
}

type CleanupDisposition string

const (
	CleanupEvidenceComplete CleanupDisposition = "CLEANUP_EVIDENCE_COMPLETE"
	CleanupEvidenceBlocked  CleanupDisposition = "CLEANUP_EVIDENCE_BLOCKED"
	CleanupEvidenceUnknown  CleanupDisposition = "INCOMPLETE"
)

// AssessCleanup checks supplied observations only. It performs no deletion or resource access.
func AssessCleanup(runID string, inventoryComplete bool, resources []ResourceObservation) CleanupDisposition {
	if !validUUIDv4(runID) || !inventoryComplete {
		return CleanupEvidenceUnknown
	}
	for _, resource := range resources {
		if !validID(resource.ID) || resource.OwnerRunID != runID {
			return CleanupEvidenceBlocked
		}
		if resource.Present {
			return CleanupEvidenceUnknown
		}
	}
	return CleanupEvidenceComplete
}

