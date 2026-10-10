# CI See

Local-first GitHub Actions execution with developer-controlled CI gating.

## Documentation and AI agents

Use [AGENTS.md](AGENTS.md) and the [documentation index](docs/index.md) to select
relevant project contracts. The repository-local
[Context Engineering skill](.agents/skills/context-engineering/SKILL.md) is installed
from [ai-skills PR #6](https://github.com/arvin-ramezani/ai-skills/pull/6)
(merge commit `6ec528e498b51a4e1a9796bb55469bc188e5fdc1`). Its [local policy](docs/context-policy.md)
requires meaning-preserving documentation changes; it does not change product scope.

Read the [PRD](docs/prd.md), [architecture](docs/architecture.md),
[three behavior specs](docs/index.md), or [implementation plan](docs/implementation-plan.md)
through their task-specific indexes. [Migration evidence](docs/migration/index.md)
records the unchanged source sections and legacy anchors.

To validate docs locally (Python standard library only):

```bash
python3 scripts/validate_docs.py --check
```

## Go CLI foundation (S1)

Requires **Go 1.27.2** for development and CI. CI See's own runtime does not
require Node.js, npm, `act`, Docker or Git for these S1 informational commands.

```sh
go build -o ci-see ./cmd/ci-see
go test ./...
go vet ./...
./ci-see --help
./ci-see --version
```

On Windows build with `go build -o ci-see.exe ./cmd/ci-see` and invoke
`.\ci-see.exe`. Go runs natively inside WSL2 as a Linux executable.

`ci-see`, `ci-see init`, and `ci-see status` are **unavailable** in S1 and fail
without writing hooks, configuration or validation state. They do not run CI.

| Exit | Meaning |
| --- | --- |
| `0` | Informational success (`--help`/`--version` only in S1) |
| `1` | Internal error |
| `2` | Invalid command/arguments |
| `3` | Operation not implemented |
| `4` | Reserved for future blocking Git gates |

Validation of exact-state Git identities, `act` execution, hook integration,
state persistence and developer approval are deferred to later slices. In S1,
all core decisions block, including caller-supplied `PASS`.
