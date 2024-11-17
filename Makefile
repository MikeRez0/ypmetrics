GOLANGCI_LINT_CACHE?=praktikum-golangci-lint-cache
current_dir := $(dir $(abspath $(firstword $(MAKEFILE_LIST))))

.PHONY: golangci-lint-run
golangci-lint-run:
	-docker run --rm -v .:/source -v $(GOLANGCI_LINT_CACHE):/root/.cache -w //source golangci/golangci-lint golangci-lint run -c .golangci.yml
	-bash -c 'cat ./.golangci-lint/report-unformatted.json | jq > ./.golangci-lint/report.json'


.PHONY: build-server
build-server:
	go build -C cmd/server

.PHONY: build-agent
build-agent:
	go build -C cmd/agent

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
	./.yptests/metricstest -test.v -test.run=$(test) \
	-binary-path cmd/server/server \
	-agent-binary-path cmd/agent/agent \
	-source-path=. \
	-server-port=8881 \
	-file-storage-path=./.tmp/metrics.json \
	-key="test" \
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