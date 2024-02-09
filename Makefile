db: ## start db
	docker-compose -f docker-compose.yml up --build --remove-orphans -d postgres
	sleep 10
	docker-compose up --build migrate

