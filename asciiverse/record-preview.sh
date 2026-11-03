#!/bin/sh
# Records the two catalog previews to hub/RECORDING.md (100x30) with hub's
# scripted recorder, hub/tools/castrec.py:
#   asciiverse/wilsons-revenge.cast: Wilson's Revenge in --demo, from the
#     title screen into a run, then a key to end it
#   asciiverse/wfc-worlds.cast: landscapes and cities growing in turn
#
#   go build -o build/ ./cmd/... && asciiverse/record-preview.sh [path to castrec.py]
set -eu
cd "$(dirname "$0")/.."
rec=${1:-../../hub/tools/castrec.py}
PATH="$PWD/build:$PATH" python3 "$rec" asciiverse/wilsons-revenge.cast build "w34;kq;w1" -- wilsons-revenge --demo
PATH="$PWD/build:$PATH" python3 "$rec" asciiverse/wfc-worlds.cast build "w34;kq;w1" -- wfc-worlds
