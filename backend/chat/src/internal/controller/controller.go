package controller

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/PolyTechProjects/chaotic_chat/chat/src/internal/client"
	"github.com/PolyTechProjects/chaotic_chat/chat/src/internal/dto"
	"github.com/PolyTechProjects/chaotic_chat/chat/src/internal/service"
	"github.com/PolyTechProjects/chaotic_chat/chat/src/internal/validator"
	"github.com/google/uuid"
)

type ChatManagementController struct {
	service    *service.ChatManagementService
	authClient *client.AuthGRPCClient
}

func NewChatManagementController(service *service.ChatManagementService, authClient *client.AuthGRPCClient) *ChatManagementController {
	return &ChatManagementController{
		service:    service,
		authClient: authClient,
	}
}

func (c *ChatManagementController) CreateChatHandler(w http.ResponseWriter, r *http.Request) {
	var chatReq dto.CreateChatRequest
	err := json.NewDecoder(r.Body).Decode(&chatReq)
	if err != nil {
		slog.Error("Failed to decode request", "error", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authResp, err := c.authClient.PerformAuthorize(r.Context(), r)
	if err != nil {
		slog.Error("Authorization error", "error", err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if len(chatReq.ParticipantsIds) < 2 {
		http.Error(w, "Not enough participants", http.StatusBadRequest)
		return
	}

	err = validator.ValidateChatName(chatReq.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	chat, err := c.service.CreateChat(&chatReq)
	if err != nil {
		slog.Error("Failed to create chat", "error", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	chatResp, err := json.Marshal(chat)
	if err != nil {
		slog.Error("Failed to marshal response", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Set-Cookie", fmt.Sprintf("Authorization=%s; HttpOnly", authResp.AccessToken))
	w.Write(chatResp)
}

func (c *ChatManagementController) DeleteChatHandler(w http.ResponseWriter, r *http.Request) {
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !params.Has("chatId") {
		http.Error(w, "URL query params are invalid", http.StatusBadRequest)
	}
	chatId, err := uuid.Parse(params.Get("chatId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authResp, err := c.authClient.PerformAuthorize(r.Context(), r)
	if err != nil {
		slog.Error("Authorization error", "error", err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	chat, err := c.service.GetChat(chatId)
	if err != nil {
		slog.Error("c.service.GetChat() returned error", "error", err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if chat.Chat.CreatorId.String() != authResp.UserId {
		http.Error(w, "You are not the creator of the chat", http.StatusForbidden)
		return
	}

	err = c.service.DeleteChat(chatId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Set-Cookie", fmt.Sprintf("Authorization=%s; HttpOnly", authResp.AccessToken))
}

func (c *ChatManagementController) GetChatHandler(w http.ResponseWriter, r *http.Request) {
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !params.Has("chatId") {
		http.Error(w, "URL query params are invalid", http.StatusBadRequest)
	}
	chatId, err := uuid.Parse(params.Get("chatId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authResp, err := c.authClient.PerformAuthorize(r.Context(), r)
	if err != nil {
		slog.Error("Authorization error", "error", err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	chat, err := c.service.GetChat(chatId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !(slices.Contains(chat.Readers, authResp.UserId) || slices.Contains(chat.Users, authResp.UserId) || slices.Contains(chat.Admins, authResp.UserId)) {
		http.Error(w, "permission denied", http.StatusForbidden)
		return
	}
	chatResp, err := json.Marshal(chat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Set-Cookie", fmt.Sprintf("Authorization=%s; HttpOnly", authResp.AccessToken))
	w.Write(chatResp)
}

func (c *ChatManagementController) UpdateChatHandler(w http.ResponseWriter, r *http.Request) {
	var chatReq dto.UpdateChatRequest
	err := json.NewDecoder(r.Body).Decode(&chatReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authResp, err := c.authClient.PerformAuthorize(r.Context(), r)
	if err != nil {
		slog.Error("Authorization error", "error", err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	chat, err := c.service.GetChat(chatReq.ChatId)
	if err != nil {
		slog.Error("c.service.GetChat() returned error", "error", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !(slices.Contains(chat.Admins, authResp.UserId) || chat.Chat.CreatorId.String() == authResp.UserId) {
		http.Error(w, "permission denied", http.StatusForbidden)
		return
	}

	err = validator.ValidateChatName(chatReq.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.service.UpdateChat(&chatReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	chat, err = c.service.GetChat(chatReq.ChatId)
	if err != nil {
		slog.Error("c.service.GetChat() returned error", "error", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	chatResp, err := json.Marshal(chat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Set-Cookie", fmt.Sprintf("Authorization=%s; HttpOnly", authResp.AccessToken))
	w.Write(chatResp)
}

func (c *ChatManagementController) DeleteUsersInChatHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateUsersRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authResp, err := c.authClient.PerformAuthorize(r.Context(), r)
	if err != nil {
		slog.Error("Authorization error", "error", err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	chat, err := c.service.GetChat(req.ChatId)
	if err != nil {
		slog.Error("c.service.GetChat() returned error", "error", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !(slices.Contains(chat.Admins, authResp.UserId) || chat.Chat.CreatorId.String() == authResp.UserId) {
		http.Error(w, "permission denied", http.StatusForbidden)
		return
	}

	err = c.service.DeleteUsers(req.ChatId, req.UserIds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Set-Cookie", fmt.Sprintf("Authorization=%s; HttpOnly", authResp.AccessToken))
}

func (c *ChatManagementController) AddUsersInChatHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateUsersRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authResp, err := c.authClient.PerformAuthorize(r.Context(), r)
	if err != nil {
		slog.Error("Authorization error", "error", err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	chat, err := c.service.GetChat(req.ChatId)
	if err != nil {
		slog.Error("c.service.GetChat() returned error", "error", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !(slices.Contains(chat.Admins, authResp.UserId) || chat.Chat.CreatorId.String() == authResp.UserId) {
		http.Error(w, "permission denied", http.StatusForbidden)
		return
	}

	err = c.service.AddUsers(req.ChatId, req.UserIds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Set-Cookie", fmt.Sprintf("Authorization=%s; HttpOnly", authResp.AccessToken))
}

func (c *ChatManagementController) MakeReadersUsersInChatHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateReadersRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authResp, err := c.authClient.PerformAuthorize(r.Context(), r)
	if err != nil {
		slog.Error("Authorization error", "error", err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	chat, err := c.service.GetChat(req.ChatId)
	if err != nil {
		slog.Error("c.service.GetChat() returned error", "error", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !(slices.Contains(chat.Admins, authResp.UserId) || chat.Chat.CreatorId.String() == authResp.UserId) {
		http.Error(w, "permission denied", http.StatusForbidden)
		return
	}

	err = c.service.MakeReadersUsers(req.ChatId, req.ReadersIds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Set-Cookie", fmt.Sprintf("Authorization=%s; HttpOnly", authResp.AccessToken))
}

func (c *ChatManagementController) MakeUsersReadersInChatHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateReadersRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authResp, err := c.authClient.PerformAuthorize(r.Context(), r)
	if err != nil {
		slog.Error("Authorization error", "error", err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	chat, err := c.service.GetChat(req.ChatId)
	if err != nil {
		slog.Error("c.service.GetChat() returned error", "error", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !(slices.Contains(chat.Admins, authResp.UserId) || chat.Chat.CreatorId.String() == authResp.UserId) {
		http.Error(w, "permission denied", http.StatusForbidden)
		return
	}

	err = c.service.MakeUsersReaders(req.ChatId, req.ReadersIds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Set-Cookie", fmt.Sprintf("Authorization=%s; HttpOnly", authResp.AccessToken))
}

func (c *ChatManagementController) AddAdminsInChatHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateAdminsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authResp, err := c.authClient.PerformAuthorize(r.Context(), r)
	if err != nil {
		slog.Error("Authorization error", "error", err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	chat, err := c.service.GetChat(req.ChatId)
	if err != nil {
		slog.Error("c.service.GetChat() returned error", "error", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if chat.Chat.CreatorId.String() != authResp.UserId {
		http.Error(w, "permission denied", http.StatusForbidden)
		return
	}

	err = c.service.AddAdmin(req.ChatId, req.AdminsIds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Set-Cookie", fmt.Sprintf("Authorization=%s; HttpOnly", authResp.AccessToken))
}

func (c *ChatManagementController) DeleteAdminsInChatHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateAdminsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authResp, err := c.authClient.PerformAuthorize(r.Context(), r)
	if err != nil {
		slog.Error("Authorization error", "error", err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	chat, err := c.service.GetChat(req.ChatId)
	if err != nil {
		slog.Error("c.service.GetChat() returned error", "error", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if chat.Chat.CreatorId.String() != authResp.UserId {
		http.Error(w, "permission denied", http.StatusForbidden)
		return
	}

	err = c.service.DeleteAdmin(req.ChatId, req.AdminsIds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Set-Cookie", fmt.Sprintf("Authorization=%s; HttpOnly", authResp.AccessToken))
}

func (c *ChatManagementController) JoinChatHandler(w http.ResponseWriter, r *http.Request) {
	joinLink := strings.Split(r.URL.Path, "/")[3]
	authResp, err := c.authClient.PerformAuthorize(r.Context(), r)
	if err != nil {
		slog.Error("Authorization error", "error", err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	userId, err := uuid.Parse(authResp.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.service.JoinChat(joinLink, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Set-Cookie", fmt.Sprintf("Authorization=%s; HttpOnly", authResp.AccessToken))
}
