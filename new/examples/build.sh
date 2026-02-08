#!/bin/bash
set -e

RELEASE_MODE=false
PROJECT_DIR="."

strip_wasm_exec() {
    local wasm_file="$1"
    local wasm_exec_file="$2"

    command -v wasm-objdump &> /dev/null || return

    local imports=$(wasm-objdump -j Import -x "$wasm_file" 2>/dev/null | \
        grep -oE '<[^>]+>' | sed 's/[<>]//g' | sort -u)
    [ -z "$imports" ] && return

    local all_funcs=$(grep -oE '"(runtime\.|syscall/js\.)[^"]+":' "$wasm_exec_file" | \
        sed 's/"//g; s/:$//' | sort -u)

    local unused=""
    local removed=0
    while IFS= read -r func; do
        [ -z "$func" ] && continue
        if ! echo "$imports" | grep -qF "$func"; then
            unused="${unused}${func}"$'\n'
            removed=$((removed + 1))
        fi
    done <<< "$all_funcs"

    [ "$removed" -eq 0 ] && return

    while IFS= read -r func; do
        [ -z "$func" ] && continue
        local escaped=$(echo "$func" | sed 's/\//\\\//g; s/\./\\./g')
        sed -i "/${escaped}\":/,/\/\/ end/d" "$wasm_exec_file"
    done <<< "$unused"

    local kept=$(($(echo "$all_funcs" | wc -l) - removed))
    echo "  Stripped $removed unused functions (kept $kept)"
}

while [[ $# -gt 0 ]]; do
    case $1 in
        --release) RELEASE_MODE=true; shift ;;
        -h|--help)
            echo "Usage: $0 [--release] [project-dir]"
            exit 0
            ;;
        -*) echo "Unknown option: $1"; exit 1 ;;
        *) PROJECT_DIR="$1"; shift ;;
    esac
done

echo "Cleaning dist folder..."
rm -rf "$PROJECT_DIR/dist"
mkdir -p "$PROJECT_DIR/dist"

echo "Generating HTML files and bindings..."
go run -tags '!wasm' "$PROJECT_DIR" 2>&1 | while read -r line; do
    if [[ "$line" == Generated:* ]]; then
        echo "  ${line#Generated: }"
    else
        echo "$line"
    fi
done

echo "Building WASM..."
TINYGO_FLAGS="-target wasm -scheduler=asyncify -gc=leaking"
if [ "$RELEASE_MODE" = true ]; then
    TINYGO_FLAGS="$TINYGO_FLAGS -panic=trap -no-debug"
fi
tinygo build -o "$PROJECT_DIR/dist/main.wasm" $TINYGO_FLAGS "$PROJECT_DIR"

echo "Copying wasm_exec.js..."
cp "$PROJECT_DIR/assets/wasm_exec.js" "$PROJECT_DIR/dist/"

if [ "$RELEASE_MODE" = true ]; then
    echo "String wasm"
    wasm-strip "$PROJECT_DIR/dist/main.wasm"
    echo "Strip wasm_exec"
    strip_wasm_exec "$PROJECT_DIR/dist/main.wasm" "$PROJECT_DIR/dist/wasm_exec.js"
    echo "Compressing..."
    for f in "$PROJECT_DIR/dist"/*.html "$PROJECT_DIR/dist"/*.wasm "$PROJECT_DIR/dist"/*.js "$PROJECT_DIR/dist"/*.json "$PROJECT_DIR/dist"/*.bin; do
        [ -f "$f" ] || continue
        (zopfli --i10 "$f" 2>/dev/null || gzip -k -f "$f" 2>/dev/null || true) &
        (brotli -q 11 -k -f "$f" 2>/dev/null || true) &
    done
    wait
    echo "Done"
fi

echo "Build complete! Output in: $PROJECT_DIR/dist/"
