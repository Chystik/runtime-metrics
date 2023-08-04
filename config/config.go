package config

type (
	ServerConfig struct {
		Http
	}

	Http struct {
		Host string
		Port uint
	}
)

func New() *ServerConfig {
	cfg := &ServerConfig{
		Http: Http{
			Host: "",
			Port: 8080,
		},
	}
	return cfg
}
