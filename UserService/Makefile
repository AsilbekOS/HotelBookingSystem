migrate-import:
	go get github.com/golang-migrate/migrate/v2/database/postgres

migrate-create:
	migrate create -ext sql -dir migrations -seq userservice

migrate-up:
	migrate -database postgres://postgres:1915@localhost:5432/userservice?sslmode=disable -path ./migrations up

migrate-down:
	migrate -database postgres://postgres:1915@localhost:5432/userservice?sslmode=disable -path ./migrations down

migrate-force:
	migrate -database postgres://postgres:1915@localhost:5432/userservice?sslmode=disable -path ./migrations force

client-proto:
	protoc --go_out=. --go-grpc_out=. proto/user.proto

docker-build:
	sudo docker build -t user-service-image .

docker-run:
	sudo docker run -p 4440:4440 user-service-image	