.PHONY: run client

run:
	go run ./cmd/sso --config=./config/config.yaml

client:
	go run examples/go/client.go --config=./config/config.yaml

db_up:
	go run ./cmd/migrator --storage-path=./storage/sso.db --migrations-path=./migrations
