# Sources and Policy Provenance

Last source verification: 2026-10-10.

## Primary official guidance

- [Rethinking skills and prompts for GPT-6 Astra](https://developers.openai.com/blog/rethinking-skills-and-prompts-for-gpt-6-astra),
  September 11, 2026: narrow triggers, concise descriptions, progressive disclosure,
  contextual routing, deliberate completion boundaries.
- [Build skills](https://learn.chatgpt.com/docs/build-skills): metadata-first loading,
  optional references/scripts, explicit and implicit invocation.
- [Custom instructions with AGENTS.md](https://learn.chatgpt.com/docs/agent-configuration/agents-md):
  broader-to-narrower instruction discovery and a configurable combined byte limit.
  Recheck version-specific behavior when configuring it.
- [Testing Agent Skills Systematically with Evals](https://developers.openai.com/blog/eval-skills):
  measurable outcomes, realistic tasks, deterministic and qualitative checks.

Use the September article as primary design guidance. Retrieve current official
docs for claims about tool behavior/configuration.
Treat recommendations/examples as guidance rather than universal protocol requirements.
Respect explicit user instructions and established project contracts.

## Supplied research

Derived from the user's `deep-research-report.md`, attached as
`deep-research-report(1).md` (1,046 lines).
SHA-256: `2c4239b8f394ac1b401b5329f275e045488159e798c82f6ec150d57b2cf97f77`.

Carry forward progressive disclosure, canonical ownership, repository-local durable
knowledge, implementation synchronization, and measurable enforcement.
Adapt CI See examples to each repository's actual requirements.
Ordinary work does not require loading the full research report or every source page.
Replace opaque research-session citation IDs with verified links when promoting
externally supported claims into portable policy.

## Local conventions

The 80–100-line preference, 100–150-line cohesive allowance, mandatory review over
150 lines, and strong splitting preference for 300–400+ lines are user/project
conventions. OpenAI does not prescribe these as universal limits.
Review decisions, preservation maps, authority/status fields, and the scanner
implement the user's goal; they are not official mandates.
Retain large cohesive documents with a recorded reason.
Apply a different explicit project policy when instructed; preserve meaning.
