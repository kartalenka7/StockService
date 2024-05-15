up:
	docker compose up stock-service -d
	$(MAKE) migrate-up

down:
	docker compose -f docker-compose.yml down -v


migrate-up:
	docker compose  --profile tools run --rm migrate up	

migrate-down:
	docker compose  --profile tools run --rm migrate down

build:
	docker compose up --build stock-service -d

test:
			go test ./...