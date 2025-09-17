#!/bin/sh

echo "[INFO] Check go version."

go_version=$(go version 2>/dev/null)

if [ $? -ne 0 ]; then
    echo "[ERROR] Go is not installed or not in PATH."
    exit 1
else
    version_number=$(echo "$go_version" | awk '{print $3}' | sed 's/^go//')
    echo "[CHECK] Go version : $version_number"

    echo "[BUILD] Multi Cluster Changer..."

    if [ ! -f "./pkg/version/version.go" ]; then
        echo "[ERROR] ./pkg/version/version.go not found." >&2
        exit 1
    else
        VERSION=$(grep 'var Version' "./pkg/version/version.go" | sed -E 's/.*"([^"]+)".*/\1/')

        go build ./cmd/mcc

        mv mcc ./mcc

        bash ./scripts/info.sh $VERSION
    fi
fi