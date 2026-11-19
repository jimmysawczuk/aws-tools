# aws-tools

Small, single-purpose CLIs for common AWS chores in CI/CD pipelines and containers: CloudFront invalidations, ECS CodeDeploy prep, Secrets Manager reads, and SSM Parameter Store sync.

All tools use the standard AWS SDK credential chain (`AWS_PROFILE`, `AWS_REGION`, env credentials, instance/task roles, etc.).

## Install

Prebuilt binaries for `darwin-amd64`, `linux-amd64`, and `linux-arm64` are attached to each [release](https://github.com/jimmysawczuk/aws-tools/releases), built with each of Go 1.25, 1.26, and 1.27 and named `<tool>-<tag>-<os>-<arch>-go<version>`.

Or build from source (Go 1.25+):

```sh
go install github.com/jimmysawczuk/aws-tools/cmd/<tool>@latest
```

## Tools

### cloudfront-invalidate

Creates an invalidation and blocks until it completes. The first argument is a distribution ID or an alias (domain); the path defaults to `/*`.

```sh
cloudfront-invalidate E1234567890ABC
cloudfront-invalidate www.example.com '/assets/*'
```

### ecs-build-appspec

Prints a CodeDeploy ECS appspec (JSON) for a service, carrying over the service's capacity provider strategy. The task definition is left as the `<TASK_DEFINITION>` placeholder.

| Env var          | Required | Default |
| ---------------- | -------- | ------- |
| `CLUSTER_NAME`   | yes      |         |
| `SERVICE_NAME`   | yes      |         |
| `CONTAINER_NAME` | yes      |         |
| `CONTAINER_PORT` | no       | `8080`  |

```sh
CLUSTER_NAME=prod SERVICE_NAME=api CONTAINER_NAME=api ecs-build-appspec > appspec.json
```

### ecs-find-template-taskdef

Finds the newest active task definition in a family carrying a given tag and prints it as JSON, with the first container's image replaced by a placeholder. Useful for producing a `taskdef.json` for CodeDeploy from a Terraform-managed template.

| Flag     | Default                | Description                             |
| -------- | ---------------------- | --------------------------------------- |
| `-tag`   | `CreatedBy=Terraform`  | `key=value` tag identifying the template |
| `-image` | `<IMAGE1_NAME>`        | placeholder for the first container image |

```sh
ecs-find-template-taskdef api > taskdef.json
```

### ecs-prune-taskdefs

For each service family, keeps every active revision up to and including the newest one tagged `CreatedBy=Terraform`, deregisters all older active revisions, then deletes all inactive revisions. **Destructive; no dry run.**

| Env var         | Required | Description                     |
| --------------- | -------- | ------------------------------- |
| `CLUSTER_NAME`  | yes      |                                 |
| `SERVICE_NAMES` | yes      | comma-separated task def families |

### retrieve-secret

Writes a Secrets Manager secret's string value to stdout or a file.

```sh
retrieve-secret my/secret
retrieve-secret -out .env my/app/env
```

### ssm-read

Reads all parameters under a path (decrypted) and writes them as `NAME="value"` lines, with names relative to the path.

```sh
ssm-read -path /myapp/prod > .env
ssm-read -path /myapp/prod -out .env
```

### ssm-load

Writes a dotenv file into Parameter Store as `SecureString` parameters under a path, overwriting existing values. Empty values are skipped. Previews only unless `-apply` is passed.

```sh
ssm-load -path /myapp/prod .env         # preview
ssm-load -path /myapp/prod -apply .env  # write
```

### ssm-delete

Deletes all parameters under a path. Lists what would be deleted unless `-apply` is passed.

```sh
ssm-delete -path /myapp/old
ssm-delete -path /myapp/old -apply
```

### ssm-exec

Loads parameters under a path into the environment, then `exec`s the given program in place of itself. Variables already set in the environment take precedence.

```sh
ssm-exec -path /myapp/prod ./server --port 8080
SSM_PATH=/myapp/prod ssm-exec ./server --port 8080
```

`-path` takes precedence over `SSM_PATH`.

## Development

```sh
go vet ./...
go test ./...
go build ./...
```

Dependencies are vendored; run `go mod tidy && go mod vendor` after changing them. Releases are built by GitHub Actions when a GitHub release is created.
