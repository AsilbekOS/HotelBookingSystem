package router

import (
	"apigateway/api/handler"
	booking "apigateway/grpc_client/booking"
	"apigateway/middleware"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/time/rate"
)

func BookingRouter() {
	port := os.Getenv("USER_PORT")

	bClient := booking.BookingApiConn("bookingservice" + port)

	bookingHandler := handler.NewBookingClient(bClient)

	limiter := rate.NewLimiter(1, 3)

	http.Handle("POST /api/bookings", middleware.CORSMiddleware(middleware.RateLimiter(limiter, http.HandlerFunc(bookingHandler.CreateBooking))))
	http.Handle("GET /api/bookings/", middleware.CORSMiddleware(middleware.RateLimiter(limiter, http.HandlerFunc(bookingHandler.GetBooking))))
	http.Handle("PUT /api/bookings/", middleware.CORSMiddleware(middleware.RateLimiter(limiter, http.HandlerFunc(bookingHandler.UpdateBooking))))
	http.Handle("DELETE /api/bookings/", middleware.CORSMiddleware(middleware.RateLimiter(limiter, http.HandlerFunc(bookingHandler.CancelBooking))))
	http.Handle("GET /api/users/", middleware.CORSMiddleware(middleware.RateLimiter(limiter, http.HandlerFunc(bookingHandler.ListUserBookings))))
}
