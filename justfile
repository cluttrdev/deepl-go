GIT_DIR := `git rev-parse --show-toplevel`

MAIN := "./cmd/deepl"
BIN_NAME := "deepl"
BIN_DIR := "bin"
DIST_DIR := "dist"

# list available recipes
default:
    @just --list

# format code
fmt:
    go fmt ./...

# lint code
lint:
    golangci-lint run ./...

# vet code
vet:
    go vet ./...

# build application
build *args="":
    {{GIT_DIR}}/scripts/build.sh -p {{MAIN}} {{args}}

# create binary distribution
dist *args="":
    {{GIT_DIR}}/scripts/dist.sh -p {{MAIN}} {{args}}

# create a new release
release *args="": clean
    {{GIT_DIR}}/scripts/release.sh {{args}}

changes from="" to="":
    #!/bin/sh
    source {{GIT_DIR}}/scripts/functions.sh
    get_changes {{from}} {{to}}

clean:
    @# build artifacts
    @echo "rm {{BIN_DIR}}/{{BIN_NAME}}"
    @-[ -f {{BIN_DIR}}/{{BIN_NAME}} ] && rm {{BIN_DIR}}/{{BIN_NAME}}
    @-[ -d {{BIN_DIR}} ] && rmdir {{BIN_DIR}}

    @# distribution archives
    @echo "rm {{DIST_DIR}}/{{BIN_NAME}}_*.tar.gz"
    @rm {{DIST_DIR}}/{{BIN_NAME}}_*.tar.gz 2>/dev/null || true
    @echo "rm {{DIST_DIR}}/{{BIN_NAME}}_*.zip"
    @rm {{DIST_DIR}}/{{BIN_NAME}}_*.zip 2>/dev/null || true
    @-[ -d {{DIST_DIR}} ] && rmdir {{DIST_DIR}}

# ---

_system-info:
    @echo "{{os()}}_{{arch()}}"
