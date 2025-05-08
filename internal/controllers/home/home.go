package home

import (
	"net/http"

	"github.com/eatmore01/light/internal/config"
	"github.com/eatmore01/light/internal/controllers/auth"
	"github.com/eatmore01/light/internal/shared/constants"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type HomeApi struct {
	Cfg *config.Config
}

func NewHomeAPi(c *config.Config) *HomeApi {
	return &HomeApi{
		Cfg: c,
	}
}

func (h *HomeApi) HomePage(c *gin.Context) {
	tokenString, err := c.Cookie(auth.CookieName)
	if err != nil {
		c.Redirect(http.StatusFound, constants.Routes["login"])
		return
	}

	// Parse and validate token
	token, err := jwt.ParseWithClaims(tokenString, &auth.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.Cfg.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		c.Redirect(http.StatusFound, constants.Routes["login"])
		return
	}

	// Get claims from token
	claims, ok := token.Claims.(*auth.CustomClaims)
	if !ok {
		c.Redirect(http.StatusFound, constants.Routes["login"])
		return
	}

	c.HTML(http.StatusOK, "home.html", gin.H{
		"Username":     claims.Username,
		"ClusterName":  h.Cfg.ClusterName,
		"logoutPath":   constants.Routes["logout"],
		"KubeDownPath": constants.Routes["donwloadkubeconfig"],
	})
}

func AddHomeHandler(r *gin.Engine, h *HomeApi) {

	homeGrp := r.Group("/home")
	homeGrp.GET("/", h.HomePage)
}
