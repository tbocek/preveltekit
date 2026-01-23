#!/bin/bash
set -e

if [ -z "$1" ]; then
    echo "usage: build.sh <main-component.go> [child-component.go ...]"
    exit 1
fi

DIR=$(dirname "$1")
SCRIPT_DIR=$(dirname "$(realpath "$0")")

# Always rebuild CLI to pick up latest changes
(cd "$SCRIPT_DIR/cmd" && go build -o reactive .)

# Always copy fresh wasm_exec.js from TinyGo (no tree-shaking)
cp "$(tinygo env TINYGOROOT)/targets/wasm_exec.js" "$SCRIPT_DIR/wasm_exec.js"

# Generate code
"$SCRIPT_DIR/cmd/reactive" "$@"

# Pre-render
echo "Pre-rendering..."
(cd "$DIR/build" && go build -tags '!js,!wasm' -o _render . && REACTIVE_BUILD=true ./_render > ../dist/prerendered.html && rm _render)

# Build WASM
echo "Building WASM..."
(cd "$DIR/build" && tinygo build -o ../dist/app.wasm -target wasm -no-debug -panic=trap -scheduler=asyncify -gc=leaking .)

# Minify wasm_exec.js if minify is available
if command -v minify &> /dev/null; then
    minify --type=js "$DIR/assets/wasm_exec.js" -o "$DIR/assets/wasm_exec.js"
fi

# Assemble final index.html
"$SCRIPT_DIR/cmd/reactive" --assemble "$DIR"

echo "Done! Serve: cd $DIR/dist && python3 -m http.server"
