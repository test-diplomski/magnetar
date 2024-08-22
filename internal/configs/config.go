package configs

import (
	"os"
)

type Config struct {
	natsAddress     string
	etcdAddress     string
	serverAddress   string
	oortAddress     string
	meridianAddress string
	gravityAddress  string
	tokenKey        string
}

func (c *Config) NatsAddress() string {
	return c.natsAddress
}

func (c *Config) EtcdAddress() string {
	return c.etcdAddress
}

func (c *Config) ServerAddress() string {
	return c.serverAddress
}

func (c *Config) OortAddress() string {
	return c.oortAddress
}

func (c *Config) MeridianAddress() string {
	return c.meridianAddress
}

func (c *Config) GravityAddress() string {
	return c.gravityAddress
}

func (c *Config) TokenKey() string {
	return c.tokenKey
}

func NewFromEnv() (*Config, error) {
	return &Config{
		natsAddress:     os.Getenv("NATS_ADDRESS"),
		etcdAddress:     os.Getenv("ETCD_ADDRESS"),
		serverAddress:   os.Getenv("MAGNETAR_ADDRESS"),
		oortAddress:     os.Getenv("OORT_ADDRESS"),
		meridianAddress: os.Getenv("MERIDIAN_ADDRESS"),
		gravityAddress:  os.Getenv("GRAVITY_ADDRESS"),
		tokenKey:        os.Getenv("SECRET_KEY"),
	}, nil
}
