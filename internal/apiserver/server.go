package apiserver

import (
	"github.com/shencw/login/internal/apiserver/config"
	genericAPIServer "github.com/shencw/login/internal/pkg/server"
	"github.com/shencw/login/pkg/shutdown"
	"github.com/shencw/login/pkg/shutdown/managers/posixsignal"
	"log"
)

type apiServer struct {
	gs               *shutdown.GracefulShutdown
	genericAPIServer *genericAPIServer.GenericAPIServer
}

type preparedAPIServer struct {
	*apiServer
}

// PrepareRun apiServer准备启动
func (s *apiServer) PrepareRun() preparedAPIServer {
	// 初始化路由
	initRouter(s.genericAPIServer.Engine)

	s.gs.AddShutdownCallback(shutdown.OnShutdownFunc(func(string) error {
		s.genericAPIServer.Close()

		return nil
	}))

	return preparedAPIServer{s}
}

// Run 真正开始执行
func (p preparedAPIServer) Run() error {
	if err := p.gs.Start(); err != nil {
		log.Fatalf("start shutdown manager failed: %s", err.Error())
	}
	return p.genericAPIServer.Run()
}

func createAPIServer(cfg *config.Config) (*apiServer, error) {
	gs := shutdown.New()
	gs.AddShutdownManager(posixsignal.NewPosixSignalManager())

	genericConfig, err := buildGenericConfig(cfg)
	if err != nil {
		return nil, err
	}

	genericServer, err := genericConfig.Completed().New()

	apiServer := &apiServer{
		genericAPIServer: genericServer,
	}

	return apiServer, nil
}

func buildGenericConfig(cfg *config.Config) (*genericAPIServer.Config, error) {
	genericConfig := genericAPIServer.NewConfig()

	if err := cfg.GenericServerRunOptions.ApplyTo(genericConfig); err != nil {
		return genericConfig, err
	}

	return genericConfig, nil
}
