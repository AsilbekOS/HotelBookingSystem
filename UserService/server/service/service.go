package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
	"userService/internal/hash"
	randcode "userService/internal/randCode"
	"userService/internal/sms"
	tokens "userService/internal/token"
	"userService/proto"

	"github.com/redis/go-redis/v9"
)

type UserService struct {
	proto.UnimplementedUserServiceServer
	DB  *sql.DB
	RDB *redis.Client
}

func NewUserService(db *sql.DB, rdb *redis.Client) *UserService {
	return &UserService{DB: db, RDB: rdb}
}

func (s *UserService) RegisterUser(ctx context.Context, req *proto.RegisterUserRequest) (*proto.RegisterUserResponse, error) {
	if !sms.IsValidEmail(req.Email) {
		log.Println("Invalid email address")
		return nil, errors.New("ERROR: invalid email address")
	}

	passwordHash, err := hash.HashPassword(req.Password)
	if err != nil {
		log.Printf("password hashing error: %v", err)
		return nil, errors.New("ERROR: password hashing failed")
	}

	verifyCode := randcode.RandomCode()
	verificationKey := fmt.Sprintf("VerifyCode:%s", req.Email)

	err = s.RDB.Set(ctx, verificationKey, verifyCode, 3*time.Minute).Err()
	if err != nil {
		log.Printf("Error saving verification code to Redis: %v", err)
		return nil, errors.New("ERROR: could not save verification code")
	}

	err = sms.SendSMS(verifyCode, req.Email)
	if err != nil {
		log.Println("Error sending SMS:", err)
		return nil, errors.New("ERROR: could not send SMS")
	}

	userData := map[string]interface{}{
		"username":      req.Username,
		"email":         req.Email,
		"password_hash": passwordHash,
	}

	userDataJSON, err := json.Marshal(userData)
	if err != nil {
		log.Printf("Error marshaling user data: %v", err)
		return nil, errors.New("ERROR: could not serialize user data")
	}

	err = s.RDB.Set(ctx, req.Email, userDataJSON, 3*time.Minute).Err()
	if err != nil {
		log.Printf("Error saving user data to Redis: %v", err)
		return nil, errors.New("ERROR: could not register user")
	}

	return &proto.RegisterUserResponse{
		Email:   req.Email,
		Message: "A verification code has been sent to you. Check your email",
	}, nil
}

func (s *UserService) VerifyUser(ctx context.Context, req *proto.VerifyRequest) (*proto.VerifyResponse, error) {
	verificationKey := fmt.Sprintf("VerifyCode:%s", req.Email)
	verCode, err := s.RDB.Get(ctx, verificationKey).Result()
	if err == redis.Nil {
		log.Println("Verification code not found in Redis")
		return nil, errors.New("verification code not found")
	} else if err != nil {
		log.Println("Error getting verification code from Redis:", err)
		return nil, errors.New("ERROR: could not retrieve verification code")
	}
	log.Println(verCode)

	if verCode != req.Verifycode {
		log.Println("Incorrect verification code")
		return nil, errors.New("incorrect verification code")
	}

	userDataJSON, err := s.RDB.Get(ctx, req.Email).Result()
	if err == redis.Nil {
		log.Println("User not found in Redis")
		return nil, errors.New("user not found")
	} else if err != nil {
		log.Println("Error getting user data from Redis:", err)
		return nil, errors.New("ERROR: could not retrieve user data")
	}

	var userData map[string]interface{}
	if err := json.Unmarshal([]byte(userDataJSON), &userData); err != nil {
		log.Println("Failed to unmarshal user data:", err)
		return nil, errors.New("invalid user data format")
	}
	log.Println(userData)
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING user_id`
	var userID string
	err = s.DB.QueryRow(query, userData["username"], userData["email"], userData["password_hash"]).Scan(&userID)
	if err != nil {
		log.Printf("Error inserting user into DB: %v", err)
		return nil, errors.New("ERROR: could not complete registration")
	}

	err = s.RDB.Del(ctx, verificationKey).Err()
	if err != nil {
		log.Printf("Error deleting verification code from Redis: %v", err)
	}

	return &proto.VerifyResponse{
		UserId:   userID,
		Username: userData["username"].(string),
		Email:    userData["email"].(string),
	}, nil
}

func (s *UserService) LoginUser(ctx context.Context, req *proto.LoginUserRequest) (*proto.LoginUserResponse, error) {
	var user proto.LoginUserRequest
	query := "select user_id, username, password_hash from users"

	var user_id int
	err := s.DB.QueryRow(query).Scan(&user_id, &user.Username, &user.Password)
	if err != nil {
		log.Println("Hech qanday foydalanuvchi topilmadi:", err)
		return nil, errors.New("error: Hech qanday foydalanuvchi topilmadi")
	}

	if !hash.CheckPasswordHash(req.Password, user.Password) {
		log.Println("ERROR: invalid password:", err)
		return nil, errors.New("ERROR: invalid password")
	}

	token, err := tokens.CreateToken(user_id)
	if err != nil {
		log.Println("Token yaratishda xatolik:", err)
		return nil, errors.New(err.Error())
	}

	resp := proto.LoginUserResponse{
		Token:     token,
		ExpiresIn: 60,
	}

	return &resp, nil
}

func (s *UserService) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	if !tokens.ValidateToken(req.Token) {
		log.Println("error: invalid token")
		return nil, errors.New("error: invalid token")
	}

	var user proto.GetUserResponse

	query := "SELECT user_id, username, email, created_at FROM users WHERE user_id = $1"
	err := s.DB.QueryRow(query, req.UserId).Scan(&user.UserId, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		log.Println("Databasedan ma'lumot olishda xatolik")
		return nil, err
	}

	log.Println("Successfully get user info")
	return &user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	if !tokens.ValidateToken(req.Token) {
		log.Println("error: invalid token")
		return nil, errors.New("error: invalid token")
	}

	query := "DELETE FROM users WHERE user_id = $1"
	_, err := s.DB.Exec(query, req.UserId)
	if err != nil {
		log.Println("Foydalanuvchini o'chirishda xatolik:", err)
		return nil, err
	}

	log.Println("Successfully deleted user")
	return &proto.DeleteUserResponse{Success: true}, nil
}
