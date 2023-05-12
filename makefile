run:
	docker-compose up
stop:
	docker-compose down --rmi all

PROJECT_NAME := bookd
MODULE_NAME := github.com/roblesoft/bookd
DEPLOYMENT := $(CURDIR)/
DOCKER_COMPOSE := $(DEPLOYMENT)/
DOCKER_COMPOSE_CMD = docker-compose -p $(PROJECT_NAME)

.PHONY: local.build local.run

local.build: ## Build Local environment
	@echo "Build local environment..."
	$(DOCKER_COMPOSE_CMD) -f $(DOCKER_COMPOSE)/docker-compose.yml build

local.run: ## Run Local environment
	@echo "Run local environment..."
	$(DOCKER_COMPOSE_CMD) -f $(DOCKER_COMPOSE)/docker-compose.yml up

local.stop: ## Run Local environment
	@echo "Run local environment..."
	$(DOCKER_COMPOSE_CMD) -f $(DOCKER_COMPOSE)/docker-compose.yml down --rmi all