package database

import (
	"fmt"
	"log/slog"

	"github.com/PolyTechProjects/chaotic_chat/media_handler/src/config"
	"github.com/PolyTechProjects/chaotic_chat/media_handler/src/internal/models"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *gorm.DB

func Init(cfg *config.Config) {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	str := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=%v",
		cfg.Db.UserName,
		cfg.Db.Password,
		cfg.Db.Host,
		cfg.Db.InnerPort,
		cfg.Db.DatabaseName,
		cfg.Db.SslMode,
	)
	slog.Info(str)
	db, err := gorm.Open(
		"postgres",
		str,
	)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	db.AutoMigrate(&models.Media{})
	DB = db
	slog.Info("Connected to DB")
}

func Close() {
	slog.Info("Disconneting from DB")
	DB.Close()
}
