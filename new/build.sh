#!/bin/bash
set -e

# =============================================================================
# Reactive Framework Build Script
# =============================================================================
# Builds a Reactive Go/WASM application from source components.
#
# Why code generation?
# --------------------
# Component .go files use a declarative syntax (similar to Svelte) that isn't
# valid Go. The CLI transforms this into plain Go code with reactive bindings,
# DOM manipulation, and event wiring.
#
# Why a separate pre-render step?
# -------------------------------
# The CLI only does text transformation - it doesn't execute your Go logic.
# Pre-rendering runs your actual component code to produce the initial HTML.
# This must be separate because your components may have initialization logic,
# computed values, or other Go code that affects the rendered output.
#
# Build pipeline:
#   1. Code generation  - CLI transforms components into valid Go code
#   2. Pre-rendering    - Executes generated code to produce static HTML (SSR)
#   3. WASM compilation - Compiles generated code to WebAssembly via TinyGo
#   4. Optimization     - Tree-shakes JS runtime, minifies, compresses
#   5. Assembly         - Combines HTML, WASM, and JS into final output
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
    echo "Builds a Reactive application from Go component files."
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
#
# The Go WASM runtime (wasm_exec.js) includes many JS<->Go bridge functions,
# but most apps only use a subset. This function:
#   1. Inspects the compiled WASM to find which functions it actually imports
#   2. Removes unused function implementations from wasm_exec.js
#
# This can reduce the JS runtime size significantly.
#
# Arguments:
#   $1 - Source wasm_exec.js path
#   $2 - Compiled WASM file path (to inspect imports)
#   $3 - Output path for tree-shaken JS
wasm_shake() {
    local SOURCE="$1"
    local WASM="$2"
    local OUTPUT="$3"

    # Extract list of syscall/js functions the WASM actually imports
    local USED
    USED=$(wasm-objdump -j Import -x "$WASM" 2>/dev/null \
        | grep 'syscall/js\.' \
        | sed 's/.*syscall\/js\.\([a-zA-Z]*\).*/\1/' \
        | sort -u)

    cp "$SOURCE" "$OUTPUT"

    # These are the optional functions that can be removed if unused
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
            # Remove the function block (from comment marker to closing brace)
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

# Use project-local assets if they exist, otherwise fall back to framework defaults.
# This allows projects to customize index.html template or wasm_exec.js.
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
# The CLI tool parses .go component files and generates:
#   - build/       Generated Go code ready for compilation
#   - dist/        Output directory for final assets

echo "Generating code from components..."
go run ./cmd/. "${COMPONENT_FILES[@]}"

# -----------------------------------------------------------------------------
# Pre-render (Server-Side Rendering)
# -----------------------------------------------------------------------------
# Compiles and runs the app in a non-WASM environment to generate static HTML.

echo "Pre-rendering..."
(
    cd "$PROJECT_DIR/build"
    REACTIVE_BUILD=true go run -tags '!wasm' . > ../dist/prerendered.html
)

# -----------------------------------------------------------------------------
# Build WASM
# -----------------------------------------------------------------------------
# Uses TinyGo for smaller output size. Flags explained:
#   -target wasm         Target WebAssembly
#   -no-debug            Strip debug info for smaller size
#   -panic=trap          Use trap instruction for panics (release only, smaller)
#   -scheduler=asyncify  Required for async operations in browser
#   -gc=leaking          Use leaking GC (smallest, fine for short-lived apps)

echo "Building WASM..."
TINYGO_FLAGS="-target wasm -no-debug -scheduler=asyncify -gc=leaking"
if [ "$RELEASE_MODE" = true ]; then
    TINYGO_FLAGS="$TINYGO_FLAGS -panic=trap"
fi
(
    cd "$PROJECT_DIR/build"
    tinygo build -o ../dist/app.wasm $TINYGO_FLAGS .
)

# Strip additional symbols from WASM for smaller size
wasm-strip "$PROJECT_DIR/dist/app.wasm"

# -----------------------------------------------------------------------------
# Optimize JavaScript Runtime
# -----------------------------------------------------------------------------
# Tree-shake wasm_exec.js to remove unused Go<->JS bridge functions

echo "Optimizing JS runtime..."
wasm_shake "$WASM_EXEC" "$PROJECT_DIR/dist/app.wasm" "$PROJECT_DIR/dist/wasm_exec.js"

# Minify if the 'minify' tool is available (npm install -g minify)
if command -v minify &> /dev/null; then
    minify --type=js "$PROJECT_DIR/dist/wasm_exec.js" -o "$PROJECT_DIR/dist/wasm_exec.js"
fi

# -----------------------------------------------------------------------------
# Assemble Final Output
# -----------------------------------------------------------------------------
# Combines prerendered HTML, WASM, and JS into the final index.html

echo "Assembling final output..."
go run ./cmd/. --assemble "$PROJECT_DIR"

# -----------------------------------------------------------------------------
# Compress Output Files
# -----------------------------------------------------------------------------
# Creates .gz and .br versions for servers that support content negotiation

if [ "$SKIP_COMPRESS" = false ]; then
    echo "Compressing output files..."
    for f in "$PROJECT_DIR/dist"/*; do
        [ -f "$f" ] || continue
        [[ "$f" == *.gz || "$f" == *.br ]] && continue
        gzip -9 -k -f "$f"
        brotli -9 -k -f "$f"
    done
fi

# -----------------------------------------------------------------------------
# Done
# -----------------------------------------------------------------------------

echo ""
echo "Build complete! Output in: $PROJECT_DIR/dist/"
ls -lh "$PROJECT_DIR/dist/" | grep -v '^\(total\|d\)'
