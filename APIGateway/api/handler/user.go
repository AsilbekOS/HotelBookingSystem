package handler

import (
	userproto "apigateway/proto/user"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"
)

type UserClient struct {
	Client userproto.UserServiceClient
}

func NewUserClient(cl userproto.UserServiceClient) *UserClient {
	return &UserClient{Client: cl}
}

// @Router		/api/users [post]
// @Summary		Register User
// @Description Bu endpoint yangi foydalanuvchini ro'yxatdan o'tkazish uchun ishlatiladi uchun ishlatiladi
// @Security	BearerAuth
// @Tags		User
// @Accept		json
// @Produce 	json
// @Param		body body userproto.RegisterUserRequest true "RegisterUserRequest"
// @Success 	201 {object} map[string]interface{} "Get verified"
// @Failure 	400 {object} models.StandartError "Bad request error"
// @Failure 	403 {object} models.ForbiddenError "Forbidden error"
// @Failure 	500 {object} models.StandartError "Internal server error"
func (u *UserClient) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user userproto.RegisterUserRequest

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("RegisterUser-io.ReadAll:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	protojson.Unmarshal(bytes, &user)

	resp, err := u.Client.RegisterUser(r.Context(), &user)
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

// @Router		/api/users/verify [POST]
// @Summary		Verify User
// @Description Bu endpoint foydalanuvchi verifikatsiyadan otish uchun ishlatiladi
// @Security	BearerAuth
// @Tags		User
// @Accept		json
// @Produce 	json
// @Param		body body userproto.VerifyRequest true "VerifyRequest"
// @Success 	200 {object} map[string]interface{} "Verify successful"
// @Failure 	400 {object} models.StandartError "Bad request error"
// @Failure 	403 {object} models.ForbiddenError "Forbidden error"
// @Failure 	500 {object} models.StandartError "Internal server error"
func (u *UserClient) VerifyUser(w http.ResponseWriter, r *http.Request) {
	var userReq userproto.VerifyRequest

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("ERROR: POST: bodydan o'qib olishda xatolik...")
		http.Error(w, "POST: bodydan o'qib olishda xatolik...", http.StatusBadRequest)
		return
	}

	err = protojson.Unmarshal(bytes, &userReq)
	if err != nil {
		log.Println("ERROR: POST: unmarshal qilishda xatolik...")
		http.Error(w, "POST: unmarshal qilishda xatolik...", http.StatusBadRequest)
		return
	}

	resp, err := u.Client.VerifyUser(r.Context(), &userReq)
	if err != nil {
		log.Println("ERROR: POST: Serverdan ma'lumot olishda xatolik...")
		http.Error(w, "POST: Serverdan ma'lumot olishda xatolik...", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("ERROR: Ma'lumotni encode qilishda xatolik...")
		http.Error(w, "Ma'lumotni encode qilishda xatolik...", http.StatusInternalServerError)
		return
	}
}

// @Router		/api/users/login [post]
// @Summary		LOGIN User
// @Description Bu endpoint foydalanuvchi profiligi kirish uchun ishlatiladi
// @Security	BearerAuth
// @Tags		User
// @Accept		json
// @Produce 	json
// @Param		body body userproto.LoginUserRequest true "LoginUserRequest"
// @Success 	200 {object} map[string]interface{} "Login successful"
// @Failure 	400 {object} models.StandartError "Bad request error"
// @Failure 	403 {object} models.ForbiddenError "Forbidden error"
// @Failure 	500 {object} models.StandartError "Internal server error"
func (u *UserClient) LoginUser(w http.ResponseWriter, r *http.Request) {
	var userReq userproto.LoginUserRequest

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("ERROR: POST: bodydan o'qib olishda xatolik...")
		http.Error(w, "POST: bodydan o'qib olishda xatolik...", http.StatusBadRequest)
		return
	}

	err = protojson.Unmarshal(bytes, &userReq)
	if err != nil {
		log.Println("ERROR: POST: unmarshal qilishda xatolik...")
		http.Error(w, "POST: unmarshal qilishda xatolik...", http.StatusBadRequest)
		return
	}

	resp, err := u.Client.LoginUser(r.Context(), &userReq)
	if err != nil {
		log.Println("ERROR: POST: Serverdan ma'lumot olishda xatolik...")
		http.Error(w, "POST: Serverdan ma'lumot olishda xatolik...", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("ERROR: Ma'lumotni encode qilishda xatolik...")
		http.Error(w, "Ma'lumotni encode qilishda xatolik...", http.StatusInternalServerError)
		return
	}
}

// @Router		/api/users/profile [get]
// @Summary		PROFILE User
// @Description Bu endpoint foydalanuvchi profilini olish uchun ishlatiladi
// @Security	BearerAuth
// @Tags		User
// @Accept		json
// @Produce 	json
// @Param 		token query string true "Token: "
// @Param 		user_id query string true "User ID: "
// @Success 	200 {object} map[string]interface{} "Successful response"
// @Failure 	400 {object} models.StandartError "Bad request error"
// @Failure 	403 {object} models.ForbiddenError "Forbidden error"
// @Failure 	500 {object} models.StandartError "Internal server error"
func (u *UserClient) GetUser(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	user_id := r.URL.Query().Get("user_id")

	if token == "" {
		log.Println("Token kiritishingiz lozim")
		http.Error(w, "Token kiritishingiz lozim", http.StatusBadRequest)
		return
	}

	if user_id == "" {
		log.Println("User ID kiritishingiz lozim")
		http.Error(w, "User ID kiritishingiz lozim", http.StatusBadRequest)
		return
	}

	userReq := userproto.GetUserRequest{
		Token:  token,
		UserId: user_id,
	}

	resp, err := u.Client.GetUser(r.Context(), &userReq)
	if err != nil {
		log.Println("ERROR: GET: Serverdan ma'lumot olishda xatolik...")
		http.Error(w, "GET: Serverdan ma'lumot olishda xatolik...", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("ERROR: Ma'lumotni encode qilishda xatolik...")
		http.Error(w, "Ma'lumotni encode qilishda xatolik...", http.StatusInternalServerError)
		return
	}
}

// @Router		/api/users/profile/delete [delete]
// @Summary		DELETE User
// @Description Bu endpoint foydalanuvchi profilini o'chirish uchun ishlatiladi
// @Security	BearerAuth
// @Tags		User
// @Accept		json
// @Produce 	json
// @Param 		token query string true "Token: "
// @Param 		user_id query string true "User ID: "
// @Success 	200 {object} map[string]interface{} "Delete successful"
// @Failure 	400 {object} models.StandartError "Bad request error"
// @Failure 	403 {object} models.ForbiddenError "Forbidden error"
// @Failure 	500 {object} models.StandartError "Internal server error"
func (u *UserClient) DeleteUser(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	user_id := r.URL.Query().Get("user_id")

	if token == "" {
		log.Println("Token kiritishingiz lozim")
		http.Error(w, "Token kiritishingiz lozim", http.StatusBadRequest)
		return
	}

	if user_id == "" {
		log.Println("User ID kiritishingiz lozim")
		http.Error(w, "User ID kiritishingiz lozim", http.StatusBadRequest)
		return
	}

	userReq := userproto.DeleteUserRequest{
		Token:  token,
		UserId: user_id,
	}

	resp, err := u.Client.DeleteUser(r.Context(), &userReq)
	if err != nil {
		log.Println("ERROR: POST: Serverdan ma'lumot olishda xatolik...")
		http.Error(w, "POST: Serverdan ma'lumot olishda xatolik...", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("ERROR: Ma'lumotni encode qilishda xatolik...")
		http.Error(w, "Ma'lumotni encode qilishda xatolik...", http.StatusInternalServerError)
		return
	}
}
