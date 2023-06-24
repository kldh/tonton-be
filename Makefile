##@ Document
PROJECT_ROOT = $(shell pwd)
PROJECT_NAME ?= tonton-be
DOCKER_COMPOSE ?= cd devstack && docker-compose -p $(PROJECT_NAME)

build: ## Build local environment
	docker build -t $(PROJECT_NAME) .

up: ## Run local environment
	@$(DOCKER_COMPOSE) up -d

down: ## Shutdown local environment
	@$(DOCKER_COMPOSE) down