package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/PolyTechProjects/chaotic_chat/user_mgmt/src/config"
	"github.com/PolyTechProjects/chaotic_chat/user_mgmt/src/database"
	"github.com/PolyTechProjects/chaotic_chat/user_mgmt/src/internal/app"
	"github.com/PolyTechProjects/chaotic_chat/user_mgmt/src/internal/client"
	"github.com/PolyTechProjects/chaotic_chat/user_mgmt/src/internal/controller"
	"github.com/PolyTechProjects/chaotic_chat/user_mgmt/src/internal/repository"
	"github.com/PolyTechProjects/chaotic_chat/user_mgmt/src/internal/server"
	"github.com/PolyTechProjects/chaotic_chat/user_mgmt/src/internal/service"
)

func main() {
	cfg := config.MustLoad()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(log)
	database.Init(cfg)
	db := database.DB
	authClient := client.NewAuthClient(cfg)
	repository := repository.New(db)
	service := service.New(repository)
	controller := controller.New(service, authClient)
	httpServer := server.NewHttpServer(controller)
	grpcServer := server.NewGRPCServer(service, authClient)
	app := app.New(httpServer, grpcServer, cfg)
	go app.MustRun()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	defer log.Info("Program successfully finished!")
	defer db.Close()
}
