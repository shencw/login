package apiserver

import (
	"github.com/shencw/login/internal/apiserver/config"
	genericAPIServer "github.com/shencw/login/internal/pkg/server"
)

type apiServer struct {
	genericAPIServer *genericAPIServer.GenericAPIServer
}

func createAPIServer(cfg *config.Config) (*apiServer, error) {

	genericConfig, err := buildGenericConfig(cfg)
	if err != nil {
		return nil, err
	}

	genericConfig.Completed().New()

	apiServer := &apiServer{
		genericAPIServer: genericConfig,
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
