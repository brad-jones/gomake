#!/usr/bin/env bash
set -euo pipefail;
make build > /dev/null;
cd ./example/;
exec ../dist/gomake "$@";
