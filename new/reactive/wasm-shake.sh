#!/bin/bash
# Tree-shake wasm_exec.js to remove unused syscall/js functions
# Run this after updating TinyGo to regenerate the slim wasm_exec.js

set -e

SCRIPT_DIR=$(dirname "$(realpath "$0")")
OUTPUT="$SCRIPT_DIR/wasm_exec.js"

# Get TinyGo's wasm_exec.js
TINYGOROOT=$(tinygo env TINYGOROOT 2>/dev/null) || {
    echo "Error: TinyGo not found"
    exit 1
}

SOURCE="$TINYGOROOT/targets/wasm_exec.js"
if [ ! -f "$SOURCE" ]; then
    echo "Error: $SOURCE not found"
    exit 1
fi

echo "Source: $SOURCE ($(wc -c < "$SOURCE") bytes)"

# Remove unused functions:
# - copyBytesToGo / copyBytesToJS (we don't copy byte arrays)
# - valueDelete (we don't delete properties)
# - valueSetIndex (we don't set array indices)
# - valueInvoke (we use valueCall instead)
# - valueNew (we don't call constructors with 'new')
# - valueInstanceOf (we don't check instanceof)

cat "$SOURCE" | \
  sed '/\/\/ func valueDelete/,/^[[:space:]]*},$/d' | \
  sed '/\/\/ valueSetIndex/,/^[[:space:]]*},$/d' | \
  sed '/\/\/ func valueInvoke/,/^[[:space:]]*},$/d' | \
  sed '/\/\/ func valueNew/,/^[[:space:]]*},$/d' | \
  sed '/\/\/ func valueInstanceOf/,/^[[:space:]]*},$/d' | \
  sed '/\/\/ func copyBytesToGo/,/^[[:space:]]*},$/d' | \
  sed '/\/\/ copyBytesToJS/,/^[[:space:]]*},$/d' \
  > "$OUTPUT"

echo "Output: $OUTPUT ($(wc -c < "$OUTPUT") bytes)"
echo "Removed: valueDelete, valueSetIndex, valueInvoke, valueNew, valueInstanceOf, copyBytesToGo, copyBytesToJS"
