package server

import (
	"github.com/gin-gonic/gin"
	"net"
	"strconv"
)

type Config struct {
	SecureServing   *SecureServingInfo
	InsecureServing *InsecureServingInfo
	Healthz         bool
	Mode            string
	Middlewares     []string
}

type InsecureServingInfo struct {
	Address string
}

type SecureServingInfo struct {
	BindAddress string
	BindPort int
	CertKey  CertKey
}

func (s *SecureServingInfo) Address() string {
	return net.JoinHostPort(s.BindAddress, strconv.Itoa(s.BindPort))
}

type CertKey struct {
	CertFile string
	KeyFile  string
}

func NewConfig() *Config {
	return &Config{
		Healthz: true,
		Mode:    gin.ReleaseMode,
	}
}

type CompletedConfig struct {
	*Config
}

func (c *Config) Completed() *CompletedConfig {
	return &CompletedConfig{c}
}

func (c *CompletedConfig) New() (*GenericAPIServer, error) {
	s := &GenericAPIServer{
		SecureServingInfo:   c.SecureServing,
		InsecureServingInfo: c.InsecureServing,
		mode:                c.Mode,
		healthz:             c.Healthz,
		middlewares:         c.Middlewares,
		Engine:              gin.New(),
	}

	initGenericAPIServer(s)

	return s, nil
}
