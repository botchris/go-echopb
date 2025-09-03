protos:
	rm -rf gen
	mkdir -p gen
	go tool buf generate --timeout 10m .

release:
	docker buildx build --no-cache --platform linux/amd64,linux/arm64,linux/386 -t botchrishub/echopb:latest --push .

client:
	go build -o dist/client ./cmd/echopb-client

server:
	go build -o dist/server ./cmd/echopb-server

install:
	go install ./cmd/echopb-client
	go install ./cmd/echopb-server
