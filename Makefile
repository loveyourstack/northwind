
# install CLI (must use this when testing so that latest config is always used)
.PHONY: cli
cli:
	sudo cp nw_config.toml /usr/local/etc
	go install ./cmd/north

# (re-)create database
.PHONY: db
db: cli
	north createDb

# regular server start
.PHONY: srv
srv:
	go run ./cmd/nwsrv

# run all tests
.PHONY: tests
tests: 
	go test -race ./...
