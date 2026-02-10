.PHONY: run restart logs clean stop build

# Default target
all: run

# Build and start the containers in detached mode
run:
	docker compose up -d --build

# Restart the containers
restart:
	docker compose restart

SHELL := /bin/bash

# Show logs for a specific service or all
logs:
	@echo "Select a service to view logs:"
	@services=$$(docker compose ps --services); \
	select service in $$services "all" "exit"; do \
		case $$service in \
			"all") docker compose logs -f; break;; \
			"exit") exit 0;; \
			"") echo "Invalid selection";; \
			*) docker compose logs -f $$service; break;; \
		esac; \
	done

# Stop and remove containers, networks, and volumes
clean:
	docker compose down --volumes --remove-orphans

# Stop the containers
stop:
	docker compose stop

# Build the containers
build:
	docker compose build
