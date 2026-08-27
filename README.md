# passman

[![Go 1.25](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Termux](https://img.shields.io/badge/Termux-friendly-brightgreen?logo=android)](https://termux.dev)

Age-encrypted password manager for Termux — SQLite WAL, `bubbletea` TUI, 9M.

I built `passman` because Termux has no good offline password manager that fits in a phone's terminal. It encrypts with `age` `scrypt`, stores in SQLite, and never leaves your device.

```bash
go install github.com/Pastalikek65/passman@latest
PASSMAN_MASTER="my-secret" passman add github alice
# Password: s3cr3t
PASSMAN_MASTER="my-secret" passman get github
# s3cr3t
PASSMAN_MASTER="my-secret" passman tui   # j/k nav, q quit
```

Single file `~/.cache/passman/pass.db` WAL `0700`, survives reboot, no daemon.

## Why not `pass`

`pass` needs `gpg` and `git`, 50MB+. This is 9M, pure Go, `age` `scrypt`, offline, TUI.

## Install (Termux)

```bash
git clone https://github.com/Pastalikek65/passman.git
cd passman
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o passman .
./passman --help
```

## Security

- Master passphrase via `PASSMAN_MASTER` env or prompt (never logged)
- `age` `scrypt` (N=32768, r=8, p=1) — 2s per decrypt on phone
- Wrong passphrase → decrypt fails, no partial leak
- `0600` db file, `0700` dir

## Dev

```bash
go vet ./...
go test ./... -count=1 -timeout 30s -cover
CGO_ENABLED=0 go build -o passman .
```

## License

MIT
