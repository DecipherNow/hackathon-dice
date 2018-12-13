#!/bin/bash

set -euxo pipefail

# assume we are in the service directory
SERVICEDIR=$PWD
SERVERDIR="${SERVICEDIR}/cmd/server"

GOPATH="${GOPATH:-$HOME/go}"
GOPATH=$(echo $GOPATH | cut -d: -f1)

(
	cd $SERVERDIR
	go build -o=$GOPATH/bin/dice
)
