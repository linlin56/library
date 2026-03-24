.PHONY: run-back run-front mg sql-start sql-stop uml

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
