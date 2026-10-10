# S0 — Requirement-to-slice and test-fixture matrix

**Status:** Proposed traceability; source wording and requirement strength remain in the linked contracts.
**Sources:** [PRD](../prd.md#6-functional-requirements), [local CI acceptance](../specs/local-ci-execution.md#14-required-acceptance-tests), [Git acceptance](../specs/git-gating.md#15-required-acceptance-tests), [approval acceptance](../specs/developer-approval-ux.md#14-required-acceptance-tests).
**Notation:** S = [planned slice](../implementation-plan/vertical-slices.md#4-dependency-ordered-vertical-slices); fixture codes = [planned scenarios](fixtures.md), not tests that have run. D-gates are [defined separately](../implementation-plan/authority-and-gates.md#3-decisions-and-feasibility-gates).

## PRD functional requirements (13/13)

| Requirement | Implementing slices | Planned evidence |
| --- | --- | --- |
| FR-1 | S1, S2, S5 | R01, W01, AC01, HK01 |
| FR-2 | S2 | W01, W02 |
| FR-3 | S1, S4 | AC01, AC02, HOST02 |
| FR-4 | S2, S4 | EV01, EV02, W02 |
| FR-5 | S3, S4 | ST01, ST02, AC05 |
| FR-6 | S2, S3, S5 | SN01, SN02, FP01, HK02 |
| FR-7 | S5 | HK01, HK02, HK03 |
| FR-8 | S6, S7 | AP01, AP03, AP05, AP07; D1 |
| FR-9 | S6, S7 | AP02, AP04, AP06, AP07; D1 |
| FR-10 | S8 | REM01; **D2 decision required** |
| FR-11 | S3, S4, S5 | AC02, AC05, ST01, HK04 |
| FR-12 | S3, S4 | FP02, AC04 |
| FR-13 | S1, S6, S9, S10 | HOST01, HOST02, AP07; D4 |

## Local CI execution — Required Acceptance Tests (12/12)

| Clause | Implementing slices | Planned evidence |
| --- | --- | --- |
| L-01 | S2 | W01 |
| L-02 | S2 | EV01 |
| L-03 | S2 | EV02 |
| L-04 | S4 | AC01 |
| L-05 | S4 | AC02 |
| L-06 | S4 | AC03 |
| L-07 | S4 | AC04 |
| L-08 | S3, S4 | ST02 |
| L-09 | S2, S4 | W02, EV02 |
| L-10 | S4 | AC05 |
| L-11 | S9 | HOST01 |
| L-12 | S1, S4 | HOST02 |

## Git gating — Required Acceptance Tests (20/20)

| Clause | Implementing slices | Planned evidence |
| --- | --- | --- |
| G-01 | S5 | HK01 |
| G-02 | S5 | HK01 |
| G-03 | S5 | HK01 |
| G-04 | S5 | HK01 |
| G-05 | S3, S5 | SN01, FP01 |
| G-06 | S2, S5 | SN01 |
| G-07 | S5 | HK02 |
| G-08 | S5 | HK03 |
| G-09 | S5 | HK04 |
| G-10 | S5 | SN02 |
| G-11 | S5 | HK03, SN02 |
| G-12 | S3, S5 | FP01, SN02 |
| G-13 | S3, S5 | FP01 |
| G-14 | S3, S5 | FP02 |
| G-15 | S5 | SN03, HK03 |
| G-16 | S3 | ST01, ST02 |
| G-17 | S3, S5 | ST02, HK04 |
| G-18 | S5 | HK04 |
| G-19 | S5 | HK05 |
| G-20 | S1, S5 | HK04 |

## Developer approval UX — Required Acceptance Tests (23/23)

| Clause | Implementing slices | Planned evidence |
| --- | --- | --- |
| A-01 | S7 | AP01 |
| A-02 | S7 | AP01 |
| A-03 | S7 | AP02 |
| A-04 | S7 | AP02 |
| A-05 | S7 | AP03 |
| A-06 | S7 | AP03, AP06 |
| A-07 | S7 | AP04 |
| A-08 | S7 | AP04 |
| A-09 | S7 | AP04 |
| A-10 | S7 | AP05 |
| A-11 | S7 | AP05 |
| A-12 | S7 | AP05 |
| A-13 | S7 | AP06 |
| A-14 | S7 | AP03 |
| A-15 | S6, S7 | AP07 |
| A-16 | S6, S7 | AP02, AP07 |
| A-17 | S6, S7 | AP04 |
| A-18 | S6, S7 | AP07 |
| A-19 | S6, S7 | AP05, AP07 |
| A-20 | S6 | AP07 |
| A-21 | S7 | AP06 |
| A-22 | S7 | AP06 |
| A-23 | S7 | AP04, AP06 |

## Coverage audit boundary

The three acceptance tables preserve the **original section-local clause numbers**; L/G/A are disambiguating trace IDs, **not** newly authoritative requirement IDs. Map 68 distinct obligations (13 + 12 + 20 + 23). All rows require implementation evidence before their owning slice closes. S0 creates **no evidence of a passing Go fixture**. D1/D2/D3/D4 decisions and platform/reliability invariants must be examined across slices even if they do not have their own FR number.
