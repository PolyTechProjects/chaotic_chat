package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/PolyTechProjects/chaotic_chat/media_handler/src/config"
	"github.com/PolyTechProjects/chaotic_chat/media_handler/src/database"
	"github.com/PolyTechProjects/chaotic_chat/media_handler/src/internal/app"
	"github.com/PolyTechProjects/chaotic_chat/media_handler/src/internal/client"
	"github.com/PolyTechProjects/chaotic_chat/media_handler/src/internal/controller"
	"github.com/PolyTechProjects/chaotic_chat/media_handler/src/internal/repository"
	"github.com/PolyTechProjects/chaotic_chat/media_handler/src/internal/server"
	"github.com/PolyTechProjects/chaotic_chat/media_handler/src/internal/service"
	"github.com/PolyTechProjects/chaotic_chat/media_handler/src/redis"
)

func main() {
	cfg := config.MustLoad()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(log)
	redis.Init(cfg)
	redisClient := redis.RedisClient
	database.Init(cfg)
	db := database.DB
	authClient := client.New(cfg)
	repository := repository.New(db, redisClient)
	service := service.New(repository, cfg)
	controller := controller.New(service, authClient)
	httpServer := server.NewHttpServer(controller)
	app := app.New(httpServer, cfg)
	go app.MustRun()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	database.Close()
	redis.Close()
}
