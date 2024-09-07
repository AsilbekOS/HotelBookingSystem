package repo

import (
	hotelproto "bookingservice/hotelProto"
	nProto "bookingservice/notifyProto"
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func HotelClientCrB(roomID int64) error {
	conn, err := grpc.NewClient("localhost:7702", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect: %v", err)
		return err
	}
	defer conn.Close()

	hotelClient := hotelproto.NewHotelServiceClient(conn)

	req := &hotelproto.UpdateRoomRequest{
		RoomID:       roomID,
		Availability: true,
	}

	resp, err := hotelClient.UpdateRoom(context.Background(), req)
	if err != nil {
		log.Printf("UpdateRoom failed: %v", err)
		return err
	}

	log.Printf("Room updated: %+v\n", resp.Room)

	return nil
}

func HotelClientCnB(roomID int64) error {
	conn, err := grpc.NewClient("localhost:7702", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect: %v", err)
		return err
	}
	defer conn.Close()

	hotelClient := hotelproto.NewHotelServiceClient(conn)

	req := &hotelproto.UpdateRoomRequest{
		RoomID:       roomID,
		Availability: false,
	}

	resp, err := hotelClient.UpdateRoom(context.Background(), req)
	if err != nil {
		log.Printf("UpdateRoom failed: %v", err)
		return err
	}

	log.Printf("Room updated: %+v\n", resp.Room)

	return nil
}

func TotalAmountRoom(StartDate, EndDate string, roomID, hotelID int64) (float64, error) {
	conn, err := grpc.NewClient("localhost:7702", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect: %v", err)
		return 0, err
	}
	defer conn.Close()

	roomClient := hotelproto.NewHotelServiceClient(conn)

	req := &hotelproto.GetRoomDetailsRequest{
		RoomID:  roomID,
		HotelID: hotelID,
	}

	resp, err := roomClient.GetRoomByID(context.Background(), req)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	roomPrice := resp.PricePerNight
	log.Println(roomPrice)

	layout := "02.01.2006"
	startDateString := StartDate
	endDateString := EndDate

	startDate, err := time.Parse(layout, startDateString)
	log.Println(startDate)
	if err != nil {
		fmt.Println("Xato:", err)
		return 0, err
	}

	endDate, err := time.Parse(layout, endDateString)
	log.Println(endDate)
	if err != nil {
		fmt.Println("Xato:", err)
		return 0, err
	}

	duration := endDate.Sub(startDate)
	log.Println(duration)

	days := int(duration.Hours() / 24)
	log.Println(days)

	totalAmount := float64(days) * roomPrice
	log.Println(totalAmount)

	return totalAmount, nil
}

func SendNotify(ctx context.Context, req *nProto.SendNotificationRequest) error {
	conn, err := grpc.NewClient("localhost:7704", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect: %v", err)
		return err
	}
	defer conn.Close()

	NotificationClient := nProto.NewNotificationServiceClient(conn)

	notificationReq := &nProto.SendNotificationRequest{
		UserID:      req.UserID,
		BookingID:   req.BookingID,
		Message:     req.Message,
		ServiceType: req.ServiceType,
	}
	_, err = NotificationClient.SendNotification(ctx, notificationReq)
	if err != nil {
		log.Printf("Failed to send notification: %v", err)
		return err
	}

	return nil
}
