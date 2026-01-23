#!/bin/bash
set -e

# Tree-shake wasm_exec.js based on what the WASM actually imports
wasm_shake() {
    local SOURCE="$1" WASM="$2" OUTPUT="$3"
    local USED=$(wasm-objdump -j Import -x "$WASM" 2>/dev/null | grep 'syscall/js\.' | sed 's/.*syscall\/js\.\([a-zA-Z]*\).*/\1/' | sort -u)
    cp "$SOURCE" "$OUTPUT"
    for func in valueDelete valueSetIndex valueInvoke valueNew valueInstanceOf copyBytesToGo copyBytesToJS; do
        if ! echo "$USED" | grep -q "^${func}$"; then
            sed -i "/\/\/ func ${func}/,/^[[:space:]]*},$/d" "$OUTPUT"
            sed -i "/\/\/ ${func}/,/^[[:space:]]*},$/d" "$OUTPUT"
        fi
    done
}

if [ -z "$1" ]; then
    echo "usage: build.sh <main-component.go> [child-component.go ...]"
    exit 1
fi

DIR=$(dirname "$1")
SCRIPT_DIR=$(dirname "$(realpath "$0")")

# Use project assets if they exist, otherwise fall back to default
if [ -f "$DIR/assets/index.html" ]; then
    INDEX_HTML="$DIR/assets/index.html"
else
    INDEX_HTML="$SCRIPT_DIR/assets/index.html"
fi
if [ -f "$DIR/assets/wasm_exec.js" ]; then
    WASM_EXEC="$DIR/assets/wasm_exec.js"
else
    WASM_EXEC="$SCRIPT_DIR/assets/wasm_exec.js"
fi

# Always rebuild CLI to pick up latest changes
(cd "$SCRIPT_DIR/cmd" && go build -o reactive .)

# Generate code
"$SCRIPT_DIR/cmd/reactive" "$@"

# Pre-render
echo "Pre-rendering..."
(cd "$DIR/build" && go build -tags '!js,!wasm' -o _render . && REACTIVE_BUILD=true ./_render > ../dist/prerendered.html && rm _render)

# Build WASM
echo "Building WASM..."
(cd "$DIR/build" && tinygo build -o ../dist/app.wasm -target wasm -no-debug -panic=trap -scheduler=asyncify -gc=leaking .)
wasm-strip "$DIR/dist/app.wasm"

# Tree-shake wasm_exec.js based on actual WASM imports
wasm_shake "$WASM_EXEC" "$DIR/dist/app.wasm" "$DIR/dist/wasm_exec.js"

# Minify wasm_exec.js if minify is available
if command -v minify &> /dev/null; then
    minify --type=js "$DIR/dist/wasm_exec.js" -o "$DIR/dist/wasm_exec.js"
fi

# Assemble final index.html
"$SCRIPT_DIR/cmd/reactive" --assemble "$DIR"

# Compress all dist files with gzip and brotli
for f in "$DIR/dist"/*; do
    [ -f "$f" ] || continue
    gzip -9 -k -f "$f"
    brotli -9 -k -f "$f"
done

echo "Done!"
