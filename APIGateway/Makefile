docker-build:
	sudo docker build -t api-gateway-image .

docker-run:
	sudo docker run -p 4440:4440 api-gateway-image

generate-user:
	protoc \
	--go_out=./ \
	--go_opt=paths=source_relative \
	--go-grpc_out=./ \
	--go-grpc_opt=paths=source_relative \
	proto/user/user.proto

generate-booking:
	protoc \
	--go_out=./ \
	--go_opt=paths=source_relative \
	--go-grpc_out=./ \
	--go-grpc_opt=paths=source_relative \
	proto/booking/booking.proto

generate-hotel:
	protoc \
	--go_out=./ \
	--go_opt=paths=source_relative \
	--go-grpc_out=./ \
	--go-grpc_opt=paths=source_relative \
	proto/hotel/hotel.proto

swag-gen:
	~/go/bin/swag init -g ./api/cmd/main.go -o api/docs force 1	