.PHONY: run client

run:
	go run cmd/main.go --config=./config/config.yaml

client:
	go run examples/client.go --config=./config/config.yaml
