.PHONY: help setup up down logs test clean

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## Setup and build all services
	docker compose build
	@echo "Services built successfully"

up: ## Start all services
	docker compose up -d
	@echo "Services starting..."
	@sleep 10
	@echo "\n=== Portfolio Platform Ready ==="
	@echo "  🎮 Portfolio & Game: http://localhost:8080"
	@echo "  🔧 Game API: http://localhost:8081/health"
	@echo "  📊 Grafana: http://localhost:3000 (admin/admin)"
	@echo "  📈 Prometheus: http://localhost:9090/targets"
	@echo "  📋 Metrics: http://localhost:8081/metrics"
	@echo "\nℹ️  Wait 30s for metrics to appear in Grafana"
	

down: ## Stop all services
	docker compose down

logs: ## Show logs from all services
	docker compose logs -f

test: ## Run tests inside containers
	docker compose run --rm game-service go test ./...
	docker compose run --rm frontend-service go test ./...

clean: ## Clean up all containers and volumes
	docker compose down -v
	docker system prune -f

terraform-init: ## Initialize Terraform
	cd terraform && terraform init

terraform-plan: ## Plan Terraform deployment
	cd terraform && terraform plan -var="environment=dev"

terraform-apply: ## Apply Terraform deployment
	cd terraform && terraform apply -var="environment=dev"

terraform-destroy: ## Destroy Terraform deployment
	cd terraform && terraform destroy -var="environment=dev"