package main

import (
	_ "apigateway/api/docs"
	"apigateway/api/router"
	"apigateway/middleware"
	"apigateway/websocket"
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"

	httpSwagger "github.com/swaggo/http-swagger"
)

// New ...
// @title Project: HOTEL BOOKING SYSTEM
// @description This swagger UI was created by Asilbek Xolmatov
// @version 1.0

// @host localhost:7777
// @schemes https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @tags user
func main() {
	port := os.Getenv("API_PORT")

	http.Handle("/swagger/", httpSwagger.WrapHandler)
	http.HandleFunc("/ws", websocket.WebsocketHandler)

	handler := middleware.CORSMiddleware(http.DefaultServeMux)

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	server := &http.Server{
		Addr:      port,
		Handler:   handler,
		TLSConfig: tlsConfig,
	}

	router.UserRouter()
	router.HotelRouter()
	router.BookingRouter()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		log.Println("O'chirish signali qabul qilindi, xushmuomalalik bilan o'chirishni boshladi...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatal("Server shutdown error:", err)
		}

		log.Println("Server qulay tarzda yopildi.")
	}()

	log.Println("Server is running on port localhost" + port)
	if err := server.ListenAndServeTLS("./tls/localhost.pem", "./tls/localhost-key.pem"); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server failed to start: ", err)
	}
}
