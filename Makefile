#.SILENT:
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
	go mod tidy
	go generate ./...
	go build -o build/${APP} shortener/cmd/shortener

test-static: ## Test static
	@echo "Testing ${APP} - static..."
	go vet -vettool="$(shell which ./tests/statictest-darwin-arm64)" ./...

.PHONY: test_all lint tests build-test test1 test2 test3 test4 test5 test6 test7
build-test: ## Build an application
	@echo "Building ${APP} ..."
	go mod tidy
	go generate ./...
	cd cmd/shortener && go build -buildvcs=false -o shortener

lint: ## Check a code by golangci-lint
	@echo "Linter checking..."
	golangci-lint run --fix -c .golangci.yml ./...

tests: ## Internal tests
	@echo "Testing ${APP} ..."
	go test ./...

test1: ## Test increment #1
	@echo "Testing ${APP} - increment 1..."
	tests/shortenertestbeta-darwin-arm64 -test.v -test.run="^TestIteration1$$" -binary-path=cmd/shortener/shortener

test2: ## Test increment #2
	@echo "Testing ${APP} - increment 2..."
	tests/shortenertestbeta-darwin-arm64 -test.v -test.run="^TestIteration2$$" -source-path=.

test3: ## Test increment #3
	@echo "Testing ${APP} - increment 3..."
	tests/shortenertestbeta-darwin-arm64 -test.v -test.run="^TestIteration3$$" -source-path=.

test4: ## Test increment #4
	@echo "Testing ${APP} - increment 4..."
	tests/shortenertestbeta-darwin-arm64 -test.v -test.run="^TestIteration4$$" -binary-path=cmd/shortener/shortener -server-port=8001

test5: ## Test increment #5
	@echo "Testing ${APP} - increment 5..."
	tests/shortenertestbeta-darwin-arm64 -test.v -test.run="^TestIteration5$$" -binary-path=cmd/shortener/shortener -server-port=8002

test6: ## Test increment #6
	@echo "Testing ${APP} - increment 6..."
	tests/shortenertestbeta-darwin-arm64 -test.v -test.run="^TestIteration6$$" -source-path=.

test7: ## Test increment #7
	@echo "Testing ${APP} - increment 7..."
	tests/shortenertestbeta-darwin-arm64 -test.v -test.run="^TestIteration7$$" -binary-path=cmd/shortener/shortener -source-path=.

test_all: lint tests build-test test1 test2 test3 test4 test5 test6 test7
	@echo "All tests completed."

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
