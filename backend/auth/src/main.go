package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/PolyTechProjects/chaotic_chat/auth/src/config"
	"github.com/PolyTechProjects/chaotic_chat/auth/src/database"
	"github.com/PolyTechProjects/chaotic_chat/auth/src/internal/app"
	"github.com/PolyTechProjects/chaotic_chat/auth/src/internal/controller"
	"github.com/PolyTechProjects/chaotic_chat/auth/src/internal/repository"
	"github.com/PolyTechProjects/chaotic_chat/auth/src/internal/server"
	"github.com/PolyTechProjects/chaotic_chat/auth/src/internal/service"
	"github.com/PolyTechProjects/chaotic_chat/auth/src/internal/validator"
)

func main() {
	goTest()
	cfg := config.MustLoad()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(log)
	database.Init(cfg)
	db := database.DB
	repository := repository.New(db)
	authService := service.New(repository)
	grpcServer := server.NewGRPCServer(authService)
	authController := controller.NewAuthController(authService)
	httpServer := server.NewHttpServer(authController)
	app := app.New(grpcServer, httpServer, cfg)
	go app.MustRun()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	defer log.Info("Program successfully finished!")
	defer db.Close()
}

func goTest() {
	var err error
	wrongNames := []string{
		"",
		" ",
		"    ",
		string(0x1F),
	}
	for _, name := range wrongNames {
		err = validator.ValidateName(name)
		if err == nil {
			panic("wrong name: {" + name + "} should be invalid but is valid")
		}
	}
	correctNames := []string{
		"normalName",
		"normal-Name",
		"normal_name",
		".normal,name",
	}
	for _, name := range correctNames {
		err = validator.ValidateName(name)
		if err != nil {
			panic(err)
		}
	}

	wrongPasswords := []string{
		"",
		" ",
		"\"",
		"1",
		"!",
		"a",
		"A",
		"ONLYCAPS",
		"0NLYCAP5W1THD1G1T5",
		"ONLY_CAPS_WITH_SPEC_SYMBOL",
		"0NLY_CAP5.W1TH_D1G1T5.AND_5P3C_5YMB075",
		"onlylow",
		"0n7y70ww1thd1g1t5",
		"only_low_with_spec_symbol",
		"0n7y_70w.w1th_d1g1t5.and_5p3c_5ymb075",
		"onlylow",
		"NoSpecialSymbolsAndDigits",
		"N05p3c1a75ymb075ButDigits",
		"12.3_4.5.6_7.89",
		"BackS1ash_At_the_END\\",
		"\"Pa55w0rd_W1th.D1G1T5.AND.5p3c_5ymb075\"",
	}
	for _, password := range wrongPasswords {
		err = validator.ValidatePassword(password)
		if err == nil {
			panic("wrong password: {" + password + "} should be invalid but is valid")
		}
	}
	correctPasswords := []string{
		"Pa55w0rd!",
		"{Pa55w0rd}",
		"Pa55[w0rd]",
		"(Pa55)w0rd",
		"Pa55-w0rd",
		"-=Pa55.w0rd=-",
		"^*P@55_W1t#=D!G!T5=&=$ymb07$*^",
		"OAoa1!",
		"oOAa1!",
		"1OAoa!",
		"!OAoa1",
		"Pass,Word|123%456+789?|~`'",
	}
	for _, password := range correctPasswords {
		err = validator.ValidatePassword(password)
		if err != nil {
			panic(err)
		}
	}
}
