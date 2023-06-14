#!/bin/bash

mkdir -p build

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

docker run --rm \
    -v "${SCRIPT_DIR}/app:/build/src/app" \
    -v "${SCRIPT_DIR}/build:/build/dst" \
    -e GOPATH=/build/ \
    golang:1.19.3 \
    sh -c "cd /build/src/app && GOOS=linux go build -tags=ycf -buildmode=plugin -o /build/dst ."
