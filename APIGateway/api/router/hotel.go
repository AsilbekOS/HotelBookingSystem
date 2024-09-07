package router

import (
	"apigateway/api/handler"
	hotel "apigateway/grpc_client/hotel"
	"apigateway/middleware"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/time/rate"
)

func HotelRouter() {
	port := os.Getenv("USER_PORT")

	hClient := hotel.HotelApiConn("hotelservice" + port)

	hotelHandler := handler.NewHotelClient(hClient)

	limiter := rate.NewLimiter(1, 3)

	http.Handle("POST /api/hotels", middleware.CORSMiddleware(middleware.RateLimiter(limiter, http.HandlerFunc(hotelHandler.CreateHotel))))
	http.Handle("GET /api/hotels", middleware.CORSMiddleware(middleware.RateLimiter(limiter, http.HandlerFunc(hotelHandler.GetHotels))))
	http.Handle("GET /api/hotelsid", middleware.CORSMiddleware(middleware.RateLimiter(limiter, http.HandlerFunc(hotelHandler.GetHotelDetails))))
	http.Handle("GET /api/hotels/rooms/availability", middleware.CORSMiddleware(middleware.RateLimiter(limiter, http.HandlerFunc(hotelHandler.CheckRoomAvailability))))
	http.Handle("DELETE /api/hotels/delete", middleware.CORSMiddleware(middleware.RateLimiter(limiter, http.HandlerFunc(hotelHandler.DeleteHotel))))

	http.Handle("POST /api/rooms", middleware.CORSMiddleware(middleware.RateLimiter(limiter, http.HandlerFunc(hotelHandler.CreateRoom))))
	http.Handle("GET /api/roomsbyid", middleware.CORSMiddleware(middleware.RateLimiter(limiter, http.HandlerFunc(hotelHandler.GetRoomByID))))
	http.Handle("GET /api/rooms", middleware.CORSMiddleware(middleware.RateLimiter(limiter, http.HandlerFunc(hotelHandler.GetRooms))))
	http.Handle("PUT /api/rooms/update", middleware.CORSMiddleware(middleware.RateLimiter(limiter, http.HandlerFunc(hotelHandler.UpdateRoom))))
	http.Handle("DELETE /api/rooms/delete", middleware.CORSMiddleware(middleware.RateLimiter(limiter, http.HandlerFunc(hotelHandler.DeleteRoom))))
}
