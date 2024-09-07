package router

import (
	"apigateway/api/handler"
	"apigateway/grpc_client/user"
	"apigateway/middleware"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/time/rate"
)

func UserRouter() {
	port := os.Getenv("USER_PORT")
	_ = port

	uClient := user.UserApiConn("userservice:8081")

	userHandler := handler.NewUserClient(uClient)

	limiter := rate.NewLimiter(1, 3)

	http.Handle("POST /api/users", middleware.CORSMiddleware(middleware.RateLimiter(limiter, http.HandlerFunc(userHandler.RegisterUser))))
	http.Handle("POST /api/users/verify", middleware.CORSMiddleware(middleware.RateLimiter(limiter, http.HandlerFunc(userHandler.VerifyUser))))
	http.Handle("POST /api/users/login", middleware.CORSMiddleware(middleware.RateLimiter(limiter, http.HandlerFunc(userHandler.LoginUser))))
	http.Handle("GET /api/users/profile", middleware.CORSMiddleware(middleware.RateLimiter(limiter, http.HandlerFunc(userHandler.GetUser))))
	http.Handle("DELETE /api/users/profile/delete", middleware.CORSMiddleware(middleware.RateLimiter(limiter, http.HandlerFunc(userHandler.DeleteUser))))
}
