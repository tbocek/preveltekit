# Tidy all Go modules (root, examples, and site)
tidy:
	go mod tidy
	cd examples && go mod tidy
	cd site && go mod tidy

# Run all tests in the project
test:
	go test ./...

# Build the site (SSR + WASM) and ensure modules are tidy
build: tidy
	bash site/build.sh site

# Build the examples (SSR + WASM) and ensure modules are tidy
build-examples: tidy
	bash examples/build.sh examples

# Comprehensive verification: tidy, vet, test, and build the site
verify: tidy
	go vet ./...
	go test ./...
	bash site/build.sh site
