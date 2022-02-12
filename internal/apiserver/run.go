package apiserver

import "github.com/shencw/login/internal/apiserver/config"

func Run(cfg *config.Config) error {
	err := createAPIServer(cfg)
	if err != nil {
		return err
	}

	return nil
}
