package service

import (
	nProto "bookingservice/notifyProto"
	"bookingservice/proto"
	"bookingservice/server/repo"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type BookingService struct {
	proto.UnimplementedBookingServiceServer
	DB *sql.DB
}

func NewBookingService(db *sql.DB) *BookingService {
	return &BookingService{DB: db}
}

func (b *BookingService) CreateBooking(ctx context.Context, req *proto.CreateBookingRequest) (*proto.BookingResponse, error) {
	totalAmount, err := repo.TotalAmountRoom(req.CheckInDate, req.CheckOutDate, req.RoomID, req.HotelID)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf(err.Error())
	}

	query := `
	INSERT INTO bookings (user_id, hotel_id, room_id, room_type, check_in_date, check_out_date, total_amount, status)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING booking_id`

	layout := "02.01.2006"
	startDateString := req.CheckInDate
	endDateString := req.CheckOutDate

	startDate, err := time.Parse(layout, startDateString)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf(err.Error())
	}

	endDate, err := time.Parse(layout, endDateString)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf(err.Error())
	}

	var bookingID int64
	err = b.DB.QueryRow(query, req.UserID, req.HotelID, req.RoomID, req.RoomType, startDate, endDate, totalAmount, "Tasdiqlangan").Scan(&bookingID)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf(err.Error())
	}

	err = repo.HotelClientCrB(req.RoomID)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf(err.Error())
	}

	message := nProto.SendNotificationRequest{
		UserID:      req.UserID,
		BookingID:   bookingID,
		Message:     "Booking successfully created",
		ServiceType: "Booking",
	}

	err = repo.SendNotify(ctx, &message)
	if err != nil {
		log.Println("1", err)
		return nil, fmt.Errorf(err.Error())
	}

	return &proto.BookingResponse{
		BookingID:    bookingID,
		UserID:       req.UserID,
		HotelID:      req.HotelID,
		RoomID:       req.RoomID,
		RoomType:     req.RoomType,
		CheckInDate:  startDateString,
		CheckOutDate: endDateString,
		TotalAmount:  totalAmount,
		Status:       "Tasdiqlangan",
	}, nil
}

func (b *BookingService) GetBooking(ctx context.Context, req *proto.GetBookingRequest) (*proto.BookingResponse, error) {
	query := `SELECT booking_id, user_id, hotel_id, room_id, room_type, check_in_date, check_out_date, total_amount, status FROM bookings WHERE booking_id = $1`

	var booking proto.BookingResponse
	err := b.DB.QueryRow(query, req.BookingID).Scan(
		&booking.BookingID,
		&booking.UserID,
		&booking.HotelID,
		&booking.RoomID,
		&booking.RoomType,
		&booking.CheckInDate,
		&booking.CheckOutDate,
		&booking.TotalAmount,
		&booking.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("bronlash topilmadi")
		}
		log.Println(err)
		return nil, fmt.Errorf(err.Error())
	}

	return &booking, nil
}

func (b *BookingService) UpdateBooking(ctx context.Context, req *proto.UpdateBookingRequest) (*proto.BookingResponse, error) {
	query := `
		UPDATE bookings
		SET check_in_date = $1, check_out_date = $2, total_amount = $3, status = $4
		WHERE booking_id = $5
		RETURNING booking_id, user_id, hotel_id, room_id, room_type, check_in_date, check_out_date, total_amount, status`

	layout := "02.01.2006"
	startDateString := req.CheckInDate
	endDateString := req.CheckOutDate

	startDate, err := time.Parse(layout, startDateString)
	if err != nil {
		log.Println("4", err)
		return nil, fmt.Errorf(err.Error())
	}

	endDate, err := time.Parse(layout, endDateString)
	if err != nil {
		log.Println("5", err)
		return nil, fmt.Errorf(err.Error())
	}

	var booking proto.BookingResponse
	err = b.DB.QueryRow(query, startDate, endDate, req.TotalAmount, req.Status, req.BookingID).Scan(
		&booking.BookingID,
		&booking.UserID,
		&booking.HotelID,
		&booking.RoomID,
		&booking.RoomType,
		&booking.CheckInDate,
		&booking.CheckOutDate,
		&booking.TotalAmount,
		&booking.Status,
	)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf(err.Error())
	}

	message := nProto.SendNotificationRequest{
		UserID:      booking.UserID,
		BookingID:   req.BookingID,
		Message:     "Booking successfully update",
		ServiceType: "Booking",
	}

	err = repo.SendNotify(ctx, &message)
	if err != nil {
		log.Println("1", err)
		return nil, fmt.Errorf(err.Error())
	}

	return &booking, nil
}

func (b *BookingService) CancelBooking(ctx context.Context, req *proto.CancelBookingRequest) (*proto.CancelBookingResponse, error) {
	err := repo.HotelClientCnB(req.RoomID)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf(err.Error())
	}
	query := `UPDATE bookings SET status = $1 WHERE booking_id = $2`
	_, err = b.DB.Exec(query, "Bekor qilingan", req.BookingID)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf(err.Error())
	}

	message := nProto.SendNotificationRequest{
		BookingID:   req.BookingID,
		Message:     "Booking successfully canceled",
		ServiceType: "Booking",
	}

	err = repo.SendNotify(ctx, &message)
	if err != nil {
		log.Println("1", err)
		return nil, fmt.Errorf(err.Error())
	}

	return &proto.CancelBookingResponse{
		Message:   "Bronlash muvaffaqiyatli bekor qilindi",
		BookingID: req.BookingID,
	}, nil
}

func (b *BookingService) ListUserBookings(ctx context.Context, req *proto.ListUserBookingsRequest) (*proto.ListUserBookingsResponse, error) {
	query := `SELECT booking_id, hotel_id, room_type, check_in_date, check_out_date, total_amount, status FROM bookings WHERE user_id = $1`

	rows, err := b.DB.QueryContext(ctx, query, req.UserID)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf(err.Error())
	}
	defer rows.Close()

	var bookings []*proto.BookingResponse
	for rows.Next() {
		var booking proto.BookingResponse
		err := rows.Scan(
			&booking.BookingID,
			&booking.HotelID,
			&booking.RoomType,
			&booking.CheckInDate,
			&booking.CheckOutDate,
			&booking.TotalAmount,
			&booking.Status,
		)
		if err != nil {
			log.Println(err)
			return nil, fmt.Errorf(err.Error())
		}
		bookings = append(bookings, &booking)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, fmt.Errorf(err.Error())
	}

	return &proto.ListUserBookingsResponse{Bookings: bookings}, nil
}
