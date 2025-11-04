include .env
export $(shell sed 's/=.*//' .env)

start:
	@go run src/main.go
lint:
	@golangci-lint run
tests:
	@go test -v ./test/...
tests-%:
	@go test -v ./test/... -run=$(shell echo $* | sed 's/_/./g')
testsum:
	@cd test && gotestsum --format testname
swagger:
	@cd src && ~/go/bin/swag  init
migration-%:
	@migrate create -ext sql -dir src/database/migrations create-table-$(subst :,_,$*)