package app

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/eatmore01/light/internal/config"
	"github.com/eatmore01/light/internal/controllers/auth"
	"github.com/eatmore01/light/internal/controllers/home"
	kubernetes_controller "github.com/eatmore01/light/internal/controllers/kubernetes"
	"github.com/eatmore01/light/internal/controllers/startup"
	"github.com/eatmore01/light/internal/shared/lg"

	"github.com/gin-gonic/gin"
)

type App struct {
	Server     *http.Server
	Logger     *slog.Logger
	AppCfg     *config.Config
	StartupApi *startup.StartUpApi
	AuthApi    *auth.AuthApi
	HomeApi    *home.HomeApi
	KubeAPi    *kubernetes_controller.KubeApi
}

func NewApp(lg *slog.Logger, appCfg *config.Config, SuApi *startup.StartUpApi, AuthApi *auth.AuthApi, h *home.HomeApi, k *kubernetes_controller.KubeApi) *App {
	return &App{
		Server:     &http.Server{},
		Logger:     lg,
		AppCfg:     appCfg,
		StartupApi: SuApi,
		AuthApi:    AuthApi,
		HomeApi:    h,
		KubeAPi:    k,
	}
}

func (a *App) Run() {
	r := gin.Default()
	r.LoadHTMLGlob(a.AppCfg.TemplatesDir + "/*")

	startup.AddStartUpHandler(r, a.StartupApi)

	auth.AddAuthHandlers(r, a.AuthApi)

	home.AddHomeHandler(r, a.HomeApi)

	kubernetes_controller.AddKubernetesHandler(r, a.KubeAPi)

	address := a.AppCfg.Host + ":" + a.AppCfg.Port
	a.Server.Addr = address
	a.Server.Handler = r

	go func() {
		err := a.Server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			a.Logger.Error("server error: %s", lg.Err(err))
			return

		}
		a.Logger.Info("Server starts at: %s address", address)
	}()
}

func (a *App) GraceFullShutDown() {
	shutDownInterval := 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), shutDownInterval)
	defer cancel()

	err := a.Server.Shutdown(ctx)
	if err != nil {
		a.Logger.Error("server couldnt shutdown with error: %s", lg.Err(err))
	}

}
