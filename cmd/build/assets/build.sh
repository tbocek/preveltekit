#!/bin/bash
set -e

PROJECT_DIR="."

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            echo "Usage: $0 [project-dir]"
            exit 0
            ;;
        -*) echo "Unknown option: $1"; exit 1 ;;
        *) PROJECT_DIR="$1"; shift ;;
    esac
done

# Resolve dependencies if needed
if [ ! -f "$PROJECT_DIR/go.sum" ] || ! grep -q 'github.com/tbocek/preveltekit/v2' "$PROJECT_DIR/go.sum"; then
    echo "Resolving dependencies..."
    (cd "$PROJECT_DIR" && go get github.com/tbocek/preveltekit/v2@latest)
fi

echo "Cleaning dist folder..."
rm -rf "$PROJECT_DIR/dist"
mkdir -p "$PROJECT_DIR/dist"

echo "Generating HTML files..."
(cd "$PROJECT_DIR" && go run -tags '!wasm' .) 2>&1 | while read -r line; do
    if [[ "$line" == Generated:* ]]; then
        echo "  ${line#Generated: }"
    else
        echo "$line"
    fi
done

echo "Building WASM..."
TINYGO_FLAGS="-target wasm -scheduler=asyncify -gc=leaking -panic=trap -no-debug"
# Compile out codec support for reflect kinds the app never decodes/encodes.
# On failure, build with all kinds.
CODEC_TAGS=$(cd "$PROJECT_DIR" && go run github.com/tbocek/preveltekit/v2/cmd/reflectscan .) || CODEC_TAGS=""
if [ -n "$CODEC_TAGS" ]; then
    echo "  codec build tags: $CODEC_TAGS"
    TINYGO_FLAGS="$TINYGO_FLAGS -tags=$CODEC_TAGS"
fi
# tinygo aborts with "requires go version ..." if it does not support the installed Go
(cd "$PROJECT_DIR" && tinygo build -o dist/main.wasm $TINYGO_FLAGS .)

echo "Copying wasm_exec.js..."
cp "$PROJECT_DIR/assets/wasm_exec.js" "$PROJECT_DIR/dist/"

echo "Strip wasm"
wasm-strip "$PROJECT_DIR/dist/main.wasm"
echo "Strip wasm_exec"
(cd "$PROJECT_DIR" && go run github.com/tbocek/preveltekit/v2/cmd/stripwasmexec dist/main.wasm dist/wasm_exec.js) || echo "  wasm_exec.js left unstripped"
echo "Compressing..."
for f in "$PROJECT_DIR/dist"/*.html "$PROJECT_DIR/dist"/*.wasm "$PROJECT_DIR/dist"/*.js; do
    [ -f "$f" ] || continue
    (zopfli --i10 "$f" 2>/dev/null || gzip -k -f "$f" 2>/dev/null || true) &
    (brotli -q 11 -k -f "$f" 2>/dev/null || true) &
done
wait
echo "Done"

echo "Build complete! Output in: $PROJECT_DIR/dist/"
