package options

import (
	"github.com/shencw/login/internal/pkg/server"
)

type SecureServingOptions struct {
	BindAddress string          `json:"bind-address" mapstructure:"bind-address"`
	BindPort    int             `json:"bind-port"    mapstructure:"bind-port"`
	ServerCert  GenerateCertKey `json:"tls"          mapstructure:"tls"`
	Required    bool
}

type CertKey struct {
	CertFile string `json:"cert-file"        mapstructure:"cert-file"`
	KeyFile  string `json:"private-key-file" mapstructure:"private-key-file"`
}

type GenerateCertKey struct {
	CertKey       CertKey `json:"cert-key"  mapstructure:"cert-key"`
	CertDirectory string  `json:"cert-dir"  mapstructure:"cert-dir"`
	PairName      string  `json:"pair-name" mapstructure:"pair-name"`
}

func NewSecureServingOptions() *SecureServingOptions {
	return &SecureServingOptions{
		BindAddress: "0.0.0.0",
		BindPort:    8081,
		Required:    true,
		ServerCert: GenerateCertKey{
			CertDirectory: "login",
			PairName:      "/var/run/login",
		},
	}
}

func (s *SecureServingOptions) ApplyTo(c *server.Config) error {
	c.SecureServing = &server.SecureServingInfo{
		BindAddress: s.BindAddress,
		BindPort:    s.BindPort,
		CertKey: server.CertKey{
			CertFile: s.ServerCert.CertKey.CertFile,
			KeyFile:  s.ServerCert.CertKey.KeyFile,
		},
	}
	return nil
}
