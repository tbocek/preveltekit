#!/bin/bash
set -e

# Initial build
echo "Initial build..."
./build.sh .

# Start livereload server
go run github.com/tbocek/preveltekit/cmd/livereload@latest &
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
    find . -name '*.go' -not -path '*/dist/*' | entr -d -s '
        echo "Rebuilding..."
        ./build.sh . && curl -s -X POST http://localhost:3001/trigger > /dev/null && echo "Reloaded."
    '
done
