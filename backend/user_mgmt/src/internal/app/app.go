package app

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/PolyTechProjects/chaotic_chat/user_mgmt/src/config"
	"github.com/PolyTechProjects/chaotic_chat/user_mgmt/src/internal/server"
)

type App struct {
	httpServer *server.HttpServer
	httpPort   int
	grpcPort   int
}

func New(httpServer *server.HttpServer, cfg *config.Config) *App {
	return &App{
		httpServer: httpServer,
		httpPort:   cfg.App.HttpInnerPort,
		grpcPort:   cfg.App.GrpcInnerPort,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err.Error())
	}
}

func (a *App) Run() error {
	go a.RunHttpServer()
	return nil
}

func (a *App) RunHttpServer() error {
	hl, err := net.Listen("tcp", fmt.Sprintf(":%d", a.httpPort))
	if err != nil {
		return err
	}
	slog.Debug("Starting HTTP server")
	slog.Debug(hl.Addr().String())
	a.httpServer.StartServer()
	if err := http.Serve(hl, nil); err != nil {
		return err
	}
	return nil
}
