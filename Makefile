.PHONY: run mg sql-start sql-stop
run: 
	echo "Le site est servi ici : http://localhost:8080"
	go run main.go

mg: 
	go mod tidy

sql-start:
	docker run --name psql-container -e POSTGRES_PASSWORD=admin_bibli -e POSTGRES_USER=postgres -e POSTGRES_DB=bibliotheque -p 5435:5432 -d postgres

sql-stop:
	docker stop psql-container && docker rm psql-container