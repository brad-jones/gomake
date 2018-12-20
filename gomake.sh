#!/usr/bin/env bash
set -euo pipefail;
rm -f ./example/.gomake/makefile_generated.go;
rm -f ./example/.gomake/runner;
make build > /dev/null;
cd ./example/;
exec ../dist/gomake "$@";
