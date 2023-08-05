package config

type (
	ServerConfig struct {
		HTTP
	}

	HTTP struct {
		Host string
		Port uint
	}
)

func New() *ServerConfig {
	cfg := &ServerConfig{
		HTTP: HTTP{
			Host: "",
			Port: 8080,
		},
	}
	return cfg
}
