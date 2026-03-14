# Contributing to go-fantasy-pl

Thanks for helping improve this project.

## Quick Start

```bash
git clone https://github.com/AbdoAnss/go-fantasy-pl.git
cd go-fantasy-pl
go mod download
go test ./...
```

Use Go `1.23+`.

## Workflow

1. Create a branch from `main`.
2. Make focused changes.
3. Add or update tests.
4. Run checks locally.
5. Open a PR with a clear description.

Branch name examples:

- `feat/async-docs`
- `fix/cache-error-handling`
- `chore/test-cleanup`

## Local Checks

Run these before pushing:

```bash
gofmt -w .
go test ./...
```

If you have `golangci-lint` installed:

```bash
golangci-lint run
```

## Code Guidelines

- Keep APIs simple and predictable.
- Prefer small, testable functions.
- Handle errors explicitly.
- Add comments only where logic is non-obvious.
- Keep examples in `examples/` working.

## Pull Requests

Please include:

- What changed
- Why it changed
- Any breaking change notes
- Linked issue (if applicable)

PR checklist:

- Tests pass
- Docs updated when behavior changes
- No unrelated refactors

## Issues

Open an issue for bugs, feature requests, or documentation improvements.

## License

By contributing, you agree your contributions are licensed under the [MIT License](LICENSE).
