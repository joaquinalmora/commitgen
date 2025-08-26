# Quick start — commitgen

This file shows the minimal steps a new user should run to try commitgen locally.

1. Build the binary

```bash
go build -o bin/commitgen ./cmd/commitgen
```

2. (Optional) Install the zsh snippet for persistent inline suggestions

```bash
# preferred command
./bin/commitgen install-shell
```

3. Run a quick diagnostic

```bash
./bin/commitgen doctor
```

4. Test in a repo with staged changes

```bash
# in any repo with staged changes
./bin/commitgen suggest --plain
# or interactively in zsh: type `git commit -m "` and accept the ghost text
```

Notes

- If you don't use zsh or do not have `zsh-autosuggestions`, the native fallback will show a preview and requires explicit accept (Ctrl+F or right-arrow depending on keymap).
- AI mode is not enabled by default and will be added behind an opt-in flag/environment variable.
