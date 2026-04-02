.PHONY: run-back run-front mg sql-start sql-stop uml restart docs

run-back:
	@echo "API REST dispo sur http://localhost:8080"
	cd back && go run main.go

run-front:
	@echo "App Angular dispo sur http://localhost:4200"
	cd front && npm start

mg:
	cd back && go mod tidy

sql-start:
	docker run --name psql-container -e POSTGRES_PASSWORD=admin_bibli -e POSTGRES_USER=postgres -e POSTGRES_DB=bibliotheque -p 5435:5432 -d postgres

sql-stop:
	docker stop psql-container && docker rm psql-container

uml:
	@echo "Génération des diagrammes dans back/diagrams/"
	cd back && go run ./cmd/gen-uml
	@echo "Conversion en PNG..."
	plantuml -tpng back/diagrams/*.puml
	@echo "Diagrammes PNG disponibles dans back/diagrams/"


restart:
	$(MAKE) sql-stop
	$(MAKE) sql-start
	@echo "Attente de PostgreSQL..."
	@until docker exec psql-container pg_isready -U postgres -q 2>/dev/null; do sleep 1; done
	@echo "PostgreSQL pret."
	(cd back && go run main.go) & (cd front && npm start) & wait

docs:
	@echo "Documentation locale sur http://localhost:4173"
	@echo "Lancement Zensical..."
	@cd docs && zensical serve -f ../zensical.toml --dev-addr localhost:4173