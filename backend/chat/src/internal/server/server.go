package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	chat_mgmt "github.com/PolyTechProjects/chaotic_chat/chat/src/gen/go/chat"
	"github.com/PolyTechProjects/chaotic_chat/chat/src/internal/client"
	"github.com/PolyTechProjects/chaotic_chat/chat/src/internal/controller"
	"github.com/PolyTechProjects/chaotic_chat/chat/src/internal/service"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HttpServer struct {
	chatMgmtController *controller.ChatManagementController
}

func NewHttpServer(chatMgmtController *controller.ChatManagementController) *HttpServer {
	return &HttpServer{chatMgmtController: chatMgmtController}
}

func (h *HttpServer) StartServer() {
	http.HandleFunc("GET /chat", h.chatMgmtController.GetAllAvailableChatsHandler)
	http.HandleFunc("POST /chat/room", h.chatMgmtController.CreateChatHandler)
	http.HandleFunc("DELETE /chat/room", h.chatMgmtController.DeleteChatHandler)
	http.HandleFunc("GET /chat/room", h.chatMgmtController.GetChatHandler)
	http.HandleFunc("PUT /chat/room", h.chatMgmtController.UpdateChatHandler)
	http.HandleFunc("DELETE /chat/room/users", h.chatMgmtController.DeleteUsersInChatHandler)
	http.HandleFunc("PUT /chat/room/users", h.chatMgmtController.AddUsersInChatHandler)
	http.HandleFunc("DELETE /chat/room/admins", h.chatMgmtController.DeleteAdminsInChatHandler)
	http.HandleFunc("PUT /chat/room/admins", h.chatMgmtController.AddAdminsInChatHandler)
	http.HandleFunc("DELETE /chat/room/readers", h.chatMgmtController.MakeReadersUsersInChatHandler)
	http.HandleFunc("PUT /chat/room/readers", h.chatMgmtController.MakeUsersReadersInChatHandler)
	http.HandleFunc("POST /chat/room/{joinLink}", h.chatMgmtController.JoinChatHandler)
}

type GRPCServer struct {
	gRPCServer *grpc.Server
	chat_mgmt.UnimplementedChatManagementServer
	service    *service.ChatManagementService
	authClient *client.AuthGRPCClient
}

func New(service *service.ChatManagementService, authClient *client.AuthGRPCClient) *GRPCServer {
	gRPCServer := grpc.NewServer()
	g := &GRPCServer{
		gRPCServer: gRPCServer,
		service:    service,
		authClient: authClient,
	}
	chat_mgmt.RegisterChatManagementServer(gRPCServer, g)
	return g
}

func (s *GRPCServer) Start(l net.Listener) error {
	slog.Debug("Starting gRPC server")
	slog.Debug(l.Addr().String())
	return s.gRPCServer.Serve(l)
}

func (s *GRPCServer) GetChat(ctx context.Context, req *chat_mgmt.GetChatRequest) (*chat_mgmt.ChatRoomResponse, error) {
	slog.Info("GetChat controller started")
	_, err := s.authClient.PerformAuthorize(ctx, nil)
	if err != nil {
		slog.Error(fmt.Sprintf("Authorization error: %v", err.Error()))
		return nil, err
	}

	chatId, err := uuid.Parse(req.ChatId)
	if err != nil {
		slog.Error("Invalid chat Id", "error", err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	chat, err := s.service.GetChat(chatId)
	if err != nil {
		slog.Error("GetChat error", "error", err.Error())
		return nil, err
	}
	return &chat_mgmt.ChatRoomResponse{
		ChatId:          chat.Chat.Id.String(),
		CreatorId:       chat.Chat.CreatorId.String(),
		Name:            chat.Chat.Name,
		Description:     chat.Chat.Description,
		ParticipantsIds: chat.Users,
		AdminsIds:       chat.Admins,
		ReadersIds:      chat.Readers,
	}, nil
}
