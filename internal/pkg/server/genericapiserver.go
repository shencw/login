package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shencw/login/internal/pkg/middleware"
	"github.com/shencw/login/pkg/version"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"strings"
	"time"
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
	s.Use(gin.Logger())

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
	s.insecureServer = &http.Server{
		Addr:    s.InsecureServingInfo.Address,
		Handler: s,
	}
	s.secureServer = &http.Server{
		Addr:    s.SecureServingInfo.Address(),
		Handler: s,
	}

	var eg errgroup.Group
	eg.Go(func() error {
		log.Printf("Start to listening the incoming requests on http address: %s", s.InsecureServingInfo.Address)
		if err := s.insecureServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Println(err.Error())
			return err
		}
		return nil
	})

	eg.Go(func() error {
		key, cert := s.SecureServingInfo.CertKey.KeyFile, s.SecureServingInfo.CertKey.CertFile
		if key == "" || cert == "" || s.SecureServingInfo.BindPort == 0 {
			return nil
		}
		log.Printf("Start to listening the incoming requests on http address: %s", s.SecureServingInfo.Address())
		if err := s.secureServer.ListenAndServeTLS(key, cert); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Println(err.Error())
			return err
		}
		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if s.healthz {
		if err := s.ping(ctx); err != nil {
			return err
		}
	}

	if err := eg.Wait(); err != nil {
		log.Fatal(err.Error())
	}

	return nil
}

func (s *GenericAPIServer) ping(ctx context.Context) error {
	url := fmt.Sprintf("http://%s/healthz", s.InsecureServingInfo.Address)
	if strings.Contains(s.InsecureServingInfo.Address, "0.0.0.0") {
		url = fmt.Sprintf("http://127.0.0.1:%s/healthz", strings.Split(s.InsecureServingInfo.Address, ":")[1])
	}

	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		resp, err := http.DefaultClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			log.Println("The router has been deployed successfully.")

			return resp.Body.Close()
		}

		log.Println("Waiting for the router, retry in 1 second.")
		time.Sleep(1 * time.Second)

		select {
		case <-ctx.Done():
			log.Fatal("can not ping http server within the specified time interval.")
		default:
		}
	}
}

func (s *GenericAPIServer) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Printf("close")

	// todo: 这里有问题 close1 可以 2失败
	if err := s.secureServer.Shutdown(ctx); err != nil {
		log.Printf("Shutdown secure server failed: %s", err.Error())
	}
	log.Printf("close1")

	if err := s.insecureServer.Shutdown(ctx); err != nil {
		log.Printf("Shutdown secure server failed: %s", err.Error())
	}
	log.Printf("close2")

}
