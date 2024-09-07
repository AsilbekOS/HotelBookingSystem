package service

import (
	"context"
	"database/sql"
	"fmt"
	proto "hotelservice/proto"
	"log"
)

type HotelService struct {
	proto.UnimplementedHotelServiceServer
	DB *sql.DB
}

func NewHotelService(db *sql.DB) *HotelService {
	return &HotelService{DB: db}
}

// Hotels
func (h *HotelService) CreateHotel(ctx context.Context, req *proto.CreateHotelRequest) (*proto.CreateHotelResponse, error) {
	query := "INSERT INTO hotels (name, location, rating, address) VALUES ($1, $2, $3, $4) RETURNING hotel_id;"

	if req.Name == "" && req.Address == "" && req.Location == "" {
		log.Println("all field are empty")
		return nil, fmt.Errorf("all field are empty")
	}

	if req.Name == "" {
		log.Println("name field is empty")
		return nil, fmt.Errorf("name field is empty")
	}

	if req.Address == "" {
		log.Println("Address field is empty")
		return nil, fmt.Errorf("address field is empty")
	}

	if req.Location == "" {
		log.Println("Location field is empty")
		return nil, fmt.Errorf("location field is empty")
	}

	var hotelID int64
	err := h.DB.QueryRow(query, req.Name, req.Location, req.Rating, req.Address).Scan(&hotelID)
	if err != nil {
		log.Printf("failed to create hotel: %v", err)
		return nil, fmt.Errorf("failed to create hotel: %v", err)
	}

	res := &proto.CreateHotelResponse{
		Hotel: &proto.Hotel{
			HotelID:  hotelID,
			Name:     req.Name,
			Location: req.Location,
			Rating:   req.Rating,
			Address:  req.Address,
		},
	}

	log.Println("Successfully created hotel!")
	return res, nil
}

func (h *HotelService) GetHotels(ctx context.Context, req *proto.GetHotelsRequest) (*proto.GetHotelsResponse, error) {
	query := "SELECT hotel_id, name, location, rating, address FROM hotels"

	rows, err := h.DB.Query(query)
	if err != nil {
		log.Printf("failed to get hotels: %v", err)
		return nil, fmt.Errorf("failed to get hotels: %v", err)
	}
	defer rows.Close()

	var hotels []*proto.Hotel

	for rows.Next() {
		var hotel proto.Hotel
		if err := rows.Scan(&hotel.HotelID, &hotel.Name, &hotel.Location, &hotel.Rating, &hotel.Address); err != nil {
			log.Printf("failed to scan hotel: %v", err)
			return nil, fmt.Errorf("failed to scan hotel: %v", err)
		}
		hotels = append(hotels, &hotel)
	}

	if err := rows.Err(); err != nil {
		log.Printf("failed to process rows: %v", err)
		return nil, fmt.Errorf("failed to process rows: %v", err)
	}

	res := &proto.GetHotelsResponse{
		Hotels: hotels,
	}

	log.Println("Successfully retrieved hotels!")
	return res, nil
}

func (h *HotelService) GetHotelDetails(ctx context.Context, req *proto.GetHotelDetailsRequest) (*proto.GetHotelDetailsResponse, error) {
	hotelQuery := `
		SELECT hotel_id, name, location, rating, address
		FROM hotels
		WHERE hotel_id = $1
	`
	roomQuery := `
		SELECT room_type, price_per_night, availability
		FROM rooms
		WHERE hotel_id = $1
	`

	var hotel proto.Hotel
	err := h.DB.QueryRow(hotelQuery, req.HotelID).Scan(&hotel.HotelID, &hotel.Name, &hotel.Location, &hotel.Rating, &hotel.Address)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("hotel with ID %d not found", req.HotelID)
		}
		log.Printf("failed to get hotel details: %v", err)
		return nil, fmt.Errorf("failed to get hotel details: %v", err)
	}

	rows, err := h.DB.Query(roomQuery, req.HotelID)
	if err != nil {
		log.Printf("failed to get rooms: %v", err)
		return nil, fmt.Errorf("failed to get rooms: %v", err)
	}
	defer rows.Close()

	var rooms []*proto.Room
	for rows.Next() {
		var room proto.Room
		if err := rows.Scan(&room.RoomType, &room.PricePerNight, &room.Availability); err != nil {
			log.Printf("failed to scan room: %v", err)
			return nil, fmt.Errorf("failed to scan room: %v", err)
		}
		rooms = append(rooms, &room)
	}

	if err := rows.Err(); err != nil {
		log.Printf("failed to process rows: %v", err)
		return nil, fmt.Errorf("failed to process rows: %v", err)
	}

	res := &proto.GetHotelDetailsResponse{
		HotelID:  hotel.HotelID,
		Name:     hotel.Name,
		Location: hotel.Location,
		Rating:   hotel.Rating,
		Address:  hotel.Address,
		Rooms:    rooms,
	}

	log.Println("Successfully retrieved hotel details!")
	return res, nil
}

func (h *HotelService) CheckRoomAvailability(ctx context.Context, req *proto.CheckRoomAvailabilityRequest) (*proto.CheckRoomAvailabilityResponse, error) {
	query := `SELECT room_type, availability 
              FROM rooms 
              WHERE hotel_id = $1 AND availability = false`

	rows, err := h.DB.Query(query, req.HotelID)
	if err != nil {
		return nil, fmt.Errorf("failed to check room availability: %v", err)
	}
	defer rows.Close()

	var roomAvailabilities []*proto.RoomAvailability
	for rows.Next() {
		var roomType string
		var availability bool

		if err := rows.Scan(&roomType, &availability); err != nil {
			return nil, fmt.Errorf("failed to scan room availability: %v", err)
		}

		log.Printf("RoomType: %s, Availability: %v", roomType, availability)

		roomAvailabilities = append(roomAvailabilities, &proto.RoomAvailability{
			RoomType:       roomType,
			AvailableRooms: availability,
		})
	}

	return &proto.CheckRoomAvailabilityResponse{RoomAvailabilities: roomAvailabilities}, nil
}

func (h *HotelService) DeleteHotel(ctx context.Context, req *proto.DeleteHotelRequest) (*proto.DeleteHotelResponse, error) {
	query := "DELETE FROM hotels WHERE hotel_id = $1"
	_, err := h.DB.Exec(query, req.HotelID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete hotel: %v", err)
	}

	log.Println("Successfully delete hotel with ID:", req.HotelID)
	return &proto.DeleteHotelResponse{Success: true}, nil
}

// Rooms
func (h *HotelService) CreateRoom(ctx context.Context, req *proto.CreateRoomRequest) (*proto.CreateRoomResponse, error) {
	query := "INSERT INTO rooms (hotel_id, room_type, price_per_night, availability) VALUES ($1, $2, $3, $4) RETURNING room_id"
	var roomID int64
	err := h.DB.QueryRow(query, req.HotelID, req.RoomType, req.PricePerNight, req.Availability).Scan(&roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to create room: %v", err)
	}

	room := &proto.Rooms{
		RoomID:        roomID,
		HotelID:       req.HotelID,
		RoomType:      req.RoomType,
		PricePerNight: req.PricePerNight,
		Availability:  req.Availability,
	}

	return &proto.CreateRoomResponse{Room: room}, nil
}

func (h *HotelService) GetRoomByID(ctx context.Context, req *proto.GetRoomDetailsRequest) (*proto.GetRoomDetailsResponse, error) {
	query := "SELECT room_id, hotel_id, room_type, price_per_night, availability FROM rooms WHERE hotel_id = $1 AND room_id = $2"

	var room proto.Rooms

	err := h.DB.QueryRow(query, req.HotelID, req.RoomID).Scan(
		&room.RoomID,
		&room.HotelID,
		&room.RoomType,
		&room.PricePerNight,
		&room.Availability,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("room with ID %d not found", req.RoomID)
		}
		log.Printf("failed to get room by ID: %v", err)
		return nil, fmt.Errorf("failed to get room: %v", err)
	}

	response := &proto.GetRoomDetailsResponse{
		RoomID:        room.RoomID,
		HotelID:       room.HotelID,
		RoomType:      room.RoomType,
		PricePerNight: room.PricePerNight,
		Availability:  room.Availability,
	}

	log.Printf("Successfully retrieved room with ID: %d", room.RoomID)
	return response, nil
}

func (h *HotelService) GetRooms(ctx context.Context, req *proto.GetRoomsRequest) (*proto.GetRoomsResponse, error) {
	query := "SELECT room_id, hotel_id, room_type, price_per_night, availability FROM rooms WHERE hotel_id = $1"
	rows, err := h.DB.Query(query, req.HotelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get rooms: %v", err)
	}
	defer rows.Close()

	var rooms []*proto.Rooms
	for rows.Next() {
		var room proto.Rooms
		if err := rows.Scan(&room.RoomID, &room.HotelID, &room.RoomType, &room.PricePerNight, &room.Availability); err != nil {
			return nil, fmt.Errorf("failed to scan room: %v", err)
		}
		rooms = append(rooms, &room)
	}

	return &proto.GetRoomsResponse{Rooms: rooms}, nil
}

func (h *HotelService) UpdateRoom(ctx context.Context, req *proto.UpdateRoomRequest) (*proto.UpdateRoomResponse, error) {
	var query string
	var params []interface{}
	var paramCount int

	query = "UPDATE rooms SET "

	if req.RoomType != "" {
		paramCount++
		query += fmt.Sprintf("room_type = $%d, ", paramCount)
		params = append(params, req.RoomType)
	}

	if req.PricePerNight != 0 {
		paramCount++
		query += fmt.Sprintf("price_per_night = $%d, ", paramCount)
		params = append(params, req.PricePerNight)
	}

	if req.Availability || !req.Availability {
		paramCount++
		query += fmt.Sprintf("availability = $%d, ", paramCount)
		params = append(params, req.Availability)
	}

	if paramCount == 0 {
		log.Println("no fields to update")
		return nil, fmt.Errorf("no fields to update")
	}

	query = query[:len(query)-2]
	paramCount++
	query += fmt.Sprintf(" WHERE room_id = $%d", paramCount)
	params = append(params, req.RoomID)

	_, err := h.DB.Exec(query, params...)
	if err != nil {
		log.Println("failed to update room")
		return nil, fmt.Errorf("failed to update room: %v", err)
	}

	room := &proto.Rooms{
		RoomID:        req.RoomID,
		RoomType:      req.RoomType,
		PricePerNight: req.PricePerNight,
		Availability:  req.Availability,
	}

	return &proto.UpdateRoomResponse{Room: room}, nil
}

func (h *HotelService) DeleteRoom(ctx context.Context, req *proto.DeleteRoomRequest) (*proto.DeleteRoomResponse, error) {
	query := "DELETE FROM rooms WHERE room_id = $1"
	_, err := h.DB.Exec(query, req.RoomID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete room: %v", err)
	}

	return &proto.DeleteRoomResponse{Success: true}, nil
}
