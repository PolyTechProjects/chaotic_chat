package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PolyTechProjects/chaotic_chat/auth/src/internal/client"
	"github.com/PolyTechProjects/chaotic_chat/auth/src/internal/dto"
	"github.com/PolyTechProjects/chaotic_chat/auth/src/internal/service"
	"github.com/PolyTechProjects/chaotic_chat/auth/src/internal/validator"
)

type AuthController struct {
	authService    *service.AuthService
	userMgmtClient *client.UserMgmtGRPCClient
}

func NewAuthController(authService *service.AuthService, userMgmtClient *client.UserMgmtGRPCClient) *AuthController {
	return &AuthController{
		authService:    authService,
		userMgmtClient: userMgmtClient,
	}
}

func (a *AuthController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = validator.ValidateLogin(req.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = validator.ValidateName(req.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = validator.ValidatePassword(req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accessToken, userId, err := a.authService.Register(req.Login, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = a.userMgmtClient.PerformAddUser(r.Context(), userId.String(), req.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	register := dto.RegisterResponse{
		UserId: userId.String(),
	}
	registerResp, err := json.Marshal(register)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Set-Cookie", fmt.Sprintf("Authorization=%s; HttpOnly", accessToken))
	w.Write(registerResp)
}

func (a *AuthController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	accessToken, err := a.authService.Login(req.Login, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Set-Cookie", fmt.Sprintf("Authorization=%s; HttpOnly", accessToken))
}
