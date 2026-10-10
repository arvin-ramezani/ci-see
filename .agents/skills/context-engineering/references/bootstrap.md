# Bootstrap Repository Knowledge

Start with the existing layout, instructions, canonical contracts, and checks.
Identify accepted, proposed, historical, and absent knowledge.
Inspect relevant code/tests when establishing architecture facts; report missing
evidence instead of inferring implemented behavior from a plan.

## Establish the smallest useful system

1. Map existing product, architecture, behavioral, decision, planning, and operational
   knowledge. Preserve paths unless a move benefits retrieval.
2. Establish a concise root operating guide and documentation index.
3. Give each entry its responsibility and the task that should read it.
4. Put detailed contracts in canonical owners; link to them from the guide.
5. Add lightweight status/authority metadata where useful.
6. Add a short repository context policy with the adopted size-review, preservation,
   routing, and synchronization rules. Keep it usable without this skill installed.
7. Review documents over the adopted threshold; execute justified splits through
   the restructure workflow when the bootstrap task authorizes them.

Use [guide fragment](../assets/agents-fragment.md) and
[index template](../assets/index-template.md) as starting points.
Replace placeholders with real canonical paths; remove unused categories.
Ground approval boundaries in actual owner decisions.
Use nested operating guides only for persistent differences in a code area.

## Establish maintenance

Add task-specific documentation routing without requiring the full tree for every edit.
Require canonical updates for changed contracts in the same PR.
Require an agent size review for edited documents over the adopted threshold.
Keep completed plans and superseded ADRs distinguishable from current truth.
Choose real verification entry points and preserve existing development workflows.

## Establish checks

Reuse existing link, Markdown, schema, and architecture checks.
Run the bundled scanner when an initial inventory is needed.
Choose documentation roots, generated/archive exclusions, and intended entry points
before adopting strict reachability in CI. Explain exclusions.
If adding CI is in scope, install a versioned checker and wire its command into the
actual workflow; avoid a personal skill path on CI runners.
Record checks still requiring implementation instead of claiming enforcement exists.
Enforce architecture only through established, measurable project invariants.

## Verify

Try representative tasks against the map: a feature contract change, a boundary
change, and an operational failure when these domains exist.
Confirm each reaches canonical context with few irrelevant hops.
Check authority/status conflicts, links, coverage, and preservation after moves.
Report unsupported assumptions and genuine owner decisions; finish unaffected work.
