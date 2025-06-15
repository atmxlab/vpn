# Vars

CLIENT_CMD_PATH = ./cmd/client
SERVER_CMD_PATH = ./cmd/server
TARGET_CMD_PATH = ./cmd/target

SERVER_BIN_PATH = ./build/server/server

TEST_ARGS = -vet=all -failfast -fullpath -cover -race -timeout=$(if $(timeout), $(count),5s) $(if $(count),-count=$(count)) $(if $(run),-run=$(run)) $(if $(package),$(package),./...)

# Targets

run\:server:
	go run $(SERVER_CMD_PATH)

run\:client:
	go run $(CLIENT_CMD_PATH)

run\:target:
	go run $(TARGET_CMD_PATH)


build\:server:
	make build output=$(SERVER_BIN_PATH) main=$(SERVER_CMD_PATH)

build:
	go build -buildvcs=false -o ${output} ${main}

up\:server:
	$(SERVER_BIN_PATH)

test-env-up:
	docker-compose up -d

test-env-down:
	docker-compose down

test-env-refresh: test-env-down test-env-up

exec\:client:
	docker-compose exec client bash

exec\:server:
	docker-compose exec server bash

exec\:target:
	docker-compose exec target bash


mockgen\:install:
	go install github.com/golang/mock/mockgen@v1.6.0

generate:
	go generate ./...

.PHONY: test
test:
	go test $(TEST_ARGS)

.PHONY: testsum
testsum:
	gotestsum -- $(TEST_ARGS)

deadcode\:install:
	go install golang.org/x/tools/cmd/deadcode@latest

.PHONY: deadcode
deadcode:
	deadcode ./cmd/... ./internal/... ./pkg/...

clean\:testcache:
	go clean -testcache