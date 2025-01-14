package server

import (
	"net/http"

	"github.com/PolyTechProjects/chaotic_chat/media_handler/src/internal/controller"
)

type HttpServer struct {
	mediaHandlerController *controller.MediaHandlerController
}

func NewHttpServer(mediaHandlerController *controller.MediaHandlerController) *HttpServer {
	return &HttpServer{mediaHandlerController: mediaHandlerController}
}

func (h *HttpServer) StartServer() {
	http.HandleFunc("POST /media/uploads", h.mediaHandlerController.UploadMediaHandler)
	http.HandleFunc("GET /media/uploads", h.mediaHandlerController.GetMediaHandler)
	http.HandleFunc("DELETE /media/uploads", h.mediaHandlerController.DeleteMediaHandler)
}
