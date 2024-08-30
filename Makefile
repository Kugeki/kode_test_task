.PHONY: run, stop, run-and-attach, run-tests

run:
	docker compose up -d --build

run-and-attach:
	docker compose up --build

stop:
	docker compose down

run-tests:
	docker compose -f compose.yml -f tests.compose.yml up --build -d
	docker logs -f "integration-tests"
	docker compose -f compose.yml -f tests.compose.yml down

install-swag:
	go install github.com/swaggo/swag/cmd/swag@latest

generate-swagger:
	swag fmt -d internal/ports/rest
	swag init -d internal/ports/rest -g rest.go