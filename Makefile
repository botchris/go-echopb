protos:
	rm -rf gen
	mkdir -p gen
	go tool buf generate --timeout 10m .
