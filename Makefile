.SILENT:
APP=shortener

.PHONY: help
help: Makefile ## Show this help
	@echo
	@echo "Choose a command run in "$(APP)":"
	@echo
	@fgrep -h "##" $(MAKEFILE_LIST) | sed -e 's/\(\:.*\#\#\)/\:\ /' | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

.PHONY: build
build: ## Build an application
	@echo "Building ${APP} ..."
	mkdir -p build
	go build -o build/${APP} shortener/cmd/shortener

.PHONY: build-test
build-test: ## Build an application
	@echo "Building ${APP} ..."
	go build -o cmd/shortener/shortener shortener/cmd/shortener

test1: ## Test increment #1
	@echo "Testing ${APP} - increment 1..."
	tests/shortenertestbeta-darwin-arm64 -test.v -test.run=^TestIteration1$ -binary-path=cmd/shortener/shortener

run: ## Run an application
	@echo "Starting ${APP} ..."
	go run main.go

test: ## Run an application
	@echo "Testing ${APP} ..."
	go test

bench: ## Run an application
	@echo "Benchmarking ${APP} ..."
	go test -bench=. .

clean: ## Clean a garbage
	@echo "Cleaning"
	go clean
	rm -rf build

lint: ## Check a code by golangci-lint
	@echo "Linter checking..."
	golangci-lint run -c .golangci.yml ./...