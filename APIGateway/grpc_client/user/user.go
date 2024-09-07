package user

import (
	userproto "apigateway/proto/user"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func UserApiConn(user_port string) userproto.UserServiceClient {
	conn, err := grpc.NewClient(user_port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Failed to connect user_service")
	}

	user := userproto.NewUserServiceClient(conn)

	return user
}
