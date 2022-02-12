package apiserver

import (
	"github.com/shencw/login/internal/apiserver/config"
	"github.com/shencw/login/internal/apiserver/options"
	"github.com/shencw/login/pkg/app"
)

const commandDesc = `login desc`

func NewApp(basename string) *app.App {
	opts := options.NewOptions()                // 先获取options
	application := app.NewApp("login",    // 创建实例化
		basename,
		app.WithDescription(commandDesc),
		app.WithRunFunc(run(opts)),
	)
	return application
}

func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		cfg, err := config.CreateConfigFromOptions(opts)
		if err != nil {
			return err
		}

		return Run(cfg)
	}
}
