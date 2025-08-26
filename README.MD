# commitgen

Developer notes

- Tests are intended to be local-only and test files are ignored by git (see `.gitignore`).

[![CI](https://github.com/joaquinalmora/commitgen/actions/workflows/ci.yml/badge.svg)](https://github.com/joaquinalmora/commitgen/actions/workflows/ci.yml)

Run tests locally:

```bash
go test ./...
```

Ghost (inline) suggestions in zsh
--------------------------------

Quick: enable ghost suggestions (no extra tools required)

- If you use zsh with `zsh-autosuggestions` the snippet will show true ghost text inline. To try it without installing anything permanently, source the snippet:

```bash
source internal/shell/commitgen.zsh
```

- In any git repo with staged changes, type:

```bash
git commit -m "
```

Pause — a suggested single-line commit subject should appear as ghost text. Accept it with your usual autosuggest accept key (right-arrow) or with Ctrl+F (fallback binding).

Install (optional)

To install the snippet so it loads automatically:

```bash
# from the repo root
go run ./cmd/commitgen install-shell
```

Privacy & safety
- The snippet only reads the working tree and staged files locally; it never uploads content.
- By default it does not modify commit messages or hooks. An optional prepare-commit-msg prefill is planned but disabled by default.



Build for local smoke tests:

```bash
go build -o bin/commitgen ./cmd/commitgen
```

Run only the end-to-end integration test:

```bash
go test ./e2e -run TestSuggestPlainIntegration -v
```

CI

- A minimal GitHub Actions workflow is provided at `.github/workflows/ci.yml`.

- Runs gofmt check and `go test ./...` on ubuntu for push/PR.

- macOS and e2e integration jobs are available as manual workflows (trigger `Run workflow` in Actions).
