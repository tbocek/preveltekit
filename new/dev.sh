#!/bin/bash
set -e

# =============================================================================
# PrevelteKit Development Server
# =============================================================================
# Watches for file changes and automatically rebuilds.
# Serves the application on localhost:8080.
#
# Usage: ./dev.sh <main-component.go> [child-component.go ...]
# =============================================================================

# -----------------------------------------------------------------------------
# Configuration
# -----------------------------------------------------------------------------

PORT="${PORT:-8080}"
DEBOUNCE_MS=100

# -----------------------------------------------------------------------------
# Usage
# -----------------------------------------------------------------------------

show_usage() {
    echo "Usage: $0 <main-component.go> [child-component.go ...]"
    echo ""
    echo "Starts a development server with automatic rebuilds on file changes."
    echo ""
    echo "Environment variables:"
    echo "  PORT    Server port (default: 8080)"
    echo ""
    echo "Examples:"
    echo "  $0 myapp/app.go"
    echo "  $0 myapp/app.go myapp/header.go myapp/footer.go"
    echo "  PORT=3000 $0 myapp/app.go"
}

# -----------------------------------------------------------------------------
# Argument Parsing
# -----------------------------------------------------------------------------

COMPONENT_FILES=()

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help) show_usage; exit 0 ;;
        -*)        echo "Unknown option: $1"; show_usage; exit 1 ;;
        *)         COMPONENT_FILES+=("$1"); shift ;;
    esac
done

if [ ${#COMPONENT_FILES[@]} -eq 0 ]; then
    echo "Error: No component files specified."
    echo ""
    show_usage
    exit 1
fi

# -----------------------------------------------------------------------------
# Resolve Paths
# -----------------------------------------------------------------------------

MAIN_COMPONENT="${COMPONENT_FILES[0]}"
PROJECT_DIR=$(dirname "$MAIN_COMPONENT")
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Compute checksums of source files for change detection
compute_checksum() {
    local files=("$@")
    cat "${files[@]}" 2>/dev/null | md5sum | cut -d' ' -f1
}

# Get all .go files in project directory
get_source_files() {
    find "$PROJECT_DIR" -maxdepth 1 -name "*.go" -type f 2>/dev/null | sort
}

# -----------------------------------------------------------------------------
# Build Function
# -----------------------------------------------------------------------------

do_build() {
    echo ""
    echo "[$(date +%H:%M:%S)] Building..."

    if "$SCRIPT_DIR/build.sh" --no-compress "${COMPONENT_FILES[@]}" 2>&1; then
        echo "[$(date +%H:%M:%S)] Build successful"
        return 0
    else
        echo "[$(date +%H:%M:%S)] Build failed"
        return 1
    fi
}

# -----------------------------------------------------------------------------
# Server Function
# -----------------------------------------------------------------------------

SERVER_PID=""

start_server() {
    # Kill existing server if running
    if [ -n "$SERVER_PID" ] && kill -0 "$SERVER_PID" 2>/dev/null; then
        kill "$SERVER_PID" 2>/dev/null || true
        wait "$SERVER_PID" 2>/dev/null || true
    fi

    # Start new server
    if command -v python3 &> /dev/null; then
        (cd "$PROJECT_DIR/dist" && python3 -m http.server "$PORT" --bind 127.0.0.1) &
        SERVER_PID=$!
    elif command -v python &> /dev/null; then
        (cd "$PROJECT_DIR/dist" && python -m SimpleHTTPServer "$PORT") &
        SERVER_PID=$!
    elif command -v caddy &> /dev/null; then
        caddy file-server --root "$PROJECT_DIR/dist" --listen ":$PORT" &
        SERVER_PID=$!
    else
        echo "Warning: No suitable server found (python3, python, or caddy)"
        echo "Serving files from: $PROJECT_DIR/dist"
        return 1
    fi

    echo "[$(date +%H:%M:%S)] Server running at http://localhost:$PORT"
}

# -----------------------------------------------------------------------------
# Cleanup
# -----------------------------------------------------------------------------

cleanup() {
    echo ""
    echo "Shutting down..."
    if [ -n "$SERVER_PID" ]; then
        kill "$SERVER_PID" 2>/dev/null || true
    fi
    exit 0
}

trap cleanup SIGINT SIGTERM

# -----------------------------------------------------------------------------
# Main Loop
# -----------------------------------------------------------------------------

echo "PrevelteKit Development Server"
echo "=============================="
echo ""
echo "Watching: $PROJECT_DIR/*.go"
echo "Output:   $PROJECT_DIR/dist/"
echo ""

# Initial build
LAST_CHECKSUM=""
SOURCE_FILES=$(get_source_files)

if do_build; then
    start_server
    LAST_CHECKSUM=$(compute_checksum $SOURCE_FILES)
else
    echo "Initial build failed. Waiting for changes..."
fi

echo ""
echo "Watching for changes... (Ctrl+C to stop)"
echo ""

# Watch loop
while true; do
    sleep 0.5

    SOURCE_FILES=$(get_source_files)
    CURRENT_CHECKSUM=$(compute_checksum $SOURCE_FILES)

    if [ "$CURRENT_CHECKSUM" != "$LAST_CHECKSUM" ]; then
        # Debounce: wait a bit for editor to finish writing
        sleep 0.1

        # Recompute in case more changes came in
        SOURCE_FILES=$(get_source_files)
        CURRENT_CHECKSUM=$(compute_checksum $SOURCE_FILES)

        if do_build; then
            # Server keeps running, just serves new files
            # For full reload, could restart server here
            :
        fi

        LAST_CHECKSUM="$CURRENT_CHECKSUM"
    fi
done
