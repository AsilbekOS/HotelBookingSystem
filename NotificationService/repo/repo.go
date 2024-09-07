package repo

// import (
// 	"log"
// 	userproto "notifyservice/protos/user"

// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// )

// func GetEmail() error {

// 	conn, err := grpc.NewClient("localhost:7701", grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		log.Printf("did not connect: %v", err)
// 		return err
// 	}
// 	defer conn.Close()

// 	userClient := userproto.NewUserServiceClient(conn)

// 	resp, err := userClient.RegisterUser()
// 	if err != nil {
// 		log.Printf("UpdateRoom failed: %v", err)
// 		return err
// 	}

// 	return nil
// }
