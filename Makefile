.PHONY: run, stop, run-and-attach

run:
	docker compose up -d --build

run-and-attach:
	docker compose up --build

stop:
	docker compose down