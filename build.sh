#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

rm -rf artifacts/*

export CGO_ENABLED=0
export GOARCH="${ARCH}"
export GOOS="${OS}"


# build linux
env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o apl-loc-deploy main.go

tar -czf artifacts/apl-loc-deploy-linux.tgz apl-loc-deploy readme.txt interviews/*.json

rm apl-loc-deploy