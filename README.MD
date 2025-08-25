# commitgen

Developer notes

- Tests are intended to be local-only and test files are ignored by git (see `.gitignore`).

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
