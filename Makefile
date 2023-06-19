##@ Document
PROJECT_ROOT = $(shell pwd)

diagram: ## Run diagram
	docker run -it -d --rm --name structurizr -p 3030:8080 -v $(PROJECT_ROOT)/docs/diagram:/usr/local/structurizr structurizr/lite
	sleep 3

	open http://localhost:3030