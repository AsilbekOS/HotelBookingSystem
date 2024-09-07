package handler

import (
	bookingproto "apigateway/proto/booking"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type BookingClient struct {
	Client bookingproto.BookingServiceClient
}

func NewBookingClient(Cl bookingproto.BookingServiceClient) *BookingClient {
	return &BookingClient{Client: Cl}
}

// CreateBooking godoc
// @Router      /api/bookings [post]
// @Summary     Create and confirm a booking
// @Description Bu endpoint yangi bronlash yaratish va tasdiqlash uchun ishlatiladi
// @Security    BearerAuth
// @Tags        Booking
// @Accept      json
// @Produce     json
// @Param       body body bookingproto.CreateBookingRequest true "Create booking request"
// @Success     200 {object} bookingproto.BookingResponse "Booking created and confirmed"
// @Failure     400 {object} models.StandartError "Bad request error"
// @Failure     401 {object} models.UnauthorizedError "Unauthorized"
// @Failure     500 {object} models.StandartError "Internal server error"
func (b *BookingClient) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var req bookingproto.CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := b.Client.CreateBooking(r.Context(), &req)
	if err != nil {
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

// GetBooking godoc
// @Router      /api/bookings/{bookingID} [get]
// @Summary     Get booking details
// @Description Bu endpoint ma'lum bir bronlash tafsilotlarini olish uchun ishlatiladi
// @Security    BearerAuth
// @Tags        Booking
// @Accept      json
// @Produce     json
// @Param       bookingID query int64 true "Booking ID"
// @Success     200 {object} bookingproto.BookingResponse "Booking details"
// @Failure     400 {object} models.StandartError "Bad request error"
// @Failure     401 {object} models.UnauthorizedError "Unauthorized"
// @Failure     500 {object} models.StandartError "Internal server error"
func (b *BookingClient) GetBooking(w http.ResponseWriter, r *http.Request) {
	bookingIDStr := r.URL.Query().Get("bookingID")
	bookingID, err := strconv.ParseInt(bookingIDStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := &bookingproto.GetBookingRequest{BookingID: bookingID}
	resp, err := b.Client.GetBooking(r.Context(), req)
	if err != nil {
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

// UpdateBooking godoc
// @Router      /api/bookings/{bookingID} [put]
// @Summary     Update booking details
// @Description Bu endpoint mavjud bronlash tafsilotlarini yangilash uchun ishlatiladi
// @Security    BearerAuth
// @Tags        Booking
// @Accept      json
// @Produce     json
// @Param       bookingID query int64 true "Booking ID"
// @Param       body body bookingproto.UpdateBookingRequest true "Update booking details"
// @Success     200 {object} bookingproto.BookingResponse "Updated booking details"
// @Failure     400 {object} models.StandartError "Bad request error"
// @Failure     401 {object} models.UnauthorizedError "Unauthorized"
// @Failure     500 {object} models.StandartError "Internal server error"
func (b *BookingClient) UpdateBooking(w http.ResponseWriter, r *http.Request) {
	bookingIDStr := r.URL.Query().Get("bookingID")
	bookingID, err := strconv.ParseInt(bookingIDStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req bookingproto.UpdateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.BookingID = bookingID

	resp, err := b.Client.UpdateBooking(r.Context(), &req)
	if err != nil {
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

// CancelBooking godoc
// @Router      /api/bookings/{bookingID} [delete]
// @Summary     Cancel a booking
// @Description Bu endpoint ma'lum bir bronlashni bekor qilish uchun ishlatiladi
// @Security    BearerAuth
// @Tags        Booking
// @Accept      json
// @Produce     json
// @Param       bookingID query int64 true "Booking ID"
// @Success     200 {object} bookingproto.CancelBookingResponse "Booking cancelled"
// @Failure     400 {object} models.StandartError "Bad request error"
// @Failure     401 {object} models.UnauthorizedError "Unauthorized"
// @Failure     500 {object} models.StandartError "Internal server error"
func (b *BookingClient) CancelBooking(w http.ResponseWriter, r *http.Request) {
	bookingIDStr := r.URL.Query().Get("bookingID")
	bookingID, err := strconv.ParseInt(bookingIDStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := &bookingproto.CancelBookingRequest{BookingID: bookingID}
	resp, err := b.Client.CancelBooking(r.Context(), req)
	if err != nil {
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

// ListUserBookings godoc
// @Router      /api/users/ [get]
// @Summary     List all bookings for a user
// @Description Bu endpoint foydalanuvchi uchun barcha bronlashlarni olish uchun ishlatiladi
// @Security    BearerAuth
// @Tags        Booking
// @Accept      json
// @Produce     json
// @Param       userID query int64 true "User ID"
// @Success     200 {object} bookingproto.ListUserBookingsResponse "List of user bookings"
// @Failure     400 {object} models.StandartError "Bad request error"
// @Failure     401 {object} models.UnauthorizedError "Unauthorized"
// @Failure     500 {object} models.StandartError "Internal server error"
func (b *BookingClient) ListUserBookings(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := &bookingproto.ListUserBookingsRequest{UserID: userID}
	resp, err := b.Client.ListUserBookings(r.Context(), req)
	if err != nil {
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
