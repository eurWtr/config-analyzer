.PHONY: build run test proto clean tidy

BINARY=analyzer.exe

build:
	go build -o $(BINARY) ./cmd/analyzer/

run: build
	$(CURDIR)\$(BINARY).exe $(ARGS)

test:
	go test ./...

proto:
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		api/proto/analyzer.proto

tidy:
	go mod tidy

clean:
	rm -f $(BINARY)

# Примеры использования:
# make run ARGS="config.yaml"
# make run ARGS="--stdin < config.json"
# make run ARGS="-s config.yaml"
# make run ARGS="--http :8080"
# make run ARGS="--grpc :9090"
# make run ARGS="-r ./configs/"