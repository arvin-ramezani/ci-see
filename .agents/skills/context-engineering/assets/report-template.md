# Documentation Change Report

Result: <PASS | PARTIAL | BLOCKED>
Mode: <BOOTSTRAP | MAINTAIN | AUDIT>
Scope: <repository, areas, exclusions>
Baseline: <commit plus snapshots/hashes of affected uncommitted contents>
Compared result: <commit or working-tree snapshot/hash evidence>
Semantic review: <SELF_REVIEWED | INDEPENDENT_PASS | BLOCKED; evidence>

## Size review and changes

| Original file | Before lines | Decision | Destinations / after lines | Reason and retrieval benefit |
| --- | ---: | --- | --- | --- |
| <path> | <count> | <KEEP/SPLIT/REORGANIZE> | <paths/counts> | <cohesion or concrete tasks> |

## Content preservation

Use IDs or baseline path, heading, and line ranges. Account for all affected units.
Group unchanged units only with explicit scope.

| Unit / original locator | Canonical destination / heading | Status | Comparison evidence |
| --- | --- | --- | --- |
| <ID, condition, subclause> | <path#heading or ID> | <PRESERVED/DEDUPLICATED/AUTHORIZED_CHANGE/UNRESOLVED> | <equivalence or authorization> |

## Findings and owner decisions

| ID | Severity | Source evidence | Action / specific unresolved decision |
| --- | --- | --- | --- |
| <ID> | <BLOCKING/REQUIRED/ADVISORY> | <locator and conflict> | <fix or needed decision> |

Use `None` for empty findings/decisions. Retain resolved findings with their fix.

## Verification

| Check | Result | Actual command or comparison evidence / limitations |
| --- | --- | --- |
| File links and index coverage | <PASS/FAIL/NOT_RUN> | <command and scope> |
| Anchors and references in other file types | <PASS/FAIL/NOT_RUN> | <manual search or capable checker> |
| Semantic comparison | <PASS/FAIL/NOT_RUN> | <baseline/result, reviewer method> |
| Implementation synchronization, if applicable | <PASS/FAIL/NOT_APPLICABLE/NOT_RUN> | <code/tests checked> |

Next action: <completed, review needed, or specific blocked item>.
Keep unmatched units visible; structural PASS does not imply semantic PASS.
