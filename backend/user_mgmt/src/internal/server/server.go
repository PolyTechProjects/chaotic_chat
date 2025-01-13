package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/PolyTechProjects/chaotic_chat/user_mgmt/src/gen/go/user_mgmt"
	"github.com/PolyTechProjects/chaotic_chat/user_mgmt/src/internal/client"
	"github.com/PolyTechProjects/chaotic_chat/user_mgmt/src/internal/controller"
	"github.com/PolyTechProjects/chaotic_chat/user_mgmt/src/internal/service"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

type UserMgmtGRPCServer struct {
	gRPCServer *grpc.Server
	user_mgmt.UnimplementedUserMgmtServer
	userMgmtService *service.UserMgmtService
	authClient      *client.AuthGRPCClient
}

func NewGRPCServer(userMgmtService *service.UserMgmtService, authClient *client.AuthGRPCClient) *UserMgmtGRPCServer {
	gRPCServcer := grpc.NewServer()
	g := &UserMgmtGRPCServer{
		gRPCServer:      gRPCServcer,
		userMgmtService: userMgmtService,
		authClient:      authClient,
	}
	user_mgmt.RegisterUserMgmtServer(g.gRPCServer, g)
	return g
}

func (s *UserMgmtGRPCServer) Start(l net.Listener) error {
	slog.Debug("Starting gRPC server")
	slog.Debug(l.Addr().String())
	return s.gRPCServer.Serve(l)
}

func (s *UserMgmtGRPCServer) AddUser(ctx context.Context, req *user_mgmt.AddUserRequest) (*user_mgmt.UserResponse, error) {
	slog.Info(fmt.Sprintf("Add User %v", req.UserId))
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		slog.Error(err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	user, err := s.userMgmtService.CreateUser(userId, req.Name)
	if err != nil {
		slog.Error(err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	resp := &user_mgmt.UserResponse{
		UserId:      user.Id.String(),
		Name:        user.Name,
		UrlTag:      user.UrlTag,
		Description: user.Description,
		ProfilePic:  user.ProfilePic,
	}
	return resp, nil
}
