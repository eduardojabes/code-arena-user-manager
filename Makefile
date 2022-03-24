.PHONY: migrate
migrate: 
	goose -dir deployment/migrations postgres "user=postgres password=postgres dbname=user-manager sslmode=disable" up