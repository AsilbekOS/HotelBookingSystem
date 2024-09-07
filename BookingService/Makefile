generate-booking:
	protoc \
	--go_out=./ \
	--go_opt=paths=source_relative \
	--go-grpc_out=./ \
	--go-grpc_opt=paths=source_relative \
	proto/booking.proto

generate-hotelproto:
	protoc \
	--go_out=./ \
	--go_opt=paths=source_relative \
	--go-grpc_out=./ \
	--go-grpc_opt=paths=source_relative \
	hotelProto/hotel.proto

generate-notify:
	protoc \
	--go_out=./ \
	--go_opt=paths=source_relative \
	--go-grpc_out=./ \
	--go-grpc_opt=paths=source_relative \
	notifyProto/notify.proto


migrate-import:
	go get github.com/golang-migrate/migrate/v2/database/postgres

migrate-create:
	migrate create -ext sql -dir migrations -seq bookingservice

migrate-up:
	migrate -database postgres://postgres:1915@localhost:5432/bookingservice?sslmode=disable -path ./migrations up

migrate-down:
	migrate -database postgres://postgres:1915@localhost:5432/bookingservice?sslmode=disable -path ./migrations down

migrate-force:
	migrate -database postgres://postgres:1915@localhost:5432/bookingservice?sslmode=disable -path ./migrations force	

docker-build:
	sudo docker build -t booking-service-image .

docker-run:
	sudo docker run -p 4440:4440 booking-service-image
