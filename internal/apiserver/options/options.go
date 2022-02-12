package options

import (
	genericoptions "github.com/shencw/login/internal/pkg/options"
	"github.com/shencw/login/internal/pkg/server"
)

type Options struct {
	GenericServerRunOptions *genericoptions.ServerRunOptions `json:"server"   mapstructure:"server"`
}

func NewOptions() *Options {
	return &Options{
		GenericServerRunOptions: genericoptions.NewServerRunOptions(),
	}
}

func (o *Options) ApplyTo(c *server.Config) error {
	return nil
}
