**Document status: HISTORICAL — superseded for product scope by [merged PRD #20](../prd.md).**

All later status labels, authority claims, imperative wording and acceptance clauses
describe the earlier design/test scope only. They are not current CLI requirements
or permission to execute old recipes. Original text is preserved below. D3 results
remain PARTIAL/UNVERIFIED; this retirement does not establish runtime safety.

# D3 Phase B — Offline preflight verifier

**Status:** preparation implementation only; not runtime verification or execution authorization.  
**Scope:** PR #18's owner-authorized exception in [authority and gates](../implementation-plan/authority-and-gates.md).  
**Required status:** **D3 = PARTIAL; real `act` = BLOCKED; PRODUCTION_READY = NO.**

## Boundary and operations

`cmd/d3-preflight` reads bounded manifest and detached-approval files, compares them with caller-supplied expected identity, verifies an Ed25519 signature against an independently supplied trust set, and checks the signed time window against an explicitly trusted UTC clock. It also checks that an `act` execution argv uses the manifest-selected runner image and the required reviewed option shape. The command in this repository has an empty trust set and no trusted clock source, so it fails closed. Tests inject only a deterministic fixture key and fake trusted clock.

Offline validation never redeems an approval. The command has no credential-store or credential API adapter and never requests or fetches credentials. It reads the caller-supplied manifest and approval files, so callers must keep credential values out of those files; the verifier does not log their raw contents and cannot determine whether an arbitrary `argv` value is a secret. It cannot start processes, invoke `act`, execute workflows, access Docker, provision a VM, change network/firewall state, or publish evidence. A valid signature and `CONFIG_VALID` result say only that the input passed offline validation; output explicitly reports `execution_authorized=false` and `approval_redeemed=false`. The tool never emits `READY`, `PASS`, or D3 acceptance.

Replay protection belongs to a future, separately reviewed trusted runtime component. That component must atomically claim `(key_id, run_id, nonce, manifest_sha256)` in durable state outside the VM, using a transaction with a uniqueness constraint or equivalent compare-and-set. A check followed by an append is not atomic and cannot prevent concurrent replay. Duplicate claims, uncertain storage, or unavailable durable state must block. Snapshot restore must not erase redemption state.

## Signed input contract

The strict JSON manifest schema is `ci-see.d3.run-manifest.v1`. It contains:

| Object | Required fields and checks |
| --- | --- |
| `identity` | Canonical repository, positive PR number, lowercase 40-character head/tree SHA, UUIDv4 run ID, lowercase 32-byte hex nonce, canonical UTC RFC3339 `issued_at` and `expires_at`; require `issued_at <= now < expires_at` and lifetime at most eight hours. |
| `environment` | Host, VM, snapshot and guest IDs; leg is `windows-native` or `wsl2`; a run-owned `container_network` exactly bound to the cleanup run label (never `host`); workspace and environment-profile SHA-256 values. |
| `workflow` | Relative workflow path, event name, content SHA-256 and event SHA-256. |
| `artifacts` | Nonempty unique IDs. Tools pin a revision and content digest; actions pin full commit, tree and content digests; images pin OCI SHA-256 digest, OCI platform and `act_platform` selector. An image's source is the exact registry/name prefix used with its digest. Sources reject whitespace, backslashes, and URL userinfo. |
| `commands` | Nonempty exact `argv` arrays, a referenced pinned tool artifact whose name matches `argv[0]`, one of acquisition/metadata/execution/teardown stages, and a bounded timeout from one second through eight hours. A run has exactly one execution command. It must be the pinned `act` tool, reference the sole image artifact using `image_artifact_id`, enumerate every declared action artifact in `action_artifact_ids`, and match the strict option policy below. No secret values or environment map is part of the schema. Producers must not put credential values in `argv`; the verifier never logs arguments, and cannot infer whether an arbitrary opaque argument is secret. |
| `destinations` | Zero or more exact stage/host/port/protocol/purpose tuples. Wildcards and CIDRs are rejected; an empty list allows no destinations. Syntax validation does not approve a destination or prove enforcement. |
| `cleanup` | Run label bound to the run ID, rollback snapshot bound to the selected snapshot, and private evidence custody ID. |
| top level | Exact `credential_policy: dummy-only`; unknown fields and missing required arrays are rejected. |

The detached envelope schema is `ci-see.d3.approval.v1` and has `key_id`, lowercase `manifest_sha256`, and standard padded base64 `signature`. Hash the exact manifest file bytes with SHA-256. The Ed25519 message is the bytes `CI-SEE-D3-MANIFEST-V1`, followed by one NUL byte, followed by those exact manifest bytes. Whitespace changes therefore invalidate the digest or signature; there is no JSON canonicalization step.

### Offline `act` command-policy checks

The one execution-stage command must use the pinned `act` tool and explicitly reference the manifest's sole image artifact by `image_artifact_id`. Exactly one platform selector must match `<act_platform>=<source>@<oci_digest>` from that artifact. Its `action_artifact_ids` must reference every and only action artifact declared by the manifest, without duplicates. The command must bind `-W` to the manifest workflow path, name the signed workflow event, set `--directory`, provide explicit `-e`, action-cache, and empty env/secret/var/input file paths, select the run-owned `container_network`, and include `--container-daemon-socket -`, `--pull=false`, `--rebuild=false`, and `--action-offline-mode`. Its only secret argument is the run-bound dummy `GITHUB_TOKEN=CISEE-DUMMY-<run_id>`. Unreviewed `act` flags are rejected, including container-option overrides.

These are offline argv consistency checks against the reviewed Phase A command shape. They do not inspect the workflow to verify that it references exactly the declared action IDs, or confirm the action cache's bytes; a separate runtime gate must compare the pinned workflow contents, resolved action references and verified cache inventory with these IDs. The verifier also cannot prove that referenced files are empty, confirm image availability or effective Docker HostConfig/network state, or prove runtime isolation. `CONFIG_VALID` reports input, signature, scope, expiry, and command-shape validation only; it is not execution-policy approval. A separately reviewed runtime security gate and separate exact-run authorization remain mandatory.

Reject invalid UTF-8, a BOM, duplicate object keys, unknown or case-mismatched fields, malformed or trailing JSON, noncanonical encodings, invalid pins, scope mismatch, invalid timestamps, expired/future approvals, and validity longer than eight hours. CLI option duplication is also rejected. The expected identity is supplied by the caller and is only compared with the signed manifest; matching values do not authorize execution. The key ID resolves only through a separately configured trust map; a key supplied by the manifest or approval envelope is never trusted. No owner-approved key fingerprint or trusted clock source is configured, so production command runs remain **BLOCKED** even though tests verify the cryptographic and time-window code with disposable fixtures. Clock uncertainty has no skew allowance and fails closed.

## Verification and remaining gates

Offline tests cover valid fixtures, schema and duplicate-key rejection, byte changes, signature/key failures, scope mismatch, expiry boundaries, untrusted time, bounded file reads, argument failure, fail-closed output, missing/mismatched runner-image selection, missing/mismatched action-artifact bindings, host-network attempts, missing daemon-socket/pull/rebuild/offline-mode controls, and unreviewed `act` flags. Package-boundary checks reject process-launch/network runtime dependencies and forbidden API calls. Deadline and cleanup tests operate on fake timestamps and supplied resource observations; they do not terminate process trees, delete resources, or prove runtime isolation or cleanup.

Keep Phase B dependency ordered: (1) offline preparation; (2) owner-authorized disposable Windows 11 Hyper-V provisioning and nested-virtualization/WSL2 setup; (3) independent verification of the outside-guest deny-by-default boundary and traffic capture, native/WSL2 legs, isolated Docker, provenance, cancellation, cleanup, rollback, and evidence custody; (4) controlled `act` only after separate explicit authorization for exact workflows and commands; (5) private raw evidence, independent sanitization review, and D3 acceptance review. Host selection and network boundary remain unavailable, so provisioning, independent runtime verification, and real `act` are blocked.

Preserve every existing `UNVERIFIED` evidence ID and its acceptance condition in the [Phase A matrix](d3-phase-a-verification-plan.md#acceptance-matrix-for-the-26-unverified-records). The historical 26 PASS / 26 UNVERIFIED counts are not Phase B results. Change no evidence status without independent evidence for that row's existing acceptance condition. This preparation does not close S4/S5 or production acceptance.
