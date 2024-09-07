package hotel

import (
	hotelproto "apigateway/proto/hotel"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func HotelApiConn(hotel_port string) hotelproto.HotelServiceClient {
	conn, err := grpc.NewClient(hotel_port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Failed to connect hotel_service")
	}

	user := hotelproto.NewHotelServiceClient(conn)

	return user
}
