package main

import (
	"log"
	"net"
	proto "notifyservice/protos/notify"
	"notifyservice/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := ":7704"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Println(err)
	}
	defer lis.Close()

	server := service.NewNotificationService()
	s := grpc.NewServer()
	proto.RegisterNotificationServiceServer(s, server)
	reflection.Register(s)

	log.Println("Server is listening on port:", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error to Server - %v", err)
	}
}
