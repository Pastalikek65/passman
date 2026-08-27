# Security Policy

Supported: 0.1.x

Report via https://github.com/Pastalikek65/passman/security/advisories/new — 48h response.

- `age` `scrypt` N=32768, master via `PASSMAN_MASTER` env or prompt, never logged
- DB `0600`, dir `0700`, `age` decrypt fails with wrong passphrase
