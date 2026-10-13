# D3 Phase B — Environment readiness preparation

**Status:** Offline preparation contract only; environment readiness has not been verified.\
**Authority:** One-time preparation scope in [authority and gates](../implementation-plan/authority-and-gates.md); operational controls in the [Phase B runtime catalogue](d3-phase-a-runtime-catalogue.md); validator boundary in the [offline preflight contract](d3-phase-b-preflight.md).\
**Required status:** **D3 = PARTIAL; real `act` = BLOCKED; PRODUCTION_READY = NO.**

## Purpose and boundary

This procedure defines the selections and evidence references needed before any live D3 environment work. It does not provision infrastructure, change configuration, run `act`, or prove isolation. Keep all gates fail-closed. A profile is untrusted configuration: `CONFIG_VALID` means only that its exact bytes match the signed manifest and its declared choices satisfy the offline schema. It never means READY, PASS, execution authorization, or D3 acceptance.

The Phase A runtime catalogue remains the operational procedure for launch, containment, timeout, inspection, cleanup, and rollback details. This document records prerequisites and the evidence needed to decide whether an owner-authorized run can be considered; it does not duplicate those commands.

## Required preparation selections

1. **Host and VM:** The owner selects a disposable Windows 11 test host and a disposable Hyper-V VM. Record the Windows/Hyper-V prerequisite evidence, VM identity, clean rollback snapshot, nested-virtualization capability, and WSL2 availability. Do not infer nested virtualization or WSL2 support from host edition alone.
2. **Separate execution legs:** Plan and later verify Windows-native and WSL2 legs independently. Each leg needs its own tool/runtime identity, evidence reference, process observer, timeout result, and cleanup inventory. One leg's result cannot establish the other's readiness.
3. **Outside-guest network:** The owner selects a deny-by-default enforcement and independent traffic-capture boundary outside the guest. Record separate references for policy enforcement and capture. Capture must be independent of the guest and correlate host `act`, helpers, WSL2, the disposable daemon, and job containers. Docker `--internal`, offline flags, guest firewall rules, and absence of observed traffic do not establish host-process egress isolation.
4. **Disposable Docker:** Select a daemon confined to the disposable VM. The intended boundary disables host socket/pipe exposure, host mounts, privileged containers, inherited credentials, and unowned process/resource access. Record daemon identity and independent effective-configuration evidence. Never use the developer's daemon or credentials.
5. **Pinned dependencies:** Pin tools by version and SHA-256; images by immutable OCI digest and platform; actions by full commit, tree, and content digest. Acquire dependencies only in a separately selected and approved boundary with its own allowlist, deny-by-default enforcement, and independent capture. Transfer verified content into execution as a read-only/immutable handoff. Execution must not fetch on cache miss or fall back to mutable tags or alternate sources.
6. **Bounded lifecycle:** Set finite per-command/run timeouts and identify the run owner for every process and resource. Cleanup may target only resources proven to carry that run's unique ownership label. Record normal stop, forced-stop escalation, and post-stop inventory separately. Any uncertain owner, timeout, cleanup, or process-tree observation blocks progress; do not remove resources by broad name.
7. **Rollback and evidence custody:** Bind rollback to the selected clean VM snapshot. Keep original logs and raw captures in private run-scoped custody, with hashes and byte counts. Publish only independently reviewed sanitized evidence; never commit credentials, raw secret values, private paths, or internal addresses.

## Readiness profile and verifier

The optional `d3-preflight --readiness-profile` mode requires the existing signed manifest, detached approval, expected run scope, trusted key, and trusted clock inputs. It validates a bounded, versioned profile with duplicate-key, unknown-field, malformed-JSON, encoding, and trailing-value rejection. SHA-256 is computed over the original profile bytes and compared with the signed manifest's `environment_profile_sha256`. Profile tool/image/action pins must equal the manifest's artifact identities. See the [preflight contract](d3-phase-b-preflight.md#optional-readiness-profile-mode) for exact behavior.

The profile contains selections, policy booleans, pins, and opaque evidence references. The verifier cannot check reference existence, evidence authenticity, effective Hyper-V/WSL2/Docker/network settings, observed traffic, actual dependency bytes, cleanup, or rollback. A positive result is configuration validation only. Missing selections or evidence references are `INCOMPLETE`; unsafe or conflicting values are `BLOCKED`.

## Owner decisions and next gate

No test host or outside-guest enforcement/capture boundary is selected in this preparation. Those choices remain explicit owner decisions. Do not provision or modify a VM, Hyper-V, WSL2, Docker, firewall, route, or network; install/download dependencies; use real credentials; or run `act` or workflows without separate explicit authorization for the exact scope.

After host and network-boundary selection, obtain independent evidence that each prerequisite is effective, review the policy and capture outside the guest, and complete both platform legs' non-`act` controls. A separately reviewed runtime security gate and separate exact-run owner authorization remain mandatory before any live `act`. Preserve **D3 = PARTIAL; real `act` = BLOCKED; PRODUCTION_READY = NO** until existing acceptance evidence supports a change.
