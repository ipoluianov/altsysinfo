#!/usr/bin/env bash
# Builds altsysinfo for every supported platform into one bin/<yyyy-mm-dd-HH-MM>-<version>/.
set -euo pipefail
DIR="$(dirname "$0")"
BUILD_STAMP="$(date +%Y-%m-%d-%H-%M)"
export BUILD_STAMP
for t in linux/amd64 linux/arm64 windows/amd64 windows/arm64 darwin/amd64 darwin/arm64; do
  "$DIR/_build.sh" "${t%/*}" "${t#*/}"
done
