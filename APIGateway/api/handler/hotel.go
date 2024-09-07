package handler

import (
	hotelproto "apigateway/proto/hotel"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"google.golang.org/protobuf/encoding/protojson"
)

type HotelClient struct {
	Client hotelproto.HotelServiceClient
}

func NewHotelClient(Cl hotelproto.HotelServiceClient) *HotelClient {
	return &HotelClient{Client: Cl}
}

// @Router		/api/hotels [post]
// @Summary		Create hotel
// @Description Bu endpoint yangi Mehmonxonani ro'yxatdan o'tkazish uchun ishlatiladi uchun ishlatiladi
// @Security	BearerAuth
// @Tags		Hotel
// @Accept		json
// @Produce 	json
// @Param		body body hotelproto.CreateHotelRequest true "CreateHotelRequest"
// @Success 	201 {object} map[string]interface{} "Post Created"
// @Failure 	400 {object} models.StandartError "Bad request error"
// @Failure 	403 {object} models.ForbiddenError "Forbidden error"
// @Failure 	500 {object} models.StandartError "Internal server error"
func (h *HotelClient) CreateHotel(w http.ResponseWriter, r *http.Request) {
	var hotel hotelproto.CreateHotelRequest

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("CreateHotel-io.ReadAll:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := protojson.Unmarshal(bytes, &hotel); err != nil {
		log.Println("CreateHotel-protojson.Unmarshal:", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	resp, err := h.Client.CreateHotel(r.Context(), &hotel)
	if err != nil {
		log.Println("u.Client.RegisterUser:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("ERROR: Ma'lumotni encode qilishda xatolik...")
		http.Error(w, "Ma'lumotni encode qilishda xatolik...", http.StatusInternalServerError)
		return
	}
}

// @Router      /api/hotels [get]
// @Summary     Get list of hotels
// @Description Bu endpoint mavjud mehmonxonalar ro'yxatini olish uchun ishlatiladi
// @Security    BearerAuth
// @Tags        Hotel
// @Accept      json
// @Produce     json
// @Success 	201 {object} map[string]interface{} "Get GetHotels"
// @Failure 	400 {object} models.StandartError "Bad request error"
// @Failure 	403 {object} models.ForbiddenError "Forbidden error"
// @Failure 	500 {object} models.StandartError "Internal server error"
func (h *HotelClient) GetHotels(w http.ResponseWriter, r *http.Request) {
	var request hotelproto.GetHotelsRequest

	resp, err := h.Client.GetHotels(r.Context(), &request)
	if err != nil {
		log.Println("h.Client.GetHotels:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("ERROR: Ma'lumotni encode qilishda xatolik...")
		http.Error(w, "Ma'lumotni encode qilishda xatolik...", http.StatusInternalServerError)
		return
	}
}

// @Router      /api/hotelsid [get]
// @Summary     Get hotel details
// @Description Bu endpoint ma'lum bir mehmonxona haqidagi batafsil ma'lumotlarni olish uchun ishlatiladi
// @Security    BearerAuth
// @Tags        Hotel
// @Accept      json
// @Produce     json
// @Param		hotel_id query string true "Hotel ID: "
// @Success     200 {object} hotelproto.GetHotelDetailsResponse "Hotel details"
// @Failure     400 {object} models.StandartError "Bad request error"
// @Failure 	403 {object} models.ForbiddenError "Forbidden error"
// @Failure     500 {object} models.StandartError "Internal server error"
func (h *HotelClient) GetHotelDetails(w http.ResponseWriter, r *http.Request) {
	hotel_id := r.URL.Query().Get("hotel_id")

	if hotel_id == "" {
		http.Error(w, "User ID kiritishingiz lozim", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(hotel_id)
	if err != nil {
		http.Error(w, "Invalid hotelID parameter", http.StatusBadRequest)
		return
	}

	hotel := hotelproto.GetHotelDetailsRequest{
		HotelID: int64(id),
	}

	resp, err := h.Client.GetHotelDetails(r.Context(), &hotel)
	if err != nil {
		log.Println("h.Client.GetHotelDetails:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("ERROR: Ma'lumotni encode qilishda xatolik...")
		http.Error(w, "Ma'lumotni encode qilishda xatolik...", http.StatusInternalServerError)
		return
	}
}

// @Router      /api/hotels/rooms/availability [get]
// @Summary     Check Room Availability
// @Description Bu endpoint ma'lum bir mehmonxona haqidagi mavjud xonalar ma'lumotlarni olish uchun ishlatiladi
// @Security    BearerAuth
// @Tags        Hotel
// @Accept      json
// @Produce     json
// @Param		hotelID query string true "Hotel ID: "
// @Success     200 {object} hotelproto.CheckRoomAvailabilityRequest "Room Availability"
// @Failure     400 {object} models.StandartError "Bad request error"
// @Failure 	403 {object} models.ForbiddenError "Forbidden error"
// @Failure     500 {object} models.StandartError "Internal server error"
func (h *HotelClient) CheckRoomAvailability(w http.ResponseWriter, r *http.Request) {
	hotelID := r.URL.Query().Get("hotelID")
	if hotelID == "" {
		http.Error(w, "Missing hotelID parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(hotelID)
	if err != nil {
		http.Error(w, "Invalid hotelID parameter", http.StatusBadRequest)
		return
	}

	grpcReq := &hotelproto.CheckRoomAvailabilityRequest{
		HotelID: int64(id),
	}

	resp, err := h.Client.CheckRoomAvailability(r.Context(), grpcReq)
	if err != nil {
		log.Println("h.Client.CheckRoomAvailability:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("ERROR: Ma'lumotni encode qilishda xatolik...")
		http.Error(w, "Ma'lumotni encode qilishda xatolik...", http.StatusInternalServerError)
		return
	}
}

// @Router      /api/hotels/delete [delete]
// @Summary     Delete room
// @Description Bu endpoint ma'lum bir mehmonxonani o'chirish uchun ishlatiladi
// @Security    BearerAuth
// @Tags        Hotel
// @Accept      json
// @Produce     json
// @Param		hotelID query string true "Hotel ID: "
// @Success     200 {object} hotelproto.DeleteHotelRequest "Delete Room"
// @Failure     400 {object} models.StandartError "Bad request error"
// @Failure 	403 {object} models.ForbiddenError "Forbidden error"
// @Failure     500 {object} models.StandartError "Internal server error"
func (h *HotelClient) DeleteHotel(w http.ResponseWriter, r *http.Request) {
	hotelID := r.URL.Query().Get("hotelID")
	if hotelID == "" {
		http.Error(w, "Missing hotelID parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(hotelID)
	if err != nil {
		http.Error(w, "Invalid hotelID parameter", http.StatusBadRequest)
		return
	}

	grpcReq := &hotelproto.DeleteHotelRequest{
		HotelID: int64(id),
	}

	resp, err := h.Client.DeleteHotel(r.Context(), grpcReq)
	if err != nil {
		log.Println("h.Client.CheckRoomAvailability:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("ERROR: Ma'lumotni encode qilishda xatolik...")
		http.Error(w, "Ma'lumotni encode qilishda xatolik...", http.StatusInternalServerError)
		return
	}
}

// >--------------ROOMS--------------<

// CreateRoom yaratadi yangi xonani
// @Router      /api/rooms [post]
// @Summary     Create room
// @Description Bu endpoint yangi xona yaratish uchun ishlatiladi
// @Security    BearerAuth
// @Tags        Room
// @Accept      json
// @Produce     json
// @Param       body body hotelproto.CreateRoomRequest true "CreateRoomRequest"
// @Success     201 {object} hotelproto.CreateRoomResponse "Room created"
// @Failure     400 {object} models.StandartError "Bad request error"
// @Failure 	403 {object} models.ForbiddenError "Forbidden error"
// @Failure     500 {object} models.StandartError "Internal server error"
func (h *HotelClient) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var room hotelproto.CreateRoomRequest

	if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
		log.Println("CreateRoom-json.Decode:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.Client.CreateRoom(r.Context(), &room)
	if err != nil {
		log.Println("h.Client.CreateRoom:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("ERROR: Ma'lumotni encode qilishda xatolik...")
		http.Error(w, "Ma'lumotni encode qilishda xatolik...", http.StatusInternalServerError)
		return
	}
}

// GetRoomByID ma'lum bir xona haqida ma'lumotlarni qaytaradi
// @Router      /api/roomsbyid [get]
// @Summary     Get room by ID
// @Description Bu endpoint ma'lum bir xona haqidagi ma'lumotlarni olish uchun ishlatiladi
// @Security    BearerAuth
// @Tags        Room
// @Accept      json
// @Produce     json
// @Param       room_ID query string true "Room ID"
// @Param       hotel_ID query string true "Hotel ID"
// @Success     200 {object} hotelproto.GetRoomDetailsResponse "Room details"
// @Failure     400 {object} models.StandartError "Bad request error"
// @Failure     401 {object} models.UnauthorizedError "Unauthorized"
// @Failure     500 {object} models.StandartError "Internal server error"
func (h *HotelClient) GetRoomByID(w http.ResponseWriter, r *http.Request) {
	room_Id := r.URL.Query().Get("room_ID")
	hotel_ID := r.URL.Query().Get("hotel_ID")

	if room_Id == "" {
		http.Error(w, "Room ID kiritishingiz lozim", http.StatusBadRequest)
		return
	}

	if hotel_ID == "" {
		http.Error(w, "Hotel ID kiritishingiz lozim", http.StatusBadRequest)
		return
	}

	rid, err := strconv.Atoi(room_Id)
	if err != nil {
		http.Error(w, "Invalid hotelID parameter", http.StatusBadRequest)
		return
	}

	hid, err := strconv.Atoi(hotel_ID)
	if err != nil {
		http.Error(w, "Invalid hotelID parameter", http.StatusBadRequest)
		return
	}

	roomID := hotelproto.GetRoomDetailsRequest{
		RoomID:  int64(rid),
		HotelID: int64(hid),
	}

	resp, err := h.Client.GetRoomByID(r.Context(), &roomID)
	if err != nil {
		log.Println("h.Client.GetRoomByID:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("ERROR: Ma'lumotni encode qilishda xatolik...")
		http.Error(w, "Ma'lumotni encode qilishda xatolik...", http.StatusInternalServerError)
		return
	}
}

// GetRooms ma'lum bir mehmonxona uchun barcha xonalarni qaytaradi
// @Router      /api/rooms [get]
// @Summary     Get rooms by hotel ID
// @Description Bu endpoint ma'lum bir mehmonxona uchun barcha xonalarni olish uchun ishlatiladi
// @Security    BearerAuth
// @Tags        Room
// @Accept      json
// @Produce     json
// @Param       hotelID query string true "Hotel ID"
// @Success     200 {object} hotelproto.GetRoomsResponse "Rooms list"
// @Failure     400 {object} models.StandartError "Bad request error"
// @Failure     401 {object} models.UnauthorizedError "Unauthorized"
// @Failure     500 {object} models.StandartError "Internal server error"
func (h *HotelClient) GetRooms(w http.ResponseWriter, r *http.Request) {
	hotelIDStr := r.URL.Query().Get("hotelID")
	if hotelIDStr == "" {
		http.Error(w, "Missing hotelID parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(hotelIDStr)
	if err != nil {
		http.Error(w, "Invalid hotelID parameter", http.StatusBadRequest)
		return
	}

	resp, err := h.Client.GetRooms(r.Context(), &hotelproto.GetRoomsRequest{HotelID: int64(id)})
	if err != nil {
		log.Println("h.Client.GetRooms:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("ERROR: Ma'lumotni encode qilishda xatolik...")
		http.Error(w, "Ma'lumotni encode qilishda xatolik...", http.StatusInternalServerError)
		return
	}
}

// UpdateRoom xonani yangilaydi
// @Router      /api/rooms/update [put]
// @Summary     Update room
// @Description Bu endpoint ma'lum bir xona ma'lumotlarini yangilash uchun ishlatiladi
// @Security    BearerAuth
// @Tags        Room
// @Accept      json
// @Produce     json
// @Param       roomID query string true "Room ID"
// @Param       body body hotelproto.UpdateRoomRequest true "UpdateRoomRequest"
// @Success     200 {object} hotelproto.UpdateRoomResponse "Room updated"
// @Failure     400 {object} models.StandartError "Bad request error"
// @Failure     401 {object} models.UnauthorizedError "Unauthorized"
// @Failure     500 {object} models.StandartError "Internal server error"
func (h *HotelClient) UpdateRoom(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("roomID")
	if roomID == "" {
		http.Error(w, "Missing hotelID parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(roomID)
	if err != nil {
		http.Error(w, "Invalid hotelID parameter", http.StatusBadRequest)
		return
	}

	var updateReq hotelproto.UpdateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		log.Println("UpdateRoom-json.Decode:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	updateReq.RoomID = int64(id)

	resp, err := h.Client.UpdateRoom(r.Context(), &updateReq)
	if err != nil {
		log.Println("h.Client.UpdateRoom:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("ERROR: Ma'lumotni encode qilishda xatolik...")
		http.Error(w, "Ma'lumotni encode qilishda xatolik...", http.StatusInternalServerError)
		return
	}
}

// DeleteRoom xonani o'chiradi
// @Router      /api/rooms/delete [delete]
// @Summary     Delete room
// @Description Bu endpoint ma'lum bir xonani o'chirish uchun ishlatiladi
// @Security    BearerAuth
// @Tags        Room
// @Accept      json
// @Produce     json
// @Param       roomID query int64 true "Room ID"
// @Success     200 {object} hotelproto.DeleteRoomResponse "Room deleted"
// @Failure     400 {object} models.StandartError "Bad request error"
// @Failure     401 {object} models.UnauthorizedError "Unauthorized"
// @Failure     500 {object} models.StandartError "Internal server error"
func (h *HotelClient) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("roomID")
	if roomID == "" {
		http.Error(w, "Missing hotelID parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(roomID)
	if err != nil {
		http.Error(w, "Invalid hotelID parameter", http.StatusBadRequest)
		return
	}

	resp, err := h.Client.DeleteRoom(r.Context(), &hotelproto.DeleteRoomRequest{RoomID: int64(id)})
	if err != nil {
		log.Println("h.Client.DeleteRoom:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("ERROR: Ma'lumotni encode qilishda xatolik...")
		http.Error(w, "Ma'lumotni encode qilishda xatolik...", http.StatusInternalServerError)
		return
	}
}
