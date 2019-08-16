#!/usr/bin/env bash
set -euo pipefail;
rm -rf ~/.gomake;
make build > /dev/null;
cd ./example/;
exec ../dist/gomake "$@";
