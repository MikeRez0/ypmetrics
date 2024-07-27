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

yptestdbstart:
	docker run --rm \
		--name=praktikum-db \
		--expose 5433 \
		-e POSTGRES_PASSWORD="postgres" \
		-e POSTGRES_USER="postgres" \
		-e POSTGRES_DB="praktikum" \
		-d \
		-p 5433:5432 \
		postgres:15.3

.PHONY: yptest
yptest: build-server build-agent
	./../metricstest-darwin-amd64 -test.v -test.run=$(test) \
	-binary-path cmd/server/server \
	-agent-binary-path cmd/agent/agent \
	-source-path=. \
	-server-port=8881 \
	-file-storage-path=./.tmp/metrics.json \
	-database-dsn=postgres://postgres:postgres@localhost:5433/praktikum?sslmode=disable

.PHONY: db-start
db-start:
	docker compose -f "scripts/db/docker-compose.yaml" up -d --build

.PHONY: db-stop
db-stop:
	docker compose -f "scripts/db/docker-compose.yaml" down

.PHONY: db-clean
db-clean:
	sudo rm -rf ./	scripts/db/data/

.PHONY: db-migration-new
db-migration-new:
	docker run --rm \
    -v $(realpath ./internal/storage/migrations):/migrations \
    migrate/migrate:v4.16.2 \
        create \
        -dir /migrations \
        -ext .sql \
        -seq -digits 5 \
        $(name)

.PHONY: db-migration-up
db-migration-up:
	docker run --rm \
    -v $(realpath ./internal/storage/migrations):/migrations \
	--network host \
    migrate/migrate:v4.16.2 \
        -path=/migrations \
        -database postgres://metrics:metrics@localhost:5432/metrics_db?sslmode=disable \
        up 

.PHONY: db-migration-down
db-migration-down:
	docker run --rm \
    -v $(realpath ./internal/storage/migrations):/migrations \
	--network host \
    migrate/migrate:v4.16.2 \
        -path=/migrations \
        -database postgres://metrics:metrics@localhost:5432/metrics_db?sslmode=disable \
        down 1