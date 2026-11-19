# CLAUDE.md

Collection of small Go CLIs wrapping AWS APIs (aws-sdk-go-v2). See README.md for what each tool does.

## Layout

- `cmd/<tool>/main.go` — one self-contained `main` package per tool. Keep tools single-file unless shared.
- `internal/ssm` — Parameter Store helpers shared by the `ssm-*` tools.
- `vendor/` — committed. After any dependency change run `go mod tidy && go mod vendor` and commit both.

## Commands

```sh
go vet ./...
go test ./...
go build ./...
```

There are no tests that hit AWS; don't add any. Tests, if added, must run offline.

## Conventions

- Go 1.25 minimum (`go.mod`); CI tests on the last three Go majors (1.25, 1.26, 1.27); release builds use only the newest. When bumping, update `test.yml` and the `go-version` in `release.yml` together.
- Config via stdlib `flag` or `envconfig` env vars; AWS config via `config.LoadDefaultConfig`.
- Errors wrapped with `fmt.Errorf("svc: op: %w", err)`; fatal via `log.Fatal*` in `main`.
- Destructive tools only preview by default and act when `-apply` is passed. Don't add `-dry-run` flags.
- Dependencies are limited to the AWS SDK, `godotenv`, and `envconfig`. Don't add more without a strong reason.

## Adding a tool

1. Create `cmd/<tool>/main.go`.
2. Add `<tool>` to `CMDS` in `.github/workflows/release.yml`.
3. Document it in README.md.

Release binaries are built with `-ldflags "-X main.version=... -X main.revision=... -X main.date=..."`; declare those vars in `main` if the tool reports its version.
