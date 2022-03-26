.PHONY: migrate
migrate: 
	goose -dir deployment/migrations postgres "user=postgres password=postgres dbname=user-manager sslmode=disable" up

generate-proto:
	protoc --go_out=. --go-grpc_out=. ./proto/user/user.proto