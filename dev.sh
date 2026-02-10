#!/bin/bash
set -e

PROJECT_DIR="${1:-.}"

# Initial build
echo "Initial build..."
cd "$PROJECT_DIR"
./build.sh .
cd - > /dev/null

# Start livereload server
go run ./cmd/livereload &
LIVERELOAD_PID=$!

# Start Caddy with dev config
caddy run --config Caddyfile.dev --adapter caddyfile &
CADDY_PID=$!

cleanup() {
    kill $LIVERELOAD_PID $CADDY_PID 2>/dev/null
    wait $LIVERELOAD_PID $CADDY_PID 2>/dev/null
    exit 0
}
trap cleanup EXIT INT TERM

echo "Dev server running on http://localhost:8080"
echo "Watching for changes..."

# Watch .go files and rebuild on change
while true; do
    find "$PROJECT_DIR" -name '*.go' -not -path '*/dist/*' -not -path '*/cmd/*' | entr -d -s "
        echo 'Rebuilding...'
        cd \"$PROJECT_DIR\" && ./build.sh . && cd - > /dev/null
        curl -s -X POST http://localhost:3001/trigger > /dev/null && echo 'Reloaded.'
    "
done
