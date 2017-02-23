#!/usr/bin/env bash

SCRIPT_PATH="$(dirname "$0")"

pushd "$SCRIPT_PATH/.." > /dev/null 2>&1

protoc -I="$SCRIPT_PATH/proto" --go_out="$SCRIPT_PATH/proto" proto/*.proto

popd > /dev/null 2>&1
