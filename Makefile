ARCH = $(shell uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')

.PHONY: run-back run-front mg sql-start sql-stop uml restart docs docker-start docker-stop docker-start-from-github docker-start-from-gitlab

run-back:
	@echo "REST API available at http://localhost:8080"
	cd back && go run main.go

run-front:
	@echo "Angular app available at http://localhost:4200"
	cd front && npm start

mg:
	cd back && go mod tidy

sql-start:
	docker run --name psql-container -e POSTGRES_PASSWORD=admin_library -e POSTGRES_USER=postgres -e POSTGRES_DB=library -p 5435:5432 -d postgres

sql-stop:
	docker stop psql-container && docker rm psql-container

uml:
	@echo "Generating diagrams in back/diagrams/"
	cd back && go run ./cmd/gen-uml
	@echo "Converting to PNG..."
	plantuml -tpng back/diagrams/*.puml
	@echo "PNG diagrams available in back/diagrams/"


restart:
	$(MAKE) sql-stop
	$(MAKE) sql-start
	@echo "Waiting for PostgreSQL..."
	@until docker exec psql-container pg_isready -U postgres -q 2>/dev/null; do sleep 1; done
	@echo "PostgreSQL ready."
	(cd back && go run main.go) & (cd front && npm start) & wait

docs:
	@echo "Local documentation at http://localhost:4173"
	@echo "Starting Zensical..."
	@cd docs && zensical serve -f ../zensical.toml --dev-addr localhost:4173

docker-start:
	@echo "Front end available at http://localhost"
	@echo "API available at http://localhost:8080"
	DOCKER_BUILDKIT=0 COMPOSE_DOCKER_CLI_BUILD=0 DOCKER_DEFAULT_PLATFORM=linux/$(ARCH) docker-compose up -d --build

docker-stop:
	docker-compose down

docker-start-from-github:
	@echo "Front end available at http://localhost"
	@echo "API available at http://localhost:8080"
	docker-compose pull back front
	docker-compose up -d

docker-start-from-gitlab:
	@echo "Front end available at http://localhost"
	@echo "API available at http://localhost:8080"
	BACK_IMAGE=registry.gitlab.com/lin56/library/back:latest \
	FRONT_IMAGE=registry.gitlab.com/lin56/library/front:latest \
	docker-compose pull back front
	BACK_IMAGE=registry.gitlab.com/lin56/library/back:latest \
	FRONT_IMAGE=registry.gitlab.com/lin56/library/front:latest \
	docker-compose up -d