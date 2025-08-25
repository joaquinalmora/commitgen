# commitgen

Developer notes

- Tests are intended to be local-only and test files are ignored by git (see `.gitignore`).

[![CI](https://github.com/joaquinalmora/commitgen/actions/workflows/ci.yml/badge.svg)](https://github.com/joaquinalmora/commitgen/actions/workflows/ci.yml)

Run tests locally:

```bash
go test ./...
```

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
