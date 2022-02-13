package apiserver

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func initRouter(g *gin.Engine) {
	installMiddleware(g)
	installController(g)
}

func installMiddleware(g *gin.Engine) {
	return
}

func installController(g *gin.Engine) *gin.Engine {
	g.POST("/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, "")
	})

	return g
}
