#!/bin/bash
set -e

# =============================================================================
# PrevelteKit Build Script
# =============================================================================
# Builds a PrevelteKit Go/WASM application from source components.
#
# Build pipeline:
#   1. Code generation  - CLI extracts metadata via reflection, generates Go code
#   2. Pre-rendering    - For each route, renders static HTML (SSR)
#   3. WASM compilation - Compiles generated code to WebAssembly via TinyGo
#   4. Optimization     - Tree-shakes JS runtime, minifies, compresses
#   5. Assembly         - Combines HTML, WASM, and JS into final output files
# =============================================================================

# -----------------------------------------------------------------------------
# Configuration
# -----------------------------------------------------------------------------

SKIP_COMPRESS=false
RELEASE_MODE=false

# -----------------------------------------------------------------------------
# Usage
# -----------------------------------------------------------------------------

show_usage() {
    echo "Usage: $0 [OPTIONS] <main-component.go> [child-component.go ...]"
    echo ""
    echo "Builds a PrevelteKit application from Go component files."
    echo ""
    echo "Arguments:"
    echo "  <main-component.go>      The main/root component file"
    echo "  [child-component.go ...] Optional child components to include"
    echo ""
    echo "Options:"
    echo "  --release         Release build: silent panics, smaller output"
    echo "  --no-compress     Skip gzip/brotli compression of output files"
    echo "  -h, --help        Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 myapp/App.go"
    echo "  $0 myapp/App.go myapp/Header.go myapp/Footer.go"
}

# -----------------------------------------------------------------------------
# Argument Parsing
# -----------------------------------------------------------------------------

COMPONENT_FILES=()

while [[ $# -gt 0 ]]; do
    case $1 in
        --release)     RELEASE_MODE=true; shift ;;
        --no-compress) SKIP_COMPRESS=true; shift ;;
        -h|--help)     show_usage; exit 0 ;;
        -*)            echo "Unknown option: $1"; show_usage; exit 1 ;;
        *)             COMPONENT_FILES+=("$1"); shift ;;
    esac
done

if [ ${#COMPONENT_FILES[@]} -eq 0 ]; then
    echo "Error: No component files specified."
    echo ""
    show_usage
    exit 1
fi

# -----------------------------------------------------------------------------
# Helper Functions
# -----------------------------------------------------------------------------

# Tree-shakes wasm_exec.js by removing unused syscall/js functions.
wasm_shake() {
    local SOURCE="$1"
    local WASM="$2"
    local OUTPUT="$3"

    local USED
    USED=$(wasm-objdump -j Import -x "$WASM" 2>/dev/null \
        | grep 'syscall/js\.' \
        | sed 's/.*syscall\/js\.\([a-zA-Z]*\).*/\1/' \
        | sort -u)

    cp "$SOURCE" "$OUTPUT"

    local OPTIONAL_FUNCS=(
        valueDelete
        valueSetIndex
        valueInvoke
        valueNew
        valueInstanceOf
        copyBytesToGo
        copyBytesToJS
    )

    for func in "${OPTIONAL_FUNCS[@]}"; do
        if ! echo "$USED" | grep -q "^${func}$"; then
            sed -i "/\/\/ func ${func}/,/^[[:space:]]*},$/d" "$OUTPUT"
            sed -i "/\/\/ ${func}/,/^[[:space:]]*},$/d" "$OUTPUT"
        fi
    done
}

# -----------------------------------------------------------------------------
# Resolve Paths and Assets
# -----------------------------------------------------------------------------

MAIN_COMPONENT="${COMPONENT_FILES[0]}"
PROJECT_DIR=$(dirname "$MAIN_COMPONENT")

if [ -f "$PROJECT_DIR/assets/index.html" ]; then
    INDEX_HTML="$PROJECT_DIR/assets/index.html"
else
    INDEX_HTML="assets/index.html"
fi

if [ -f "$PROJECT_DIR/assets/wasm_exec.js" ]; then
    WASM_EXEC="$PROJECT_DIR/assets/wasm_exec.js"
else
    WASM_EXEC="assets/wasm_exec.js"
fi

# -----------------------------------------------------------------------------
# Generate Code
# -----------------------------------------------------------------------------
# The CLI tool uses reflection to extract component metadata and generates:
#   - build/       Generated Go code ready for compilation
#   - build/routes.txt  List of routes to pre-render
#   - dist/        Output directory for final assets

echo "Generating code from components..."
go run ./cmd/. "${COMPONENT_FILES[@]}"

# -----------------------------------------------------------------------------
# Pre-render (Server-Side Rendering)
# -----------------------------------------------------------------------------
# For each route in routes.txt, render the HTML with PRERENDER_PATH set

echo "Pre-rendering..."

# Read routes from routes.txt (format: path:filename)
ROUTES_FILE="$PROJECT_DIR/build/routes.txt"
if [ -f "$ROUTES_FILE" ]; then
    while IFS=: read -r path htmlfile; do
        [ -z "$path" ] && continue
        [[ "$path" == \#* ]] && continue

        outfile="${htmlfile%.html}_prerendered.html"
        echo "  $path -> $outfile"
        (
            cd "$PROJECT_DIR/build"
            PRERENDER_PATH="$path" go run -tags '!wasm' . > "../dist/$outfile"
        ) || { echo "Error pre-rendering $path"; exit 1; }
    done < "$ROUTES_FILE"
else
    # Default: single route
    echo "  / -> index_prerendered.html"
    (
        cd "$PROJECT_DIR/build"
        PRERENDER_PATH="/" go run -tags '!wasm' . > ../dist/index_prerendered.html
    )
fi

# -----------------------------------------------------------------------------
# Build WASM
# -----------------------------------------------------------------------------

echo "Building WASM..."
TINYGO_FLAGS="-target wasm -no-debug -scheduler=asyncify -gc=leaking"
if [ "$RELEASE_MODE" = true ]; then
    TINYGO_FLAGS="$TINYGO_FLAGS -panic=trap"
fi
(
    cd "$PROJECT_DIR/build"
    tinygo build -o ../dist/app.wasm $TINYGO_FLAGS .
)

wasm-strip "$PROJECT_DIR/dist/app.wasm"

# -----------------------------------------------------------------------------
# Optimize JavaScript Runtime
# -----------------------------------------------------------------------------

echo "Optimizing JS runtime..."
wasm_shake "$WASM_EXEC" "$PROJECT_DIR/dist/app.wasm" "$PROJECT_DIR/dist/wasm_exec.js"

if command -v minify &> /dev/null; then
    minify --type=js "$PROJECT_DIR/dist/wasm_exec.js" -o "$PROJECT_DIR/dist/wasm_exec.js"
fi

# -----------------------------------------------------------------------------
# Assemble Final Output
# -----------------------------------------------------------------------------
# Combines prerendered HTML, WASM, and JS into final HTML files

echo "Assembling final output..."
go run ./cmd/. --assemble "$PROJECT_DIR"

# -----------------------------------------------------------------------------
# Compress Output Files
# -----------------------------------------------------------------------------

if [ "$SKIP_COMPRESS" = false ]; then
    echo "Compressing output files..."
    for f in "$PROJECT_DIR/dist"/*.html "$PROJECT_DIR/dist"/*.wasm "$PROJECT_DIR/dist"/*.js; do
        [ -f "$f" ] || continue
        gzip -9 -k -f "$f"
        brotli -9 -k -f "$f"
    done
fi

# -----------------------------------------------------------------------------
# Done
# -----------------------------------------------------------------------------

echo ""
echo "Build complete! Output in: $PROJECT_DIR/dist/"
ls -lh "$PROJECT_DIR/dist/"/*.html 2>/dev/null | awk '{print "  " $NF ": " $5}'
