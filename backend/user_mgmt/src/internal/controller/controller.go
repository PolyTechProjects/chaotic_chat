package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/PolyTechProjects/chaotic_chat/user_mgmt/src/internal/client"
	"github.com/PolyTechProjects/chaotic_chat/user_mgmt/src/internal/dto"
	"github.com/PolyTechProjects/chaotic_chat/user_mgmt/src/internal/service"
	"github.com/PolyTechProjects/chaotic_chat/user_mgmt/src/internal/validator"
	"github.com/google/uuid"
)

type UserMgmtController struct {
	userMgmtService *service.UserMgmtService
	authClient      *client.AuthGRPCClient
}

func New(userMgmtService *service.UserMgmtService, authClient *client.AuthGRPCClient) *UserMgmtController {
	return &UserMgmtController{userMgmtService: userMgmtService, authClient: authClient}
}

func (c *UserMgmtController) UpdateAvatarHandler(w http.ResponseWriter, r *http.Request) {
	authResp, err := c.authClient.PerformAuthorize(r.Context(), r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var req dto.UploadProfilePicRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := c.userMgmtService.UpdateAvatar(userId, req.FileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Set-Cookie", fmt.Sprintf("Authorization=%s; HttpOnly", authResp.AccessToken))
	w.Write(resp)
}

func (c *UserMgmtController) InfoUpdateHandler(w http.ResponseWriter, r *http.Request) {
	dto := dto.UpdateInfoRequest{}
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authResp, err := c.authClient.PerformAuthorize(r.Context(), r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userId, err := uuid.Parse(dto.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = validator.ValidateName(dto.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = validator.ValidateDescription(dto.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = validator.ValidateUrlTag(dto.UrlTag)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := c.userMgmtService.UpdateUser(userId, dto.Name, dto.UrlTag, dto.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Set-Cookie", fmt.Sprintf("Authorization=%s; HttpOnly", authResp.AccessToken))
	w.Write(resp)
}

func (c *UserMgmtController) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !params.Has("urlTag") {
		http.Error(w, "URL query params are invalid", http.StatusBadRequest)
	}

	authResp, err := c.authClient.PerformAuthorize(r.Context(), r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	urlTag, err := uuid.Parse(params.Get("urlTag"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := c.userMgmtService.GetUser(urlTag)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Set-Cookie", fmt.Sprintf("Authorization=%s; HttpOnly", authResp.AccessToken))
	w.Write(resp)
}

func (c *UserMgmtController) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	dto := dto.UpdateInfoRequest{}
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authResp, err := c.authClient.PerformAuthorize(r.Context(), r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userId, err := uuid.Parse(dto.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = c.userMgmtService.DeleteUser(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Add("Set-Cookie", fmt.Sprintf("Authorization=%s; HttpOnly", authResp.AccessToken))
}
