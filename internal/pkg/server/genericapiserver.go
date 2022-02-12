package server

import (
	"github.com/gin-gonic/gin"
	"github.com/shencw/login/internal/pkg/middleware"
	"github.com/shencw/login/pkg/version"
	"log"
	"net/http"
)

func initGenericAPIServer(s *GenericAPIServer) {
	s.Setup()
	s.InstallMiddlewares()
	s.InstallAPIs()
}

type GenericAPIServer struct {
	middlewares         []string
	mode                string
	healthz             bool
	SecureServingInfo   *SecureServingInfo
	InsecureServingInfo *InsecureServingInfo

	*gin.Engine

	insecureServer, secureServer *http.Server
}

func (s *GenericAPIServer) Setup() {
	gin.SetMode(s.mode)
}

func (s *GenericAPIServer) InstallMiddlewares() {
	s.Use(middleware.RequestID())
	s.Use(middleware.Context())

	for _, m := range s.middlewares {
		mw, ok := middleware.Middlewares[m]
		if !ok {
			log.Printf("can not find middleware: %s", m)

			continue
		}

		log.Printf("install middleware: %s", m)
		s.Use(mw)
	}
}

func (s *GenericAPIServer) InstallAPIs() {
	if s.healthz {
		s.GET("/healthz", func(c *gin.Context) {
			c.JSON(http.StatusOK, map[string]string{"status": "OK"})
		})
	}

	s.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, version.Get())
	})
}

func (s *GenericAPIServer) Run() error {

	return nil
}
