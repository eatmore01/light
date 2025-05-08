package startup

import (
	"net/http"

	"github.com/eatmore01/light/internal/shared/constants"
	"github.com/gin-gonic/gin"
)

type StartUpApi struct{}

func NewStartUpApi() *StartUpApi {
	return &StartUpApi{}
}

func (su *StartUpApi) StartUpHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "startup.html", gin.H{
		"loginPath": constants.Routes["login"],
	})
}

func AddStartUpHandler(r *gin.Engine, su *StartUpApi) {

	StartUpGrp := r.Group("/")
	StartUpGrp.GET("/", su.StartUpHandler)
}
