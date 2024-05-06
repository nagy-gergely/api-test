package user

import (
	"encoding/json"
	"net/http"

	"github.com/nagy-gergely/api-test/services/auth"
	"github.com/nagy-gergely/api-test/types"
	"github.com/nagy-gergely/api-test/utils"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (handler *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /login", handler.login)
	router.HandleFunc("POST /register", handler.register)
}

func (handler *Handler) login(w http.ResponseWriter, r *http.Request) {
	var loginPayload types.LoginPayload

	if err := json.NewDecoder(r.Body).Decode(&loginPayload); err != nil {
		utils.Logger.Println("error decoding request body:", err)
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := utils.Validate.Struct(loginPayload); err != nil {
		utils.Logger.Println("error validating request body:", err)
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := handler.store.GetUserByEmail(loginPayload.Email)
	if err != nil {
		utils.Logger.Println("error user not found:", loginPayload.Email)
		utils.WriteError(w, http.StatusNotFound, "user not found")
		return
	}

	if err := auth.ComparePasswords(user.Password, loginPayload.Password); err != nil {
		utils.Logger.Println("error invalid password:", err)
		utils.WriteError(w, http.StatusUnauthorized, "invalid password")
		return
	}

	token, err := auth.GenerateToken(user)
	if err != nil {
		utils.Logger.Println("error generating token:", err)
		utils.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (handler *Handler) register(w http.ResponseWriter, r *http.Request) {
	var registerPayload types.RegisterPayload

	if err := json.NewDecoder(r.Body).Decode(&registerPayload); err != nil {
		utils.Logger.Println("error decoding request body:", err)
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := utils.Validate.Struct(registerPayload); err != nil {
		utils.Logger.Println("error validating request body:", err)
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if user, err := handler.store.GetUserByEmail(registerPayload.Email); err == nil {
		utils.Logger.Println("error user already exists:", user.Email)
		utils.WriteError(w, http.StatusConflict, "user already exists")
		return
	}

	hashedPassword, err := auth.HashPassword(registerPayload.Password)
	if err != nil {
		utils.Logger.Println("error hashing password:", err)
		utils.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := handler.store.CreateUser(types.User{
		FirstName: registerPayload.FirstName,
		LastName:  registerPayload.LastName,
		Email:     registerPayload.Email,
		Password:  hashedPassword,
	}); err != nil {
		utils.Logger.Println("error creating user:", err)
		utils.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "user created"})
}
