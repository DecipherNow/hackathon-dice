#!/bin/bash

set -euxo pipefail

# assume we are in the service directory
SERVICEDIR=$PWD
DOCKERDIR="${SERVICEDIR}/docker"
BUILDFILES="${DOCKERDIR}/build"

mkdir -p "$BUILDFILES"

(
	cd "cmd/server"
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w' -o "$BUILDFILES/dice" .
)

(
	cp "settings.toml" "$BUILDFILES/."
	#cd $DOCKERDIR
	docker build -t deciphernow/hackathon:dice .
	rm -rf $BUILDFILES
)

