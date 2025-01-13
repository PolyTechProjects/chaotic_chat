package server

import (
	"context"
	"log/slog"
	"net"
	"net/http"

	"github.com/PolyTechProjects/chaotic_chat/auth/src/gen/go/auth"
	"github.com/PolyTechProjects/chaotic_chat/auth/src/internal/controller"
	"github.com/PolyTechProjects/chaotic_chat/auth/src/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HttpServer struct {
	authController *controller.AuthController
}

func NewHttpServer(authController *controller.AuthController) *HttpServer {
	return &HttpServer{
		authController: authController,
	}
}

func (h *HttpServer) StartServer() {
	http.HandleFunc("POST /register", h.authController.RegisterHandler)
	http.HandleFunc("POST /login", h.authController.LoginHandler)
}

type GRPCServer struct {
	gRPCServer *grpc.Server
	auth.UnimplementedAuthServer
	authService *service.AuthService
}

func NewGRPCServer(authService *service.AuthService) *GRPCServer {
	gRPCServer := grpc.NewServer()
	g := &GRPCServer{
		gRPCServer:  gRPCServer,
		authService: authService,
	}
	auth.RegisterAuthServer(gRPCServer, g)
	return g
}

func (s *GRPCServer) Start(l net.Listener) error {
	slog.Debug("Starting gRPC server")
	slog.Debug(l.Addr().String())
	return s.gRPCServer.Serve(l)
}

func (s *GRPCServer) Authorize(ctx context.Context, req *auth.AuthorizeRequest) (*auth.AuthorizeResponse, error) {
	accessToken, err := s.authService.Authorize(req.GetAccessToken())
	if err != nil {
		slog.Error(err.Error())
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	return &auth.AuthorizeResponse{AccessToken: accessToken}, nil
}
