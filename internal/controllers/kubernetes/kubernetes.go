package kubernetes_controller

import (
	"fmt"
	"net/http"

	"github.com/eatmore01/light/internal/config"
	"github.com/eatmore01/light/internal/controllers/auth"
	"github.com/eatmore01/light/internal/services/kubernetes"
	"github.com/eatmore01/light/internal/shared/constants"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"sigs.k8s.io/yaml"
)

type KubeApi struct {
	KubeService *kubernetes.KubeService
	Cfg         *config.Config
}

func NewKubeApi(c *config.Config) *KubeApi {
	kubeService := kubernetes.NewKubeService()
	return &KubeApi{
		KubeService: kubeService,
		Cfg:         c,
	}
}

func (k *KubeApi) DownloadKubeconfigHandler(c *gin.Context) {
	tokenString, err := c.Cookie(auth.CookieName)
	if err != nil {
		c.Redirect(http.StatusFound, constants.Routes["login"])
		return
	}

	token, err := jwt.ParseWithClaims(tokenString, &auth.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(k.Cfg.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		c.Redirect(http.StatusFound, constants.Routes["login"])
		return
	}

	claims, ok := token.Claims.(*auth.CustomClaims)
	if !ok {
		c.Redirect(http.StatusFound, constants.Routes["login"])
		return
	}

	info := k.KubeService.GenerateInfo(k.Cfg, claims)
	fmt.Println(info)

	cfg := k.KubeService.GenerateKubeConfig(info)

	yamlData, err := yaml.Marshal(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate YAML"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=kubeconfig.yaml")
	c.Data(http.StatusOK, "application/x-yaml", yamlData)
	c.Redirect(http.StatusTemporaryRedirect, constants.Routes["home"])
}

func AddKubernetesHandler(r *gin.Engine, k *KubeApi) {

	kubeGrp := r.Group("/kube")
	kubeGrp.GET("/cfgdownload", k.DownloadKubeconfigHandler)
}
