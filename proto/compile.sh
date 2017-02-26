#!/usr/bin/env bash

if test ! $(which protoc); then
    echo "The ProtoBuf Compiler must be installed. Exiting."
    exit 1
fi

SCRIPT_PATH="$(dirname "$0")"

pushd "$SCRIPT_PATH/.." > /dev/null 2>&1

# Clean up old generated files
rm -f ./proto/*.pb.go

# Generate new files
protoc -I="./proto" --go_out="./proto" ./proto/*.proto

popd > /dev/null 2>&1
