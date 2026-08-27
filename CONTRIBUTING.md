# Contributing to passman

```bash
git clone https://github.com/Pastalikek65/passman.git
cd passman
go vet ./...
go test ./... -count=1 -timeout 30s -cover
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o passman .
PASSMAN_MASTER=test go test ./store -v
```

PRs: fork → feat/name → test → open PR.
