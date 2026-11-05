SWAG ?= $(shell go env GOPATH)/bin/swag

.PHONY: run swagger lint test fmt tools all

## tools: установить swag CLI
tools:
	go install github.com/swaggo/swag/cmd/swag@v1.16.4

## swagger: сгенерировать docs/ из аннотаций и экспортировать спецификацию в api/
swagger:
	$(SWAG) fmt -g cmd/api/main.go
	$(SWAG) init -g cmd/api/main.go -o docs --parseInternal
	cp docs/swagger.yaml api/swagger.yaml
	npx --yes swagger2openapi@7 api/swagger.yaml --yaml --outfile api/openapi.yaml

## lint: проверить спецификацию линтером Redocly
lint:
	npx --yes @redocly/cli@1 lint api/openapi.yaml

## run: запустить сервер (API_TOKEN=secret make run — включить Bearer-авторизацию)
run:
	go run ./cmd/api

test:
	go test ./... -count=1

fmt:
	gofmt -w .

## all: полная регенерация документации, проверка и тесты
all: swagger lint test
