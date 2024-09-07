generate-notify:
	protoc \
	--go_out=./ \
	--go_opt=paths=source_relative \
	--go-grpc_out=./ \
	--go-grpc_opt=paths=source_relative \
	protos/notify/notify.proto

generate-user:
	protoc \
	--go_out=./ \
	--go_opt=paths=source_relative \
	--go-grpc_out=./ \
	--go-grpc_opt=paths=source_relative \
	protos/user/user.proto

docker-build:
	sudo docker build -t notification-service-image .

docker-run:
	sudo docker run -p 4440:4440 notification-service-image