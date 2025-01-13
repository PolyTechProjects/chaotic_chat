package server

import (
	"net/http"

	"github.com/PolyTechProjects/chaotic_chat/user_mgmt/src/internal/controller"
)

type HttpServer struct {
	userMgmtController *controller.UserMgmtController
}

func NewHttpServer(userMgmtController *controller.UserMgmtController) *HttpServer {
	return &HttpServer{userMgmtController: userMgmtController}
}

func (h *HttpServer) StartServer() {
	http.HandleFunc("POST /user/profile/pic", h.userMgmtController.UpdateAvatarHandler)
	http.HandleFunc("PUT /user/profile", h.userMgmtController.InfoUpdateHandler)
	http.HandleFunc("GET /user/profile", h.userMgmtController.GetUserHandler)
	http.HandleFunc("DELETE /user/profile", h.userMgmtController.DeleteUserHandler)
}
