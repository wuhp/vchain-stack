#!/bin/bash

ROOT_DIR=$(readlink -f $(dirname $0))

# Build
export GOPATH=${ROOT_DIR}
go get vstack
go install vstack
