#!/bin/bash
set -e

cleanup() {
    trap - EXIT INT TERM
    kill 0 2>/dev/null
    wait 2>/dev/null
    exit 0
}
trap cleanup EXIT INT TERM

# Initial build
echo "Initial build..."
./build.sh .

# Start livereload server
go run github.com/tbocek/preveltekit/v2/cmd/livereload@latest

# Start Caddy with dev config
caddy run --config Caddyfile.dev --adapter caddyfile &

echo "Dev server running on http://localhost:8080"
echo "Watching for changes..."

# Watch .go files and rebuild on change
while true; do
    find . -name '*.go' -not -path '*/dist/*' | entr -dn -s '
        echo "Rebuilding..."
        ./build.sh . && curl -s -X POST http://localhost:3001/trigger > /dev/null && echo "Reloaded."
    '
done
