# Vars

CLIENT_CMD_PATH = ./cmd/server/main.go
SERVER_CMD_PATH = ./cmd/server/main.go
TARGET_CMD_PATH = ./cmd/server/main.go

# Targets

test:
	go test -failfast -fullpath -cover -race -timeout=$(if $(timeout), $(count), 5s) $(if $(count),-count=$(count)) $(if $(run),-run=$(run)) $(if $(package),$(package),./...)

run\:server:
	go run $(SERVER_CMD_PATH)

run\:client:
	go run $(CLIENT_CMD_PATH)

run\:target:
	go run $(TARGET_CMD_PATH)


SERVER_BIN_PATH = ./build/server/server

build\:server:
	go build  -o $(SERVER_BIN_PATH) $(SERVER_CMD_PATH)

up\:server\:bin:
	$(SERVER_BIN_PATH)

up:
	docker-compose up -d

down:
	docker-compose down

mockgen\:install:
	go install github.com/golang/mock/mockgen@v1.6.0

mockgen:
	make mockgen:install

	mockgen -source=internal/http/handler/auth.go -destination=internal/http/handler/mocks/auth.go -package=mockhandler
	mockgen -source=internal/http/http.go -destination=internal/http/mocks/http.go -package=mockhttp
