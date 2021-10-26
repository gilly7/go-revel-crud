TAG=v0.1.0
BINARY=go-revel-crud
NAME=ektowett/$(BINARY)
IMAGE=$(NAME):$(TAG)
LATEST=$(NAME):latest
LDFLAGS := -ldflags ""
DB_URL='postgres://go-revel-crud:go-revel-crud@127.0.0.1:5432/go-revel-crud?sslmode=disable'


run:
	@revel run

build:
	@docker-compose build

up:
	@docker-compose up -d

build_live:
	@echo "Building the image $(IMAGE)"
	@docker build -t $(IMAGE) . -f Dockerfile
	@echo "Tagging the image $(IMAGE) to latest"
	@docker tag $(IMAGE) $(LATEST)
	@echo "Remove the go binary"
	@rm $(BINARY)
	@echo "Done!"

logs:
	docker-compose logs -f

ps:
	@docker-compose ps

stop:
	@docker-compose stop

rm: stop
	@docker-compose rm

build_cli:
	@echo "Building cli to /tmp/go-revel-crud-cli"
	@go build -o /tmp/go-revel-crud-cli ./scripts/cli/main.go
	@echo "Done!"

# make migration name=create_users
migration:
	@echo "Creating migration $(name)!"
	@goose -dir migrations create $(name) sql
	@echo "Done!"

migrate_up:
	@echo "Migrating up!"
	@goose -dir migrations postgres $(DB_URL) up
	@echo "Done!"

migrate_down:
	@echo "Migrating down!"
	@goose -dir migrations postgres $(DB_URL) down
	@echo "Done!"

migrate_status:
	@echo "Getting migration status!"
	@goose -dir migrations postgres $(DB_URL) status
	@echo "Done!"

migrate_reset:
	@echo "Resetting migrations!"
	@goose -dir migrations postgres $(DB_URL) reset
	@echo "Done!"

migrate_version:
	@echo "Getting migration version!"
	@goose -dir migrations postgres $(DB_URL) version
	@echo "Done!"

migrate_redo: migrate_reset migrate_up
