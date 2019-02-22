package restsrv

import "flag"

type Config struct {
	HTTPConn string
}

func NewConfig() *Config {
	cfg := &Config{
		HTTPConn: "127.0.0.1:8777",
	}
	flag.StringVar(&cfg.HTTPConn, "http", cfg.HTTPConn, "REST API http connection")
	return cfg
}
