#!/bin/bash
set -e

# =============================================================================
# PrevelteKit Declarative Build Script
# =============================================================================
# Builds a PrevelteKit application using the p.Hydrate() API.
#
# Build pipeline:
#   1. SSR Phase    - Generates HTML files and collects bindings
#   2. WASM Build   - Compiles Go code to WebAssembly via TinyGo
#   3. Optimization - Compresses output files (optional)
# =============================================================================

SKIP_COMPRESS=false
RELEASE_MODE=false
OUTPUT_DIR=""
PROJECT_DIR=""

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --release)     RELEASE_MODE=true; shift ;;
        --no-compress) SKIP_COMPRESS=true; shift ;;
        -o|--output)   OUTPUT_DIR="$2"; shift 2 ;;
        -h|--help)
            echo "Usage: $0 [OPTIONS] [project-dir]"
            echo ""
            echo "Arguments:"
            echo "  project-dir    Directory containing main.go (default: current dir)"
            echo ""
            echo "Options:"
            echo "  --release      Release build: smaller WASM, stripped debug info"
            echo "  --no-compress  Skip gzip/brotli compression"
            echo "  -o, --output   Output directory (default: project-dir/dist)"
            echo "  -h, --help     Show this help"
            exit 0
            ;;
        -*)
            echo "Unknown option: $1"; exit 1 ;;
        *)
            # First non-option argument is the project directory
            if [ -z "$PROJECT_DIR" ]; then
                PROJECT_DIR="$1"
            fi
            shift
            ;;
    esac
done

# Default to current directory if no project specified
if [ -z "$PROJECT_DIR" ]; then
    PROJECT_DIR="."
fi

# Handle main.go path - extract directory
if [[ "$PROJECT_DIR" == *.go ]]; then
    PROJECT_DIR=$(dirname "$PROJECT_DIR")
fi

# Default output dir
if [ -z "$OUTPUT_DIR" ]; then
    OUTPUT_DIR="$PROJECT_DIR/dist"
fi

# Change to project directory
cd "$PROJECT_DIR"

# -----------------------------------------------------------------------------
# Step 1: SSR - Generate HTML files and collect bindings
# -----------------------------------------------------------------------------
echo "Generating HTML files..."

mkdir -p "dist"

# Run SSR phase - generates HTML files and outputs bindings to stderr
HYDRATE_MODE=generate-all go run -tags '!wasm' . 2>&1 | while read -r line; do
    if [[ "$line" == BINDINGS:* ]]; then
        # Extract bindings JSON
        echo "${line#BINDINGS:}" > "dist/bindings.json"
        echo "  Saved bindings.json"
    elif [[ "$line" == Generated:* ]]; then
        echo "  ${line#Generated: }"
    else
        echo "$line"
    fi
done

# -----------------------------------------------------------------------------
# Step 2: Build WASM with TinyGo
# -----------------------------------------------------------------------------
echo "Building WASM..."

TINYGO_FLAGS="-target wasm -scheduler=asyncify -gc=leaking"
if [ "$RELEASE_MODE" = true ]; then
    TINYGO_FLAGS="$TINYGO_FLAGS -panic=trap -no-debug"
fi

tinygo build -o "dist/main.wasm" $TINYGO_FLAGS .

if [ "$RELEASE_MODE" = true ] && command -v wasm-strip &> /dev/null; then
    wasm-strip "dist/main.wasm"
fi

# -----------------------------------------------------------------------------
# Step 3: Copy wasm_exec.js
# -----------------------------------------------------------------------------
echo "Copying wasm_exec.js..."

WASM_EXEC=$(tinygo env TINYGOROOT)/targets/wasm_exec.js
if [ -f "$WASM_EXEC" ]; then
    cp "$WASM_EXEC" "dist/"
else
    echo "Warning: Could not find TinyGo wasm_exec.js at $WASM_EXEC"
fi

# -----------------------------------------------------------------------------
# Step 4: Compress (optional)
# -----------------------------------------------------------------------------
if [ "$SKIP_COMPRESS" = false ]; then
    echo "Compressing output files..."
    for f in dist/*.html dist/*.wasm dist/*.js dist/*.json; do
        [ -f "$f" ] || continue
        if command -v zopfli &> /dev/null; then
            zopfli --i10 "$f" 2>/dev/null || gzip -k -f "$f" 2>/dev/null || true
        else
            gzip -k -f "$f" 2>/dev/null || true
        fi
        if command -v brotli &> /dev/null; then
            brotli -q 11 -k -f "$f" 2>/dev/null || true
        fi
    done
fi

# -----------------------------------------------------------------------------
# Done
# -----------------------------------------------------------------------------
echo ""
echo "Build complete! Output in: $(pwd)/dist/"
echo ""

# Show file sizes
WASM_SIZE=$(ls -lh "dist/main.wasm" 2>/dev/null | awk '{print $5}')
WASM_GZ_SIZE=$(ls -lh "dist/main.wasm.gz" 2>/dev/null | awk '{print $5}')
WASM_BR_SIZE=$(ls -lh "dist/main.wasm.br" 2>/dev/null | awk '{print $5}')

echo "WASM size: $WASM_SIZE"
[ -n "$WASM_GZ_SIZE" ] && echo "  gzip:   $WASM_GZ_SIZE"
[ -n "$WASM_BR_SIZE" ] && echo "  brotli: $WASM_BR_SIZE"
echo ""

ls -lh dist/*.html 2>/dev/null | awk '{printf "  %-30s %s\n", $NF, $5}'
