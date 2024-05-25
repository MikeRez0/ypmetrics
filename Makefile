GOLANGCI_LINT_CACHE?=/tmp/praktikum-golangci-lint-cache

.PHONY: golangci-lint-run
golangci-lint-run: _golangci-lint-rm-unformatted-report

.PHONY: _golangci-lint-reports-mkdir
_golangci-lint-reports-mkdir:
	mkdir -p ./golangci-lint

.PHONY: _golangci-lint-run
_golangci-lint-run: _golangci-lint-reports-mkdir
	-docker run --rm \
    -v $(shell pwd):/app \
    -v $(GOLANGCI_LINT_CACHE):/root/.cache \
    -w /app \
    golangci/golangci-lint:v1.57.2 \
        golangci-lint run \
            -c .golangci.yml \
	> ./golangci-lint/report-unformatted.json

.PHONY: _golangci-lint-format-report
_golangci-lint-format-report: _golangci-lint-run
	cat ./golangci-lint/report-unformatted.json | jq > ./golangci-lint/report.json

.PHONY: _golangci-lint-rm-unformatted-report
_golangci-lint-rm-unformatted-report: _golangci-lint-format-report
	rm ./golangci-lint/report-unformatted.json

.PHONY: golangci-lint-clean
golangci-lint-clean:
	sudo rm -rf ./golangci-lint 


.PHONY: build-server
build-server:
	go build -o cmd/server/server cmd/server/*.go

.PHONY: build-agent
build-agent:
	go build -o cmd/agent/agent cmd/agent/*.go

.PHONY: test
test:
	go test ./...
.PHONY: test-cover
test-cover: test
	go test ./... -cover

.PHONY: yptest
yptest: build-server build-agent
	./../metricstest-darwin-amd64 -test.v -test.run=$(test) -binary-path cmd/server/server -agent-binary-path cmd/agent/agent -source-path=. -server-port=8888