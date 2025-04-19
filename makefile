# Vars

CLIENT_CMD_PATH = ./cmd/client
SERVER_CMD_PATH = ./cmd/server
TARGET_CMD_PATH = ./cmd/target

SERVER_BIN_PATH = ./build/server/server

# Targets

.PHONY: test
test:
	go test -failfast -fullpath -cover -race -timeout=$(if $(timeout), $(count),5s) $(if $(count),-count=$(count)) $(if $(run),-run=$(run)) $(if $(package),$(package),./...)

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

mockgen\:install:
	go install github.com/golang/mock/mockgen@v1.6.0

generate:
	go generate ./...

exec\:client:
	docker-compose exec client bash

exec\:server:
	docker-compose exec client bash

exec\:target:
	docker-compose exec client bash