package hotel

import (
	bookingproto "apigateway/proto/booking"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func BookingApiConn(booking_port string) bookingproto.BookingServiceClient {
	conn, err := grpc.NewClient(booking_port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Failed to connect hotel_service")
	}

	user := bookingproto.NewBookingServiceClient(conn)

	return user
}
