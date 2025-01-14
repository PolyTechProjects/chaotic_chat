package controller

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/PolyTechProjects/chaotic_chat/media_handler/src/internal/client"
	"github.com/PolyTechProjects/chaotic_chat/media_handler/src/internal/models"
	"github.com/PolyTechProjects/chaotic_chat/media_handler/src/internal/service"
	"github.com/google/uuid"
)

type MediaHandlerController struct {
	mediaHandlerService *service.MediaHandlerService
	authClient          *client.AuthGRPCClient
}

func New(mediaHandlerService *service.MediaHandlerService, authClient *client.AuthGRPCClient) *MediaHandlerController {
	return &MediaHandlerController{mediaHandlerService: mediaHandlerService, authClient: authClient}
}

func (m *MediaHandlerController) UploadMediaHandler(w http.ResponseWriter, r *http.Request) {
	authResp, err := m.authClient.PerformAuthorize(r.Context(), r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	objectType := r.Header.Get("ObjectType")
	objectIdHeader := r.Header.Get("ObjectId")
	objectId, err := uuid.Parse(objectIdHeader)
	if err != nil {
		slog.Error(fmt.Sprintf("uuid.Parse(%s) returned error: %s", objectIdHeader, err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		slog.Error(fmt.Sprintf("r.FormFile returned error: %s", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mediaId, err := m.mediaHandlerService.UploadMedia(objectType, objectId, file, fileHeader)
	if err != nil {
		slog.Error(fmt.Sprintf("m.mediaHandlerService.UploadMedia(%s, %s, file, fileHeader) returned error: %s", objectType, objectIdHeader, err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Set-Cookie", fmt.Sprintf("Authorization=%s; HttpOnly", authResp.AccessToken))
	if objectType != "OBJECT_MESSAGE" {
		resp := &models.UploadMediaResponse{
			ObjectType: objectType,
			ObjectId:   objectIdHeader,
			FileId:     mediaId.String(),
		}
		respBytes, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(respBytes)
	}
}

func (m *MediaHandlerController) GetMediaHandler(w http.ResponseWriter, r *http.Request) {
	authResp, err := m.authClient.PerformAuthorize(r.Context(), r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !params.Has("fileId") {
		http.Error(w, "URL query params are invalid", http.StatusBadRequest)
	}
	fileId, err := uuid.Parse(params.Get("fileId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := m.mediaHandlerService.GetMedia(fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Set-Cookie", fmt.Sprintf("Authorization=%s; HttpOnly", authResp.AccessToken))
	w.Write(res)
}

func (m *MediaHandlerController) DeleteMediaHandler(w http.ResponseWriter, r *http.Request) {
	authResp, err := m.authClient.PerformAuthorize(r.Context(), r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !params.Has("fileId") {
		http.Error(w, "URL query params are invalid", http.StatusBadRequest)
	}
	fileId, err := uuid.Parse(params.Get("fileId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = m.mediaHandlerService.DeleteMedia(fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Set-Cookie", fmt.Sprintf("Authorization=%s; HttpOnly", authResp.AccessToken))
}
