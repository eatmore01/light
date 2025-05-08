package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/eatmore01/light/internal/app"
	"github.com/eatmore01/light/internal/config"
	"github.com/eatmore01/light/internal/controllers/auth"
	"github.com/eatmore01/light/internal/controllers/home"

	kubernetes_controller "github.com/eatmore01/light/internal/controllers/kubernetes"
	"github.com/eatmore01/light/internal/controllers/startup"
	"github.com/eatmore01/light/internal/shared/lg"
)

func main() {
	cfg := config.MustLoad()

	log := lg.SetupLogger("local")

	suapi := startup.NewStartUpApi()
	aa := auth.NewAuthApi(cfg)
	ha := home.NewHomeAPi(cfg)
	k := kubernetes_controller.NewKubeApi(cfg)

	app := app.NewApp(log, cfg, suapi, aa, ha, k)

	app.Run()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	app.Logger.Info("server start shutdown...")
	app.GraceFullShutDown()
	app.Logger.Info("server shutdown success")
}
