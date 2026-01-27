#!/bin/bash
# Use the main build script from the project root
cd "$(dirname "$0")"
exec ../../build.sh "$@"
