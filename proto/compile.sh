#!/usr/bin/env bash

if test ! $(which protoc); then
    echo "The ProtoBuf Compiler must be installed. Exiting."
    exit 1
fi

SCRIPT_PATH="$(dirname "$0")"

pushd "$SCRIPT_PATH/.." > /dev/null 2>&1

protoc -I="./proto" --go_out="./proto" ./proto/*.proto

popd > /dev/null 2>&1
